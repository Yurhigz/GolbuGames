package solo

import (
	"encoding/json"
	"fmt"
	"golbugames/backend/internal/websocket"
	"golbugames/backend/internal/websocket/client"
	"io"
	"log"
	"net"
	"time"
)

// structure client pour solo

type SoloClient struct {
	baseClient *client.BaseClient
}

func newSoloClient(conn net.Conn) *SoloClient {
	return &SoloClient{
		baseClient: client.NewBaseClient(conn),
	}
}

func (c *SoloClient) handleFrame(frame websocket.Frame) {
	switch frame.Opcode {
	case websocket.OpcodeClose:
		log.Printf("[INFO] Client %s closed the connection", c.baseClient.ClientId)
		log.Printf("[INFO] Fermeture du client")
		c.baseClient.Send <- websocket.CloseFrame(1000, "Normal Closure")
		return

	case websocket.OpcodePing:
		log.Printf("[INFO] Received ping from client %s", c.baseClient.ClientId)
		pongFrame := websocket.Pong(frame.Payload)
		c.baseClient.Mu.Lock()
		_, err := c.baseClient.Conn.Write(pongFrame)
		c.baseClient.Mu.Unlock()
		if err != nil {
			log.Printf("[ERR] Error sending pong to client %s: %v", c.baseClient.ClientId, err)
			return
		}

	case websocket.OpcodePong:
		log.Printf("[INFO] Received pong from client %s", c.baseClient.ClientId)

	case websocket.OpcodeText, websocket.OpcodeBinary:

		var move struct {
			Position int  `json:"position"`
			Value    byte `json:"value"`
		}
		if err := json.Unmarshal(frame.Payload, &move); err != nil {
			log.Printf("[ERR] Invalid JSON: %v", err)
			return
		}

		valid := c.baseClient.ValidateMove(move.Position, move.Value)

		resp := &websocket.Frame{
			Opcode: websocket.OpcodeText,
			FIN:    true,
			Payload: []byte(fmt.Sprintf(
				`{"position":%d,"value":%d,"valid":%t}`,
				move.Position, move.Value, valid,
			)),
		}

		c.baseClient.Send <- resp

	default:
		log.Printf("[INFO] Received unknown frame type (0x%02x) from client %s", frame.Opcode, c.baseClient.ClientId)
	}
}

func (c *SoloClient) writePump() {
	log.Printf("[INFO] writePump started for client %s", c.baseClient.ClientId)
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.baseClient.Conn.Close()
	}()

	for {
		select {
		case frame, ok := <-c.baseClient.Send:
			if !ok {
				log.Printf("[INFO] send channel closed for client %s", c.baseClient.ClientId)
				return
			}
			c.baseClient.Mu.Lock()
			_, err := c.baseClient.Conn.Write(frame.ToBytes())
			c.baseClient.Mu.Unlock()
			if err != nil {
				log.Printf("[ERR] Error writing to client %s: %v", c.baseClient.ClientId, err)
				return
			}

			if frame.Opcode == websocket.OpcodeClose {
				log.Printf("[INFO] Sent close frame to client %s, closing writePump", c.baseClient.ClientId)
				return
			}
		case <-ticker.C:
			c.baseClient.Mu.Lock()
			_, err := c.baseClient.Conn.Write(websocket.Ping([]byte("ping")))
			c.baseClient.Mu.Unlock()

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
		log.Printf("[INFO] Closing connection for client %s", c.baseClient.ClientId)
		c.baseClient.Conn.Close()
		close(c.baseClient.Send)
	}()

	for {
		temp := make([]byte, 1024)
		n, err := c.baseClient.Conn.Read(temp)
		if err != nil {
			if err == io.EOF {
				log.Printf("[INFO] Client %s disconnected", c.baseClient.ClientId)
			} else {
				log.Printf("[ERR] <readPump> Error reading from client %s: %v", c.baseClient.ClientId, err)
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
				log.Printf("[ERR] Error parsing frame from client %s: %v", c.baseClient.ClientId, err)
				return
			}

			buffer = buffer[framelen:]
			c.handleFrame(frame)
		}
	}
}
