package protocol

import (
	"encoding/json"
	"time"
)

const (
	MessageTypeValidateMove = "validate_move"
	MessageTypeChat         = "chat"
	MessageTypeSystem       = "system"
	MessageTypeError        = "error"
)

const (
	SystemGameStart = 1000
	SystemGameOver  = 1001
	PlayerJoined    = 1002
	PlayerLeft      = 1003
	InternalError   = 3000
	ConnectionLost  = 3001
)

type TypeOnly struct {
	Type string `json:"type"`
}

type ChatMessage struct {
	Type    string    `json:"type"`
	Sender  string    `json:"sender"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sentat"`
}

type GameMessage struct {
	Type     string `json:"type"`
	Position int    `json:"position"`
	Value    int    `json:"value"`
}

type ValidationMove struct {
	Type     string `json:"type"`
	Position int    `json:"position"`
	Value    int    `json:"value"`
	Valid    bool   `json:"valid"`
}

type SystemMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewSystemMessage(message string, code int) ([]byte, error) {
	msg := SystemMessage{
		Type:    MessageTypeSystem,
		Message: message,
		Code:    code,
	}
	marshallized, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return marshallized, nil
}
