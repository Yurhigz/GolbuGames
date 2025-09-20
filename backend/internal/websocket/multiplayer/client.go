package multiplayer

import (
	"golbugames/backend/internal/websocket"
	"golbugames/backend/internal/websocket/client"
	"log"
	"net"
	"sync"
	"time"
)

// structure client pour multijoueurs
type Client struct {
	baseClient *client.BaseClient
	hub        *Hub
	hubManager *HubManager
	matchId    string
	queueTime  time.Time
	mu         sync.Mutex
	closed     bool
}

// Création d'un nouveau client

func newClient(conn net.Conn, hubManager *HubManager) *Client {
	return &Client{
		baseClient: client.NewBaseClient(conn),
		hubManager: hubManager,
		closed:     false,
	}
}

// cleanup gère la fermeture propre du client
func (c *Client) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}
	c.closed = true

	log.Printf("Cleaning up client %s", c.baseClient.ClientId)

	// Retirer le client de la queue si il y est
	c.hubManager.RemoveClientFromQueue(c)

	// Désenregistrer du hub si assigné
	if c.hub != nil {
		select {
		case c.hub.unregister <- c:
		default:
			log.Printf("Hub unregister channel blocked for client %s", c.baseClient.ClientId)
		}
	}

	// Fermer la connexion
	if c.baseClient != nil && c.baseClient.Conn != nil {
		c.baseClient.Conn.Close()
	}

	// Fermer le canal Send si il existe
	if c.baseClient != nil && c.baseClient.Send != nil {
		close(c.baseClient.Send)
	}
}

// Gestion des frames reçues

func (c *Client) handleFrame(frame websocket.Frame) {
	switch frame.Opcode {
	case websocket.OpcodeClose:
		log.Printf("Client %s closed the connection", c.baseClient.ClientId)
		if c.hub != nil {
			// vérifier pourquoi le hub n'est pas défini
			c.hub.unregister <- c
		}
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
			log.Printf("Received complete %s message from client %s", websocket.OpcodeToString(frame.Opcode), c.baseClient.ClientId)
			c.baseClient.Send <- c.baseClient.DuplicateFrame(&frame)
			c.baseClient.ResetFragmentation()
		} else {
			// Début d'un message fragmenté
			log.Printf("Received first frame of fragmented %s message from client %s",
				websocket.OpcodeToString(frame.Opcode), c.baseClient.ClientId)

			// Vérification de la taille
			if len(frame.Payload) > client.MaxMessageSize {
				log.Printf("First frame too large from client %s", c.baseClient.ClientId)
				c.baseClient.ResetFragmentation()
				return
			}

			c.baseClient.FrameBuffer = append(c.baseClient.FrameBuffer[:0], frame.Payload...)
		}

	case websocket.OpcodeContinuation:
		// Frame de continuation
		if c.baseClient.FrameBuffer == nil || c.baseClient.CurrentOpcode == 0 {
			log.Printf("Received continuation frame without initial frame from client %s", c.baseClient.ClientId)
			return
		}

		// Vérification de la taille totale
		if len(c.baseClient.FrameBuffer)+len(frame.Payload) > client.MaxMessageSize {
			log.Printf("Message too large from client %s", c.baseClient.ClientId)
			c.baseClient.ResetFragmentation()
			return
		}

		c.baseClient.FrameBuffer = append(c.baseClient.FrameBuffer, frame.Payload...)

		if frame.FIN {
			// Message complet
			log.Printf("Received final continuation frame from client %s", c.baseClient.ClientId)
			c.baseClient.Send <- c.baseClient.DuplicateFrame(&frame)
			c.baseClient.ResetFragmentation()
		} else {
			log.Printf("Received continuation frame from client %s", c.baseClient.ClientId)
		}

	default:
		log.Printf("Received unknown frame type (0x%02x) from client %s", frame.Opcode, c.baseClient.ClientId)
	}
}

// WritePump envoie les messages du canal Send au client

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		if c.hub != nil {
			c.hub.unregister <- c
			c.baseClient.Conn.Close()
		}
	}()

	for {
		select {
		case message, ok := <-c.baseClient.Send:
			if !ok {
				return
			}
			c.baseClient.Mu.Lock()
			_, err := c.baseClient.Conn.Write(message.Payload)
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

// ReadPump lit les messages du client et les traite

func (c *Client) readPump() {

	defer func() {
		if c.hub != nil {
			c.hub.unregister <- c
			c.baseClient.Conn.Close()
		}

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
			frame, frameLen, err := websocket.ParseFrame(buffer)
			if err != nil {
				if err == websocket.ErrIncompleteFrame {
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
