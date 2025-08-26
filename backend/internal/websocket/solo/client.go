package solo

import (
	"encoding/json"
	"fmt"
	"golbugames/backend/internal/websocket"
	"io"
	"log"
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
	send        chan *websocket.Frame
	solution    []int
	playable    []int
}

func newSoloClient(conn net.Conn) *SoloClient {
	return &SoloClient{
		conn: conn,
		send: make(chan *websocket.Frame, 256),
	}
}

func (c *SoloClient) ValidateMove(index int, value byte) bool {
	// Vérifie si la valeur correspond à la solution
	if c.solution[index] == int(value) {
		return true
	}
	return false
}

func (c *SoloClient) handleFrame(frame websocket.Frame) {
	switch frame.Opcode {
	case websocket.OpcodeClose:
		log.Printf("[INFO] Client %s closed the connection", c.clientId)
		log.Printf("[INFO] Fermeture du client")
		c.send <- websocket.CloseFrame(1000, "Normal Closure")
		return

	case websocket.OpcodePing:
		log.Printf("[INFO] Received ping from client %s", c.clientId)
		pongFrame := websocket.Pong(frame.Payload)
		c.mu.Lock()
		_, err := c.conn.Write(pongFrame)
		c.mu.Unlock()
		if err != nil {
			log.Printf("[ERR] Error sending pong to client %s: %v", c.clientId, err)
			return
		}

	case websocket.OpcodePong:
		log.Printf("[INFO] Received pong from client %s", c.clientId)

	case websocket.OpcodeText, websocket.OpcodeBinary:

		var move struct {
			Position int  `json:"position"`
			Value    byte `json:"value"`
		}
		if err := json.Unmarshal(frame.Payload, &move); err != nil {
			log.Printf("[ERR] Invalid JSON: %v", err)
			return
		}

		valid := c.ValidateMove(move.Position, move.Value)

		resp := &websocket.Frame{
			Opcode: websocket.OpcodeText,
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

			if frame.Opcode == websocket.OpcodeClose {
				log.Printf("[INFO] Sent close frame to client %s, closing writePump", c.clientId)
				return
			}
		case <-ticker.C:
			c.mu.Lock()
			_, err := c.conn.Write(websocket.Ping([]byte("ping")))
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
			frame, framelen, err := websocket.ParseFrame(buffer)
			if err != nil {
				log.Printf("parseFrame error: %v", err)
				if err == websocket.ErrIncompleteFrame {
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
