package multiplayer

import (
	"encoding/json"
	"golbugames/internal/websocket"
	"golbugames/internal/websocket/client"
	"golbugames/internal/websocket/protocol"
	"log"
	"net"
	"time"
)

type Client struct {
	baseClient *client.BaseClient
	hub        *Hub
	hubManager *HubManager
	queueTime  time.Time
	userID     string
	username   string
	closed     bool
}

func (c *Client) setHub(hub *Hub) bool {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()
	if c.closed {
		return false
	}
	c.hub = hub
	return true
}

func (c *Client) getHub() *Hub {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()
	return c.hub
}

func (c *Client) isClosed() bool {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()
	return c.closed
}

func newClient(conn net.Conn, hubManager *HubManager, userID, username string) *Client {
	return &Client{
		baseClient: client.NewBaseClient(conn),
		hubManager: hubManager,
		userID:     userID,
		username:   username,
		closed:     false,
	}
}

func (c *Client) cleanup() {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()

	if c.closed {
		return
	}
	c.closed = true

	log.Printf("Cleaning up client %s", c.baseClient.ClientId)

	c.hubManager.RemoveClientFromQueue(c)

	if c.hub != nil {
		select {
		case c.hub.unregister <- c:
		default:
			log.Printf("Hub unregister channel blocked for client %s", c.baseClient.ClientId)
		}
	}

	if c.baseClient != nil && c.baseClient.Conn != nil {
		c.baseClient.Conn.Close()
	}

	if c.baseClient != nil && c.baseClient.Send != nil {
		close(c.baseClient.Send)
	}
}

func (c *Client) SendGameStart() {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()
	if c.closed {
		return
	}

	c.baseClient.Playable = c.hub.playable
	c.baseClient.Solution = c.hub.solution

	msg := protocol.GameStartNotification{
		Type:      protocol.MessageTypeSystem,
		Grid:      c.baseClient.Playable,
		Message:   protocol.OutboundGameStart,
		Countdown: protocol.OutboundCountdown,
	}

	marshallized, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marhalizing : %v", err)
		return
	}

	resp := &websocket.Frame{
		Opcode:  websocket.OpcodeText,
		FIN:     true,
		Payload: marshallized,
	}

	select {
	case c.baseClient.Send <- resp:
		log.Printf("Opponent found confirmation message sent to client %s", c.baseClient.ClientId)

	default:
		log.Printf("Failed to send confirmation message to client %s", c.baseClient.ClientId)

	}

}

func (c *Client) SendWaitingOpponent() {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()
	if c.closed {
		return
	}

	msg := protocol.SystemMessage{
		Type:    protocol.MessageTypeSystem,
		Message: protocol.MessageWaitingOpponent,
		Code:    protocol.WaitingOpponent,
	}

	marshallized, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marhalizing : %v", err)
		return
	}

	resp := &websocket.Frame{
		Opcode:  websocket.OpcodeText,
		FIN:     true,
		Payload: marshallized,
	}

	select {
	case c.baseClient.Send <- resp:
		log.Printf("Waiting message sent to client %s", c.baseClient.ClientId)

	default:
		log.Printf("Failed to send waiting message to client %s", c.baseClient.ClientId)

	}
}

func (c *Client) SendLeavingOpponent() {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()

	if c.closed {
		return
	}
	payload, err := json.Marshal(protocol.SystemMessage{
		Type:    protocol.MessageTypeSystem,
		Message: protocol.OutboundOpponentLeft,
		Code:    protocol.PlayerDisconnected,
	})

	if err != nil {
		log.Printf("Error marhalizing : %v", err)
		return
	}

	resp := &websocket.Frame{
		FIN:     true,
		Opcode:  websocket.OpcodeText,
		Payload: payload,
	}

	select {
	case c.baseClient.Send <- resp:
		log.Printf("ending game sent to client %s", c.baseClient.ClientId)

	default:
		log.Printf("Failed to send ending message to client %s", c.baseClient.ClientId)

	}

}

func (c *Client) SendGameOver(winner string) {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()

	if c.closed {
		return
	}
	payload, err := json.Marshal(protocol.GameEndNotification{
		Type:    protocol.MessageTypeSystem,
		Message: protocol.OutboundGameEnd,
		Winner:  winner,
		Reason:  protocol.OutboundGameEnd,
	})

	if err != nil {
		log.Printf("Error marhalizing : %v", err)
		return
	}

	resp := &websocket.Frame{
		FIN:     true,
		Opcode:  websocket.OpcodeText,
		Payload: payload,
	}

	select {
	case c.baseClient.Send <- resp:
		log.Printf("ending game sent to client %s", c.baseClient.ClientId)

	default:
		log.Printf("Failed to send ending message to client %s", c.baseClient.ClientId)

	}

}

func (c *Client) SendMessage(payload []byte) {
	c.baseClient.Mu.Lock()
	defer c.baseClient.Mu.Unlock()

	if c.closed {
		return
	}

	frame := &websocket.Frame{
		FIN:     true,
		Opcode:  websocket.OpcodeText,
		Payload: payload,
	}

	select {
	case c.baseClient.Send <- frame:
		log.Printf("Message queued for client %s", c.baseClient.ClientId)
	default:
		log.Printf("Failed to queue message for client %s - channel full", c.baseClient.ClientId)
	}
}

func (c *Client) HandlerChat(msg protocol.ChatMessage) {

	hub := c.getHub()
	if hub == nil {
		log.Printf("Client %s not in hub, cannot broadcast chat", c.baseClient.ClientId)
		return
	}

	if len(msg.Message) == 0 || len(msg.Message) > 500 {
		log.Printf("Ignoring invalid chat message from %s", c.baseClient.ClientId)
		return
	}
	msg.Type = protocol.MessageTypeChat
	msg.Sender = c.username
	msg.SentAt = time.Now()
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error Marshalling msg : %v", err)
		return
	}

	frame := websocket.Frame{
		Opcode:  websocket.OpcodeText,
		FIN:     true,
		Payload: payload,
	}
	select {
	case hub.broadcast <- frame:
		log.Printf("Chat message broadcasted from %s", msg.Sender)
	default:
		log.Printf("Failed to broadcast chat - channel full or closed")
	}

}

func (c *Client) HandlerMove(msg protocol.GameMessage) {
	hub := c.getHub()
	if hub == nil {
		log.Printf("Ignoring move from client %s outside a game", c.baseClient.ClientId)
		return
	}
	hub.HandleMove(c, msg)
}

func (c *Client) HandlerSystemMessage(msg protocol.VictoryClaim) {
	if msg.Action == "quit" {
		c.cleanup()
		return
	}
	if msg.Message == protocol.InboundVictoryClaim {
		log.Printf("Ignoring client-supplied victory claim from %s", c.baseClient.ClientId)
		return
	}
	log.Printf("Unknown system message from client %s", c.baseClient.ClientId)
}

func (c *Client) ProcessMessage(payload []byte) {
	var TypeOnly protocol.TypeOnly
	err := json.Unmarshal(payload, &TypeOnly)
	if err != nil {
		log.Printf("Invalid message format from client %s: %v", c.baseClient.ClientId, err)
		return
	}

	switch TypeOnly.Type {
	case protocol.MessageTypeChat:
		var chatMessage protocol.ChatMessage
		if err := json.Unmarshal(payload, &chatMessage); err != nil {
			log.Printf("Invalid chat message from client %s: %v", c.baseClient.ClientId, err)
			return
		}
		c.HandlerChat(chatMessage)
	case protocol.MessageTypeValidateMove:
		var gameMessage protocol.GameMessage
		if err := json.Unmarshal(payload, &gameMessage); err != nil {
			log.Printf("Invalid move message from client %s: %v", c.baseClient.ClientId, err)
			return
		}
		c.HandlerMove(gameMessage)
	case protocol.MessageTypeSystem:
		var victoryClaim protocol.VictoryClaim
		if err := json.Unmarshal(payload, &victoryClaim); err != nil {
			log.Printf("Invalid system message from client %s: %v", c.baseClient.ClientId, err)
			return
		}
		c.HandlerSystemMessage(victoryClaim)
	default:
		log.Printf("Unknown message type from client %s: %s", c.baseClient.ClientId, TypeOnly.Type)
		return

	}
}

func (c *Client) handleFrame(frame websocket.Frame) {
	switch frame.Opcode {
	case websocket.OpcodeClose:
		log.Printf("Client %s closed the connection", c.baseClient.ClientId)
		c.cleanup()
		return

	case websocket.OpcodePing:
		log.Printf("Received ping from client %s", c.baseClient.ClientId)
		pongFrame := websocket.Pong(frame.Payload)
		c.baseClient.Mu.Lock()
		_, err := c.baseClient.Conn.Write(pongFrame)
		c.baseClient.Mu.Unlock()
		if err != nil {
			log.Printf("Error sending pong to client %s: %v", c.baseClient.ClientId, err)
			c.cleanup()
			return
		}

	case websocket.OpcodePong:
		log.Printf("Received pong from client %s", c.baseClient.ClientId)

	case websocket.OpcodeText, websocket.OpcodeBinary:
		c.baseClient.CurrentOpcode = frame.Opcode

		if frame.FIN {
			log.Printf("Received complete %s message from client %s", websocket.OpcodeToString(frame.Opcode), c.baseClient.ClientId)
			c.ProcessMessage(c.baseClient.FrameBuffer)
			c.baseClient.ResetFragmentation()
		} else {
			log.Printf("Received first frame of fragmented %s message from client %s",
				websocket.OpcodeToString(frame.Opcode), c.baseClient.ClientId)

			if uint64(len(frame.Payload)) > websocket.MaxPayloadSize() {
				log.Printf("First frame too large from client %s", c.baseClient.ClientId)
				c.baseClient.ResetFragmentation()
				return
			}

			c.baseClient.FrameBuffer = append(c.baseClient.FrameBuffer[:0], frame.Payload...)
		}

	case websocket.OpcodeContinuation:
		if c.baseClient.FrameBuffer == nil || c.baseClient.CurrentOpcode == 0 {
			log.Printf("Received continuation frame without initial frame from client %s", c.baseClient.ClientId)
			return
		}

		if uint64(len(c.baseClient.FrameBuffer)+len(frame.Payload)) > websocket.MaxPayloadSize() {
			log.Printf("Message too large from client %s", c.baseClient.ClientId)
			c.baseClient.ResetFragmentation()
			return
		}

		c.baseClient.FrameBuffer = append(c.baseClient.FrameBuffer, frame.Payload...)

		if frame.FIN {
			log.Printf("Received final continuation frame from client %s", c.baseClient.ClientId)
			c.ProcessMessage(frame.Payload)
			c.baseClient.ResetFragmentation()
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
		c.cleanup()
	}()

	for {
		select {
		case message, ok := <-c.baseClient.Send:
			if !ok {
				log.Printf("Send channel closed for client %s", c.baseClient.ClientId)
				return
			}

			frameBytes := message.ToBytes()

			c.baseClient.Mu.Lock()
			_, err := c.baseClient.Conn.Write(frameBytes)
			c.baseClient.Mu.Unlock()

			if err != nil {
				log.Printf("Error writing to client %s: %v", c.baseClient.ClientId, err)
				return
			}
		case <-ticker.C:
			c.baseClient.Mu.Lock()
			_, err := c.baseClient.Conn.Write(websocket.Ping([]byte("ping")))
			c.baseClient.Mu.Unlock()

			if err != nil {
				log.Printf("Error sending ping to client %s: %v", c.baseClient.ClientId, err)
				return
			}
		}
	}

}

func (c *Client) readPump() {

	defer c.cleanup()

	buffer := make([]byte, 0, 4096)

	for {

		c.baseClient.Mu.Lock()
		if c.closed {
			c.baseClient.Mu.Unlock()
			return
		}
		c.baseClient.Mu.Unlock()
		temp := make([]byte, 1024)

		if tcpConn, ok := c.baseClient.Conn.(*net.TCPConn); ok {
			tcpConn.SetReadDeadline(time.Now().Add(60 * time.Second))
		}

		n, err := c.baseClient.Conn.Read(temp)
		if err != nil {
			log.Printf("Error reading from client %s: %v", c.baseClient.ClientId, err)
			return
		}

		if n == 0 {
			continue
		}

		buffer = append(buffer, temp[:n]...)

		for len(buffer) > 0 {
			frame, frameLen, err := websocket.ParseFrame(buffer)
			if err != nil {
				if err == websocket.ErrIncompleteFrame {
					break
				}
				log.Printf("Error parsing frame from client %s: %v", c.baseClient.ClientId, err)
				return
			}
			c.handleFrame(frame)

			buffer = buffer[frameLen:]
		}
	}
}
