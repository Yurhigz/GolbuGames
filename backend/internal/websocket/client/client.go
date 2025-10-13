package client

import (
	"context"
	"fmt"
	"golbugames/backend/internal/websocket"
	"net"
	"sync"
	"time"
)

const (
	PongWait       = 60 * time.Second // Durée d'attente pour un pong
	Newline        = "\n"
	Space          = " "
	PingPeriod     = (PongWait * 9) / 10 // Période de ping pour garder la connexion active
	PongTimeout    = 60 * time.Second    // Durée d'attente pour un pong avant de fermer la connexion
	MaxMessageSize = 1024 * 1024
)

type BaseClient struct {
	ClientId      string
	Pseudo        string
	Elo           int
	Conn          net.Conn
	Mu            sync.Mutex
	Send          chan *websocket.Frame
	Solution      []int
	Playable      []int
	FrameBuffer   []byte
	CurrentOpcode byte
	Ctx           context.Context
}

func NewBaseClient(conn net.Conn) *BaseClient {
	return &BaseClient{
		Conn:     conn,
		Send:     make(chan *websocket.Frame, 256),
		ClientId: createId(),
	}
}

func createId() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

func (c *BaseClient) ResetFragmentation() {
	c.FrameBuffer = nil
	c.CurrentOpcode = 0
}

func (c *BaseClient) ValidateMove(index int, value byte) bool {
	// Vérifie si la valeur correspond à la solution
	if c.Solution[index] == int(value) {
		return true
	}
	return false
}

func (c *BaseClient) DuplicateFrame(frame *websocket.Frame) *websocket.Frame {
	f := *frame
	f.Payload = make([]byte, len(frame.Payload))
	copy(f.Payload, frame.Payload)
	return &f
}
