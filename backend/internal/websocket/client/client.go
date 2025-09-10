package client

import (
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
}

func NewBaseClient(conn net.Conn) *BaseClient {
	return &BaseClient{
		Conn: conn,
		Send: make(chan *websocket.Frame, 256),
	}
}

func (c *BaseClient) ResetFragmentation() {
	c.FrameBuffer = nil
	c.CurrentOpcode = 0
}
