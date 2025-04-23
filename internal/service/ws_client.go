package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/pkg/auths"
)

const (
	// Время, отведенное на написание сообщения партнеру.
	writeWait = 3600 * time.Second

	// Время, отведенное на прочтение очередного сообщения от партнера.
	pongWait = 10 * time.Second

	// Посылать пинги на пир с этим периодом. Должно быть меньше pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 //512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// HandshakeTimeout: 240,
	Subprotocols: []string{"JSON"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client struct for websocket connection and message sending
type Client struct {
	UserId   string
	RoomId   string
	Conn     *websocket.Conn
	send     chan domain.MessageSocket
	hub      *Hub
	Services *Services
}

// NewClient creates a new client
func NewClient(userId string, roomId string, conn *websocket.Conn, hub *Hub, services *Services) *Client {
	return &Client{UserId: userId, RoomId: roomId, Conn: conn, send: make(chan domain.MessageSocket, 256), hub: hub, Services: services}
}

// Client goroutine to read messages from client
func (c *Client) Read() {

	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	// fmt.Println("Read: ", c.UserId)
	for {
		var msg domain.MessageSocket
		c.Conn.SetReadDeadline(time.Now().Add(writeWait))
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error Read: ", err)
			break
		}

		// fmt.Println("----->", msg)

		switch msg.Type {
		case "exit":
			fmt.Println("Need close socket")
		case "jwt":
			// fmt.Println("jwt: ", &msg, err)
			c.hub.register <- c
			tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
			if err != nil {
				// c.AbortWithError(http.StatusUnauthorized, err)
				c.hub.HandleMessage(domain.MessageSocket{
					Type:      "error",
					Content:   errors.New("Access denied!").Error(),
					ID:        "room1",
					Sender:    "anonymous",
					Recipient: "anonymous",
				})
				// c.hub.RemoveClient(c)
				return
			}

			claims, err := tokenManager.Parse(msg.Content.(string))
			if err != nil {
				// c.AbortWithError(http.StatusUnauthorized, err)
				// appG.ResponseError(http.StatusUnauthorized, err, nil)
				c.hub.HandleMessage(domain.MessageSocket{
					Type:      "error",
					Content:   errors.New("Access denied!").Error(),
					ID:        "room1",
					Sender:    "anonymous",
					Recipient: "anonymous",
				})
				// c.hub.RemoveClient(c)
				return
			}

			_, err = c.Services.Authorization.GetAuth(claims.Subject)
			if err != nil {
				// appG.ResponseError(http.StatusUnauthorized, err, nil)
				c.hub.HandleMessage(domain.MessageSocket{
					Type:      "error",
					Content:   errors.New("Access denied!").Error(),
					ID:        "room1",
					Sender:    "anonymous",
					Recipient: "anonymous",
				})
				// c.hub.RemoveClient(c)
				return
			}

			// fmt.Println("authData: ", authData)
			c.UserId = claims.Uid

			// c.Conn.SetReadDeadline(time.Now().Add(1 * time.Hour))

			status := true
			_, err =
				c.Services.User.UpdateUser(c.UserId, &domain.UserInput{Online: &status, LastTime: time.Now()})
			if err != nil {
				// c.hub.HandleMessage(domain.Message{Type: "message", Sender: c.UserId, Recipient: "user2", Content: user, ID: "room1", Service: "user"})
				c.hub.RemoveClient(c)
			} else {
				c.Conn.SetReadLimit(maxMessageSize)
				c.Conn.SetReadDeadline(time.Now().Add(pongWait))
				// c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
			}

		}

		c.Conn.SetPongHandler(func(string) error {
			c.Conn.SetReadDeadline(time.Now().Add(pongWait))
			// c.hub.HandleMessage(domain.MessagePingPongSocket{Type: "pong"})
			// fmt.Println("Pong")
			return nil
		})
		for {
			_, _, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("error: %v", err)
				}
				break
			}
			// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
			c.hub.broadcast <- msg
		}
		// c.hub.broadcast <- msg
	}
}

// Client goroutine to write messages to client
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	// time.Sleep(12 * time.Second)

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.Conn.WriteJSON(message)
				if err != nil {
					fmt.Println("Error Write: ", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// fmt.Println("Ping")
			// c.hub.HandleMessage(domain.MessageSocket{Type: "ping", Recipient: "", ID: "room1", Content: bson.D{}})
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

// Client закрывает канал для отмены регистрации клиента
func (c *Client) Close() {
	close(c.send)
}

// Функция для обработки соединения через веб-сокет, регистрации клиента в концентраторе и запуска горутин.
func ServeWS(ctx *gin.Context, roomId string, hub *Hub, services *Services) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// userId, _ := ctx.Get("uid")
	// fmt.Println("Connect to room: ", roomId, userId)
	client := NewClient("anonymous", roomId, ws, hub, services)
	hub.register <- client
	go client.Write()
	go client.Read()
}
