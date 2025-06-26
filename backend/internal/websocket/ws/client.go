package ws

import (
	"log"
	"net"
	"sync"
	"time"
)

type Client struct {
	clientId      string
	nickname      string
	conn          net.Conn
	mu            sync.Mutex
	send          chan []byte
	hub           *Hub
	matchId       string
	frameBuffer   []byte
	currentOpcode byte
}

const (
	pongWait       = 60 * time.Second // Durée d'attente pour un pong
	newline        = "\n"
	space          = " "
	pingPeriod     = (pongWait * 9) / 10 // Période de ping pour garder la connexion active
	pongTimeout    = 60 * time.Second    // Durée d'attente pour un pong avant de fermer la connexion
	MaxMessageSize = 1024 * 1024
)

func newClient(conn net.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
		hub:  hub,
	}
}

func (c *Client) resetFragmentation() {
	c.frameBuffer = nil
	c.currentOpcode = 0
}

// Ajouter une logique de traitement de messages si nécessaire
// processMessage traite les messages reçus des clients
// Il peut être utilisé pour gérer les messages de jeu, les commandes, etc...
func (c *Client) processMessage(payload []byte) {

}

func (c *Client) handleFrame(frame Frame) {
	switch frame.Opcode {
	case OpcodeClose:
		log.Printf("Client %s closed the connection", c.clientId)
		c.hub.unregister <- c
		return

	case OpcodePing:
		log.Printf("Received ping from client %s", c.clientId)
		pongFrame := Pong(frame.Payload)
		c.mu.Lock()
		_, err := c.conn.Write(pongFrame)
		c.mu.Unlock()
		if err != nil {
			log.Printf("Error sending pong to client %s: %v", c.clientId, err)
			return
		}

	case OpcodePong:
		log.Printf("Received pong from client %s", c.clientId)

	case OpcodeText, OpcodeBinary:
		// Premier frame d'un nouveau message
		c.currentOpcode = frame.Opcode

		if frame.FIN {
			// Message complet en un seul frame
			log.Printf("Received complete %s message from client %s",
				opcodeToString(frame.Opcode), c.clientId)
			c.send <- frame.Payload
			c.resetFragmentation()
		} else {
			// Début d'un message fragmenté
			log.Printf("Received first frame of fragmented %s message from client %s",
				opcodeToString(frame.Opcode), c.clientId)

			// Vérification de la taille
			if len(frame.Payload) > MaxMessageSize {
				log.Printf("First frame too large from client %s", c.clientId)
				c.resetFragmentation()
				return
			}

			c.frameBuffer = append(c.frameBuffer[:0], frame.Payload...)
		}

	case OpcodeContinuation:
		// Frame de continuation
		if c.frameBuffer == nil || c.currentOpcode == 0 {
			log.Printf("Received continuation frame without initial frame from client %s", c.clientId)
			return
		}

		// Vérification de la taille totale
		if len(c.frameBuffer)+len(frame.Payload) > MaxMessageSize {
			log.Printf("Message too large from client %s", c.clientId)
			c.resetFragmentation()
			return
		}

		c.frameBuffer = append(c.frameBuffer, frame.Payload...)

		if frame.FIN {
			// Message complet
			log.Printf("Received final continuation frame from client %s", c.clientId)
			c.send <- frame.Payload
			c.resetFragmentation()
		} else {
			log.Printf("Received continuation frame from client %s", c.clientId)
		}

	default:
		log.Printf("Received unknown frame type (0x%02x) from client %s", frame.Opcode, c.clientId)
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
			_, err := c.conn.Write(Ping([]byte("ping")))
			c.mu.Unlock()

			if err != nil {
				return
			}
		}
	}

}

func (c *Client) readPump() {

	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	buffer := make([]byte, 0, 4096)

	for {
		temp := make([]byte, 1024)
		n, err := c.conn.Read(temp)
		if err != nil {
			log.Printf("Error reading from client %s: %v", c.clientId, err)
			return
		}

		if n == 0 {
			continue
		}

		buffer = append(buffer, temp[:n]...)

		// Traiter toutes les frames complètes dans le buffer
		for len(buffer) > 0 {
			frame, frameLen, err := parseFrame(buffer)
			if err != nil {
				if err == ErrIncompleteFrame {
					// Frame incomplète, attendre plus de données
					break
				}
				log.Printf("Error parsing frame from client %s: %v", c.clientId, err)
				return
			}

			// Traiter la frame parsée
			c.handleFrame(frame)

			// Retirer la frame traitée du buffer
			buffer = buffer[frameLen:]
		}
	}
}
