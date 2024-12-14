package wsocket

// import (
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// var clients = make(map[*websocket.Conn]bool) // Connected clients
// var broadcast = make(chan Message)           // Channel for broadcasting messages

// // Message defines the structure of the messages exchanged
// type Message struct {
// 	Username  string `json:"username"`
// 	Message   string `json:"message"`
// 	Timestamp string `json:"timestamp"`
// 	Typing    bool   `json:"typing"` // Indicates if the user is typing
// }

// // HandleConnections handles new WebSocket requests from clients
// func HandleConnections(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil) // Upgrade HTTP to WebSocket
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer ws.Close() // Close the WebSocket connection when the function returns

// 	clients[ws] = true // Register the new client
// 	log.Printf("clients: %v", clients)

// 	for {
// 		var msg Message
// 		// Read a new message as JSON and map it to a Message object
// 		err := ws.ReadJSON(&msg)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			delete(clients, ws) // Remove the client from the list if there is an error
// 			break
// 		}
// 		// Send the message to the broadcast channel
// 		fmt.Println(msg)
// 		broadcast <- msg
// 		HandleMessages()
// 	}
// }

// // HandleMessages broadcasts incoming messages to all clients
// func HandleMessages() {
// 	log.Println("HandleMessages running")
// 	for {
// 		// Get the next message from the broadcast channel
// 		msg := <-broadcast
// 		// Send it to every connected client
// 		for client := range clients {
// 			err := client.WriteJSON(msg) // Write message to the client
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				client.Close()          // Close the connection if there's an error
// 				delete(clients, client) // Remove the client
// 			}
// 		}
// 	}
// }
