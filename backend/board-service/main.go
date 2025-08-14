// backend/board-service/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors" // We need CORS here too
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// Use the same secret key as the auth-service.
var jwtKey = []byte("my_super_secret_key_that_is_long_and_secure")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow our frontend origin
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// ChatMessage defines the structure of our real-time messages.
type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	hub := newHub()
	go hub.run()

	router := gin.Default()
	
    // We need CORS here because the initial connection is over HTTP
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	router.Use(cors.New(config))

	router.GET("/ws/board/:boardId", func(c *gin.Context) {
		serveWs(hub, c)
	})
	log.Println("Board service starting on port 8082")
	router.Run(":8082")
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		chatMsg := ChatMessage{
			Username: c.username,
			Message:  string(message),
		}
		jsonMsg, _ := json.Marshal(chatMsg)
		c.hub.broadcast <- jsonMsg
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, c *gin.Context) {
	// --- AUTHENTICATION LOGIC ---
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Auth token not provided"})
		return
	}
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	// --- END AUTHENTICATION LOGIC ---

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: claims.Subject, // Use username from the validated token
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}