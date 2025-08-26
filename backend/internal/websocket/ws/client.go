package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

// structure client pour solo

type SoloClient struct {
	clientId    string
	pseudo      string
	conn        net.Conn
	mu          sync.Mutex
	elo         int
	frameBuffer []byte
	send        chan *Frame
	solution    string
	playable    string
}

func newSoloClient(conn net.Conn) *SoloClient {
	return &SoloClient{
		conn: conn,
		send: make(chan *Frame, 256),
	}
}

func (c *SoloClient) ValidateMove(index int, value byte) bool {
	// Vérifie si la valeur correspond à la solution
	if c.solution[index] == value {
		return true
	}
	return false
}

func (c *SoloClient) handleFrame(frame Frame) {
	switch frame.Opcode {
	case OpcodeClose:
		log.Printf("[INFO] Client %s closed the connection", c.clientId)
		log.Printf("[INFO] Fermeture du client")
		c.send <- CloseFrame(1000, "Normal Closure")
		return

	case OpcodePing:
		log.Printf("[INFO] Received ping from client %s", c.clientId)
		pongFrame := Pong(frame.Payload)
		c.mu.Lock()
		_, err := c.conn.Write(pongFrame)
		c.mu.Unlock()
		if err != nil {
			log.Printf("[ERR] Error sending pong to client %s: %v", c.clientId, err)
			return
		}

	case OpcodePong:
		log.Printf("[INFO] Received pong from client %s", c.clientId)

	case OpcodeText, OpcodeBinary:

		var move struct {
			Position int  `json:"position"`
			Value    byte `json:"value"`
		}
		if err := json.Unmarshal(frame.Payload, &move); err != nil {
			log.Printf("[ERR] Invalid JSON: %v", err)
			return
		}

		valid := c.ValidateMove(move.Position, move.Value)

		resp := &Frame{
			Opcode: OpcodeText,
			FIN:    true,
			Payload: []byte(fmt.Sprintf(
				`{"position":%d,"value":%d,"valid":%t}`,
				move.Position, move.Value, valid,
			)),
		}

		c.send <- resp

	default:
		log.Printf("[INFO] Received unknown frame type (0x%02x) from client %s", frame.Opcode, c.clientId)
	}
}

func (c *SoloClient) writePump() {
	log.Printf("[INFO] writePump started for client %s", c.clientId)
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case frame, ok := <-c.send:
			if !ok {
				log.Printf("[INFO] send channel closed for client %s", c.clientId)
				return
			}
			c.mu.Lock()
			_, err := c.conn.Write(frame.ToBytes())
			c.mu.Unlock()
			if err != nil {
				log.Printf("[ERR] Error writing to client %s: %v", c.clientId, err)
				return
			}

			if frame.Opcode == OpcodeClose {
				log.Printf("[INFO] Sent close frame to client %s, closing writePump", c.clientId)
				return
			}
		case <-ticker.C:
			c.mu.Lock()
			_, err := c.conn.Write(Ping([]byte("ping")))
			c.mu.Unlock()

			if err != nil {
				fmt.Printf("[ERR] Erreur dans le select du writepump - ticker")
				return
			}

		}
	}
}

func (c *SoloClient) readPump() {
	buffer := make([]byte, 0, 4096)

	defer func() {
		log.Printf("[INFO] Closing connection for client %s", c.clientId)
		c.conn.Close()
		close(c.send)
	}()

	for {
		temp := make([]byte, 1024)
		n, err := c.conn.Read(temp)
		if err != nil {
			if err == io.EOF {
				log.Printf("[INFO] Client %s disconnected", c.clientId)
			} else {
				log.Printf("[ERR] <readPump> Error reading from client %s: %v", c.clientId, err)
			}
			return
		}

		if n == 0 {
			continue
		}

		buffer = append(buffer, temp[:n]...)

		// Traiter toutes les frames complètes dans le buffer
		for len(buffer) > 0 {
			frame, framelen, err := parseFrame(buffer)
			if err != nil {
				log.Printf("parseFrame error: %v", err)
				if err == ErrIncompleteFrame {
					log.Printf("[INFO] Incomplete frame, waiting for more data...")
					break
				}
				log.Printf("[ERR] Error parsing frame from client %s: %v", c.clientId, err)
				return
			}

			buffer = buffer[framelen:]
			c.handleFrame(frame)
		}
	}
}

// structure client pour multijoueurs
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
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	pongTimeout    = 60 * time.Second
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

func (c *Client) handleFrame(frame Frame) {
	switch frame.Opcode {
	case OpcodeClose:
		log.Printf("[INFO] Client %s closed the connection", c.clientId)
		log.Printf("[INFO] Fermeture du client")
		if frame.Opcode == OpcodeClose && len(frame.Payload) >= 2 {
			// Les 2 premiers bytes d'un close frame contiennent le code de fermeture
			closeCode := (uint16(frame.Payload[0]) << 8) | uint16(frame.Payload[1])
			reason := ""
			if len(frame.Payload) > 2 {
				reason = string(frame.Payload[2:])
			}
			log.Printf("[INFO] Close code: %d, reason: %s", closeCode, reason)
		}
		log.Printf("[DEBUG] Fermeture du client %s", c.clientId)
		c.hub.unregister <- c
		return

	case OpcodePing:
		log.Printf("[INFO] Received ping from client %s", c.clientId)
		pongFrame := Pong(frame.Payload)
		c.mu.Lock()
		_, err := c.conn.Write(pongFrame)
		c.mu.Unlock()
		if err != nil {
			log.Printf("[ERR] Error sending pong to client %s: %v", c.clientId, err)
			return
		}

	case OpcodePong:
		log.Printf("[INFO] Received pong from client %s", c.clientId)

	case OpcodeText, OpcodeBinary:
		// Premier frame d'un nouveau message
		c.currentOpcode = frame.Opcode

		if frame.FIN {
			// Message complet en un seul frame
			log.Printf("[INFO] Received complete %s message from client %s",
				opcodeToString(frame.Opcode), c.clientId)
			if c.hub != nil {
				c.hub.broadcast <- &frame
			}
			// c.send <- &frame
			c.resetFragmentation()
		} else {
			// Début d'un message fragmenté
			log.Printf("[INFO] Received first frame of fragmented %s message from client %s",
				opcodeToString(frame.Opcode), c.clientId)

			// Vérification de la taille
			if len(frame.Payload) > MaxMessageSize {
				log.Printf("[ERR] First frame too large from client %s", c.clientId)
				c.resetFragmentation()
				return
			}

			c.frameBuffer = append(c.frameBuffer[:0], frame.Payload...)
		}

	case OpcodeContinuation:
		// Frame de continuation
		if c.frameBuffer == nil || c.currentOpcode == 0 {
			log.Printf("[INFO] Received continuation frame without initial frame from client %s", c.clientId)
			return
		}

		// Vérification de la taille totale
		if len(c.frameBuffer)+len(frame.Payload) > MaxMessageSize {
			log.Printf("[INFO] Message too large from client %s", c.clientId)
			c.resetFragmentation()
			return
		}

		c.frameBuffer = append(c.frameBuffer, frame.Payload...)

		if frame.FIN {
			// Message complet
			log.Printf("[INFO] Received final continuation frame from client %s", c.clientId)
			if c.hub != nil {
				c.hub.broadcast <- &frame
			}
			// c.send <- &frame
			c.resetFragmentation()
		} else {
			log.Printf("[INFO] Received continuation frame from client %s", c.clientId)
		}

	default:
		log.Printf("[INFO] Received unknown frame type (0x%02x) from client %s", frame.Opcode, c.clientId)
	}
}

func (c *Client) writePump() {
	log.Printf("[INFO] writePump started for client %s", c.clientId)
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		if c.hub != nil {
			log.Printf("[INFO] writePump closing for client %s", c.clientId)
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
				log.Printf("[INFO] send channel closed for client %s", c.clientId)
				return
			}
			c.mu.Lock()
			_, err := c.conn.Write(frame.ToBytes())
			c.mu.Unlock()

			if err != nil {
				fmt.Printf("[ERR] Erreur dans le select du writepump - message")
				return
			}
		case <-ticker.C:
			// Envoyer un ping
			c.mu.Lock()
			_, err := c.conn.Write(Ping([]byte("ping")))
			c.mu.Unlock()

			if err != nil {
				fmt.Printf("[ERR] Erreur dans le select du writepump - ticker")
				return
			}
		}
	}
}

func (c *Client) readPump() {
	log.Printf("[INFO] readPump started for client %s", c.clientId)
	defer func() {
		log.Printf("[INFO] readPump closing for client %s", c.clientId)
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

		n, err := c.conn.Read(temp)
		readCount++
		if err != nil {
			log.Printf("[ERR] <readPump> Error reading from client %s: %v", c.clientId, err)
			return
		}

		log.Printf("[INFO] Read %d bytes from client %s (read #%d)", n, c.clientId, readCount)
		if n == 0 {
			log.Printf("[WARN] Read 0 bytes from client %s, continuing...", c.clientId)
			continue
		}

		buffer = append(buffer, temp[:n]...)

		// Traiter toutes les frames complètes dans le buffer
		frameCount := 0
		for len(buffer) > 0 {
			frameCount++
			frame, frameLen, err := parseFrame(buffer)
			if err != nil {
				log.Printf("parseFrame error: %v", err)
				if err == ErrIncompleteFrame {
					// Frame incomplète, attendre plus de données
					log.Printf("[INFO] Incomplete frame, waiting for more data...")
					break
				}
				log.Printf("[ERR] Error parsing frame from client %s: %v", c.clientId, err)
				return
			}
			c.handleFrame(frame)
			buffer = buffer[frameLen:]
		}
	}
}
