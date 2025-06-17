package ws

import (
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	clientId string
	conn     net.Conn
	mu       sync.Mutex
	send     chan []byte
	hub      *Hub
	matchId  string
}

const (
	pongWait    = 60 * time.Second // Durée d'attente pour un pong
	newline     = "\n"
	space       = " "
	pingPeriod  = (pongWait * 9) / 10 // Période de ping pour garder la connexion active
	pongTimeout = 60 * time.Second    // Durée d'attente pour un pong avant de fermer la connexion
)

func newClient(conn net.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  hub,
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			c.mu.Lock()
			_, err := c.conn.Write(message)
			c.mu.Unlock()

			if err != nil {
				return
			}
		case <-ticker.C:
			// Envoyer un ping
			c.mu.Lock()
			_, err := c.conn.Write([]byte{0x9})
			c.mu.Unlock()

			if err != nil {
				return
			}
		}
	}

}

func (c *Client) readPump() {
	// Il faut intégrer le décodage des messages et la gestion des erreurs
	// ainsi que la gestion des messages de type ping/pong
	// ----- WIP -----
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		frame := make([]byte, 1024)
		n, err := c.conn.Read(frame)
		if err != nil {
			log.Printf("Error reading from client %s: %v", c.clientId, err)
		}
		message := frame[:n]
		c.hub.broadcast <- message
	}
}

//  Implémenter l'attribution du matchmaking
//  l'attribution des ID et la création des hubs en conséquence
