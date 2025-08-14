// backend/board-service/main.go
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// In production this would restrict to only allow the frontend's domain:
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// The central Hub:
// clients - The list of all connected Users. A map from pointer to the
//           client's WebSocket connection object to a string representing
//           the board they are on.
type ClientManager struct {
	clients    map[*websocket.Conn]string
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

var Hub = ClientManager{
	clients:    make(map[*websocket.Conn]string),
	broadcast:  make(chan []byte),
	register:   make(chan *websocket.Conn),
	unregister: make(chan *websocket.Conn),
	// The mutex defaults to unlocked.
}

func main() {
	// Create the Hub in a goroutine:
	go Hub.run()
	
	//Create a Gin web server:
	router := gin.Default()
	router.GET("/ws/board/:boardId", serveWs)
	log.Println("Board service starting on port 8082")
	router.Run(":8082")
}

func (manager *ClientManager) run() {
	for {
		select {
		case connection := <-manager.register:
			manager.mu.Lock()
			
			//For now, there's just one board named "general":
			manager.clients[connection] = "general"
			manager.mu.Unlock()
			log.Println("Connection registered")

		case connection := <-manager.unregister:
			manager.mu.Lock()
			if _, ok := manager.clients[connection]; ok {
				delete(manager.clients, connection)
				log.Println("Connection unregistered")
			}
			manager.mu.Unlock()

		case message := <-manager.broadcast:
			manager.mu.Lock()
			for connection := range manager.clients {
				if err := connection.WriteMessage(websocket.TextMessage, message); err != nil {
					manager.unregister <- connection
					connection.Close()
				}
			}
			manager.mu.Unlock()
		}
	}
}

func serveWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	Hub.register <- conn

	// Ensure that the client is always unregistered and the connection is
	// closed, even if an error occurs:
	defer func() {
		Hub.unregister <- conn
		conn.Close()
	}()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		Hub.broadcast <- msg
	}
}