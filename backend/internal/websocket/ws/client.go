package ws

import (
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

type Client struct {
	clientId      string
	nickname      string
	conn          net.Conn
	mu            sync.Mutex
	send          chan *Frame
	hub           *Hub
	hubManager    *HubManager
	matchId       string
	frameBuffer   []byte
	currentOpcode byte
	elo           int
	queueTime     time.Time
}

const (
	pongWait       = 60 * time.Second // Durée d'attente pour un pong
	newline        = "\n"
	space          = " "
	pingPeriod     = (pongWait * 9) / 10 // Période de ping pour garder la connexion active
	pongTimeout    = 60 * time.Second    // Durée d'attente pour un pong avant de fermer la connexion
	MaxMessageSize = 1024 * 1024
)

func newClient(conn net.Conn, hubmanager *HubManager) *Client {
	return &Client{
		conn:       conn,
		send:       make(chan *Frame, 256),
		hubManager: hubmanager,
	}
}

func (c *Client) compatibleRanking(other *Client) bool {
	return math.Abs(float64(c.elo)-float64(other.elo)) <= 100
}

func (c *Client) resetFragmentation() {
	c.frameBuffer = nil
	c.currentOpcode = 0
}

// Ajouter une logique de traitement de messages si nécessaire
// processMessage traite les messages reçus des clients
// Il peut être utilisé pour gérer les messages de jeu, les commandes, etc...
// créer une frame à partir des messages reçus dans le channel car sinon cela créé des erreurs protocoles car le message est mal formaté

func (c *Client) handleFrame(frame Frame) {
	log.Printf("=== FRAME RECEIVED ===")
	log.Printf("Client: %s", c.clientId)
	log.Printf("Opcode: 0x%x (%s)", frame.Opcode, opcodeToString(frame.Opcode))
	log.Printf("FIN: %t", frame.FIN)
	log.Printf("Payload length: %d", len(frame.Payload))
	log.Printf("Opcode reçu: 0x%x", frame.Opcode)
	switch frame.Opcode {
	case OpcodeClose:
		log.Printf("Client %s closed the connection", c.clientId)
		fmt.Printf("Fermeture du client")
		if frame.Opcode == OpcodeClose && len(frame.Payload) >= 2 {
			// Les 2 premiers bytes d'un close frame contiennent le code de fermeture
			closeCode := (uint16(frame.Payload[0]) << 8) | uint16(frame.Payload[1])
			reason := ""
			if len(frame.Payload) > 2 {
				reason = string(frame.Payload[2:])
			}
			log.Printf("Close code: %d, reason: %s", closeCode, reason)
		}
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
			c.send <- &frame
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
			c.send <- &frame
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
		if c.hub != nil {
			log.Printf("writePump closing for client %s", c.clientId)
			c.hub.unregister <- c
		} else {
			c.hubManager.mu.Lock()
			c.hubManager.RemoveClientFromQueue(c)
			c.hubManager.mu.Unlock()
		}

		c.conn.Close()
	}()

	for {
		select {
		case frame, ok := <-c.send:
			if !ok {
				log.Printf("send channel closed for client %s", c.clientId)
				return
			}
			c.mu.Lock()
			// fmt.Printf("WritePump: about to write frame for client %s", c.clientId)
			_, err := c.conn.Write(frame.ToBytes())
			// fmt.Printf("WritePump: finished writing frame for client %s", c.clientId)
			c.mu.Unlock()

			if err != nil {
				fmt.Printf("Erreur dans le select du writepump - message")
				return
			}
		case <-ticker.C:
			// Envoyer un ping
			c.mu.Lock()
			_, err := c.conn.Write(Ping([]byte("ping")))
			c.mu.Unlock()

			if err != nil {
				fmt.Printf("Erreur dans le select du writepump - ticker")
				return
			}
		}
	}

}

func (c *Client) readPump() {
	log.Printf("readPump started for client %s", c.clientId)
	defer func() {
		log.Printf("readPump closing for client %s", c.clientId)
		if c.hub != nil {
			c.hub.unregister <- c
		} else {
			c.hubManager.mu.Lock()
			c.hubManager.RemoveClientFromQueue(c)
			c.hubManager.mu.Unlock()
		}

		c.conn.Close()
	}()

	buffer := make([]byte, 0, 4096)
	readCount := 0

	for {
		temp := make([]byte, 1024)
		// log.Printf("📖 Attempting to read from client %s (read #%d)", c.clientId, readCount+1)

		n, err := c.conn.Read(temp)
		readCount++
		if err != nil {
			// log.Printf("❌ Error reading from client %s (read #%d): %v", c.clientId, readCount, err)
			// log.Printf("📊 Total successful reads before error: %d", readCount-1)
			log.Printf("❌ <readpump> Error reading from client %s: %v", c.clientId, err)
			return
		}

		log.Printf("✅ Read %d bytes from client %s (read #%d)", n, c.clientId, readCount)
		if n == 0 {
			log.Printf("⚠️ Read 0 bytes from client %s, continuing...", c.clientId)
			continue
		}

		buffer = append(buffer, temp[:n]...)
		// log.Printf("📝 Buffer now contains %d bytes", len(buffer))

		// Traiter toutes les frames complètes dans le buffer
		frameCount := 0
		for len(buffer) > 0 {
			frameCount++
			// log.Printf("🔍 Parsing frame #%d from buffer (%d bytes available)", frameCount, len(buffer))
			frame, frameLen, err := parseFrame(buffer)
			if err != nil {
				log.Printf("parseFrame error: %v", err)
				if err == ErrIncompleteFrame {
					// Frame incomplète, attendre plus de données
					log.Printf("⏳ Incomplete frame, waiting for more data...")
					break
				}
				log.Printf("Error parsing frame from client %s: %v", c.clientId, err)
				return
			}
			// log.Printf("✅ Parsed frame #%d: opcode=0x%x, length=%d bytes", frameCount, frame.Opcode, frameLen)
			// Traiter la frame parsée
			c.handleFrame(frame)

			// Retirer la frame traitée du buffer
			buffer = buffer[frameLen:]
			// log.Printf("📝 Removed %d bytes from buffer, %d bytes remaining", frameLen, len(buffer))
		}
	}
}
