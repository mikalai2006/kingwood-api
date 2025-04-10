package service

import (
	"fmt"
	"time"

	"github.com/mikalai2006/kingwood-api/internal/domain"
)

// Hub — это структура, содержащая всех клиентов и отправляемые им сообщения.
type Hub struct {
	// Зарегистрированные клиенты.
	clients map[string]map[*Client]bool
	// Незарегистрированные клиенты.
	unregister chan *Client
	// Регистрация заявок от клиентов.
	register chan *Client
	// Входящие сообщения от клиентов.
	broadcast chan domain.MessageSocket
}

// // Message struct to hold message data
// type Message struct {
// 	Type      string `json:"type"`
// 	Sender    string `json:"sender"`
// 	Recipient string `json:"recipient"`
// 	Content   string `json:"content"`
// 	ID        string `json:"id"`
// }

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		unregister: make(chan *Client),
		register:   make(chan *Client),
		broadcast:  make(chan domain.MessageSocket),
	}
}

// Основная функция для запуска хаба
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.register:
			h.RegisterNewClient(client)
			// Unregister a client.
		case client := <-h.unregister:
			h.RemoveClient(client)
			// Broadcast a message to all clients.
		case message := <-h.broadcast:

			//Check if the message is a type of "message"
			h.HandleMessage(message)

		}
	}
}

// функция проверяет, существует ли комната, и если нет, создайте ее и добавьте в нее клиента
func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.clients[client.RoomId]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.RoomId] = connections
	}
	h.clients[client.RoomId][client] = true

	fmt.Println("RegisterNewClient: Size of clients: ", len(h.clients[client.RoomId]))
}

// func (h *Hub) ValidateClient(client *Client) {
// 	status := true
// 	user, err :=
// 		client.Services.User.UpdateUser(client.UserId, &domain.UserInput{Online: &status})
// 	if err == nil {
// 		h.HandleMessage(domain.Message{Type: "message", Sender: client.UserId, Recipient: "user2", Content: user, ID: "room1", Service: "user"})
// 	}
// 	fmt.Println("ValidateClient: ", len(h.clients[client.RoomId]))
// }

// function to remvoe client from room
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.clients[client.RoomId]; ok {
		status := false
		user, err :=
			client.Services.User.UpdateUser(client.UserId, &domain.UserInput{Online: &status, LastTime: time.Now()})
		if err == nil {
			h.HandleMessage(domain.MessageSocket{Type: "message", Sender: client.UserId, Recipient: "", Content: user, ID: "room1", Service: "user"})
		}

		delete(h.clients[client.RoomId], client)
		close(client.send)
		fmt.Println("Removed client", client.UserId)
	}
}

// function to handle message based on type of message
func (h *Hub) HandleMessage(message domain.MessageSocket) {

	//Check if the message is a type of "message"
	if message.Type == "message" || message.Type == "error" {
		clients := h.clients[message.ID]
		// fmt.Println("===============MESSAGE====================")
		// fmt.Println("len clients=", len(clients))
		// fmt.Println("Recipient=", message.Recipient)
		for client := range clients {
			if client.UserId == message.Recipient || message.Recipient == "" {
				// fmt.Println("Send message client=====>", client.UserId)
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients[message.ID], client)
				}
			}
		}
		// fmt.Println("===========================================")
	}

	//Check if the message is a type of "notification"
	if message.Type == "notification" {
		// fmt.Println("Notification: ", message.Content)
		clients := h.clients[message.Recipient]
		for client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients[message.Recipient], client)
			}
		}
	}

}
