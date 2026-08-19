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
	MaxMessageSize = 1024 * 1024
)

type BaseClient struct {
	ClientId      string
	Pseudo        string
	Elo           int
	Conn          net.Conn
	Mu            sync.Mutex
	Send          chan *websocket.Frame
	Done		  chan struct{}
	closeOnce	 sync.Once
	Solution      []int
	Playable      []int
	FrameBuffer   []byte
	CurrentOpcode byte
}

func NewBaseClient(conn net.Conn) *BaseClient {
	return &BaseClient{
		Conn:     conn,
		Send:     make(chan *websocket.Frame, 256),
		Done:     make(chan struct{}),
		ClientId: createId(),
	}
}

func createId() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}

func (c *BaseClient) Close() {
	c.closeOnce.Do(func() {
		close(c.Done)
		c.Conn.Close()
	})
}

// TrySend tente d'envoyer une frame sans jamais bloquer indéfiniment :
// si Done est fermé (writePump arrêté / connexion en cours de fermeture),
// on abandonne au lieu de deadlock.
func (c *BaseClient) TrySend(frame *websocket.Frame) bool {
	select {
	case c.Send <- frame:
		return true
	case <-c.Done:
		return false
	}
}

func (c *BaseClient) ResetFragmentation() {
	c.FrameBuffer = nil
	c.CurrentOpcode = 0
}

func (c *BaseClient) ValidateMove(index int, value byte) bool {
	if index < 0 || index >= len(c.Solution) {
		return false
	}
	return c.Solution[index] == int(value)
}

func (c *BaseClient) DuplicateFrame(frame *websocket.Frame) *websocket.Frame {
	f := *frame
	f.Payload = make([]byte, len(frame.Payload))
	copy(f.Payload, frame.Payload)
	return &f
}
