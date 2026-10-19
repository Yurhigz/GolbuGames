package protocol

import (
	"time"
)

const (
	MessageTypeValidateMove = "validate_move"
	MessageTypeChat         = "chat"
	MessageTypeSystem       = "system"
	MessageTypeError        = "error"
)

// Messages entrants (Frontend -> backend)
const (
	InboundChat         = "chat"
	InboundGameMove     = "game_move"
	InboundVictoryClaim = "victory_claim"
)

// Messages sortants (Backend -> Frontend)
const (
	OutboundGameStart     = "gameStarting"
	OutboundWaiting       = "waitingOpponent"
	OutboundGameEnd       = "gameEnding"
	OutboundMoveResponse  = "move_response"
	OutboundChatBroadcast = "chat_broadcast"
	OutboundOpponentLeft  = "opponentLeft"
	OutboundCountdown     = 5
)

const (
	SystemGameStart        = 2000
	SystemGameOver         = 2001
	PlayerJoined           = 1000
	PlayerLeft             = 1001
	WaitingOpponent        = 1002
	PlayerDisconnected     = 1003
	MessageGameStart       = "Opponent found... Game will start"
	MessageWaitingOpponent = OutboundWaiting
)

// ===== MESSAGES ENTRANTS =====

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

type VictoryClaim struct {
	Type           string `json:"type"`
	Action         string `json:"action"`
	Message        string `json:"message"`
	CompletionTime int    `json:"completiontime"`
	Code           int    `json:"code"`
}

// ===== MESSAGES SORTANTS =====

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

type MoveValidationResponse struct {
	Type     string `json:"type"`
	Position int    `json:"position"`
	Value    int    `json:"value"`
	Valid    bool   `json:"valid"`
}

type GameStartNotification struct {
	Type      string `json:"type"`
	Grid      []int  `json:"grid"`
	Message   string `json:"message"`
	Countdown int    `json:"countdown"`
}

type GameEndNotification struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Winner  string `json:"winner"`
	Reason  string `json:"reason"`
}
