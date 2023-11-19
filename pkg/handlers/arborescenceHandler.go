package handlers

import (
	"github.com/gofiber/contrib/websocket"
	filesystemmanager "github.com/mehdiseddik.com/pkg/services/fileSystemManager"
)

// fiber websocket mamager
var Clients = make(map[*websocket.Conn]bool)

// Register a new client
func Register(c *websocket.Conn) {
	Clients[c] = true
}

// Unregister a client
func Unregister(c *websocket.Conn) {
	c.Close()
	delete(Clients, c)
}

// Broadcast a message to all clients
func Broadcast(message interface{}) {
	for client := range Clients {
		client.WriteJSON(message)
	}
}

// Send a message to a specific client
func Send(c *websocket.Conn, message []byte) {
	c.WriteMessage(1, message)
}

// Send a message to a specific client
func SendJSON(c *websocket.Conn, message interface{}) {
	c.WriteJSON(message)
}

func Arborescence(c *websocket.Conn) {
	defer Unregister(c)
	Register(c)
	c.WriteJSON(filesystemmanager.RootFolder)
	for {
		// when receive a message from client, broadcast to all clients
		_, _, err := c.ReadMessage()
		if err != nil {
			return
		}
		c.WriteJSON(filesystemmanager.RootFolder)
	}
}
