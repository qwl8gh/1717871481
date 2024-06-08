package ws

import (
	"log"
	"net/http"
	"time"
	"web-messaging-service/pkg/api/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	clients    = make(map[*websocket.Conn]bool)
	broadcast  = make(chan models.Message)
	register   = make(chan *websocket.Conn)
	unregister = make(chan *websocket.Conn)
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	log.Println("WebSocket connection established:", ws.RemoteAddr())

	clients[ws] = true

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			delete(clients, ws)
			break
		}
		msg.Timestamp = time.Now()
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		select {
		case client := <-register:
			clients[client] = true
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				delete(clients, client)
				err := client.Close()
				if err != nil {
					log.Fatalf("Error writing JSON: %v", err)
				}
			}
		case message := <-broadcast:
			log.Printf("Broadcasting message: %v", message)
			for client := range clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Fatalf("Error writing JSON: %v", err)
					delete(clients, client)
					_ = client.Close()
				}
			}
		}
	}
}

func NotifyClients(message models.Message) {
	broadcast <- message
}

func SetupWebSocket(router *gin.Engine) {
	router.GET("/ws", handleConnections)

	go handleMessages()
}
