package multiplayer

import (
	"golbugames/backend/internal/websocket"
	"golbugames/backend/internal/websocket/client"
	"log"
	"net"
	"time"
)

// structure client pour multijoueurs
type Client struct {
	baseClient client.BaseClient
	hub        *Hub
	matchId    string
}

func newClient(conn net.Conn, hub *Hub) *Client {
	return &Client{
		baseClient: *client.NewBaseClient(conn),
		hub:        hub,
	}
}

// Ajouter une logique de traitement de messages si nécessaire
// processMessage traite les messages reçus des clients
// Il peut être utilisé pour gérer les messages de jeu, les commandes, etc...
func (c *Client) processMessage(payload []byte) {

}

func (c *Client) handleFrame(frame websocket.Frame) {
	switch frame.Opcode {
	case websocket.OpcodeClose:
		log.Printf("Client %s closed the connection", c.baseClient.ClientId)
		c.hub.unregister <- c
		return

	case websocket.OpcodePing:
		log.Printf("Received ping from client %s", c.baseClient.ClientId)
		pongFrame := websocket.Pong(frame.Payload)
		c.baseClient.Mu.Lock()
		_, err := c.baseClient.Conn.Write(pongFrame)
		c.baseClient.Mu.Unlock()
		if err != nil {
			log.Printf("Error sending pong to client %s: %v", c.baseClient.ClientId, err)
			return
		}

	case websocket.OpcodePong:
		log.Printf("Received pong from client %s", c.baseClient.ClientId)

	case websocket.OpcodeText, websocket.OpcodeBinary:
		// Premier frame d'un nouveau message
		c.baseClient.CurrentOpcode = frame.Opcode

		if frame.FIN {
			// Message complet en un seul frame
			log.Printf("Received complete %s message from client %s",
				websocket.opcodeToString(frame.Opcode), c.baseClient.ClientId)
			c.baseClient.Send <- frame.Payload
			c.resetFragmentation()
		} else {
			// Début d'un message fragmenté
			log.Printf("Received first frame of fragmented %s message from client %s",
				websocket.opcodeToString(frame.Opcode), c.baseClient.ClientId)

			// Vérification de la taille
			if len(frame.Payload) > MaxMessageSize {
				log.Printf("First frame too large from client %s", c.baseClient.ClientId)
				c.resetFragmentation()
				return
			}

			c.frameBuffer = append(c.frameBuffer[:0], frame.Payload...)
		}

	case websocket.OpcodeContinuation:
		// Frame de continuation
		if c.frameBuffer == nil || c.baseClient.CurrentOpcode == 0 {
			log.Printf("Received continuation frame without initial frame from client %s", c.baseClient.ClientId)
			return
		}

		// Vérification de la taille totale
		if len(c.frameBuffer)+len(frame.Payload) > MaxMessageSize {
			log.Printf("Message too large from client %s", c.baseClient.ClientId)
			c.resetFragmentation()
			return
		}

		c.frameBuffer = append(c.frameBuffer, frame.Payload...)

		if frame.FIN {
			// Message complet
			log.Printf("Received final continuation frame from client %s", c.baseClient.ClientId)
			c.send <- frame.Payload
			c.resetFragmentation()
		} else {
			log.Printf("Received continuation frame from client %s", c.baseClient.ClientId)
		}

	default:
		log.Printf("Received unknown frame type (0x%02x) from client %s", frame.Opcode, c.baseClient.ClientId)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.hub.unregister <- c
		c.baseClient.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			c.baseClient.Mu.Lock()
			_, err := c.baseClient.Conn.Write(message)
			c.baseClient.Mu.Unlock()

			if err != nil {
				return
			}
		case <-ticker.C:
			// Envoyer un ping
			c.baseClient.Mu.Lock()
			_, err := c.baseClient.Conn.Write(websocket.Ping([]byte("ping")))
			c.baseClient.Mu.Unlock()

			if err != nil {
				return
			}
		}
	}

}

func (c *Client) readPump() {

	defer func() {
		c.hub.unregister <- c
		c.baseClient.Conn.Close()
	}()

	buffer := make([]byte, 0, 4096)

	for {
		temp := make([]byte, 1024)
		n, err := c.baseClient.Conn.Read(temp)
		if err != nil {
			log.Printf("Error reading from client %s: %v", c.baseClient.ClientId, err)
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
				log.Printf("Error parsing frame from client %s: %v", c.baseClient.ClientId, err)
				return
			}

			// Traiter la frame parsée
			c.handleFrame(frame)

			// Retirer la frame traitée du buffer
			buffer = buffer[frameLen:]
		}
	}
}
