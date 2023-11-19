package handlers

import (
	"github.com/gofiber/contrib/websocket"
	filesystemmanager "github.com/mehdiseddik.com/pkg/services/fileSystemManager"
)

// fiber websocket mamager
var ArborescenceClients = make(map[*websocket.Conn]bool)

// RegisterArborescence a new client
func RegisterArborescence(c *websocket.Conn) {
	ArborescenceClients[c] = true
}

// UnregisterArborescence a client
func UnregisterArborescence(c *websocket.Conn) {
	c.Close()
	delete(ArborescenceClients, c)
}

// arborescenceBroadCast a message to all clients
func arborescenceBroadCast(message interface{}) {
	for client := range ArborescenceClients {
		client.WriteJSON(message)
	}
}

func ArborescenceWsHandler(c *websocket.Conn) {
	defer UnregisterArborescence(c)
	RegisterArborescence(c)
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
