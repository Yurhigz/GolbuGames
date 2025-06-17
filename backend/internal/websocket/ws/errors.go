package ws

import "errors"

// Frame errors
var (
	ErrIncompleteFrame     = errors.New("ws: incomplete frame received")
	ErrMissingMask         = errors.New("ws: client frame is not masked as required by RFC6455")
	ErrInvalidOpcode       = errors.New("ws: invalid or unsupported opcode")
	ErrPayloadTooLarge     = errors.New("ws: payload exceeds allowed size")
	ErrUnmaskedServerFrame = errors.New("ws: server frame must not be masked")
)

// Handshake errors
var (
	ErrInvalidUpgradeHeader    = errors.New("ws: invalid or missing Upgrade header")
	ErrInvalidConnectionHeader = errors.New("ws: invalid or missing Connection header")
	ErrMissingWebSocketKey     = errors.New("ws: missing Sec-WebSocket-Key header")
	ErrUnsupportedVersion      = errors.New("ws: unsupported WebSocket version")
	ErrHandshakeFailed         = errors.New("ws: WebSocket handshake failed")
)

// Internal errors
var (
	ErrConnectionClosed = errors.New("ws: connection already closed")
	ErrWriteFailed      = errors.New("ws: failed to write to connection")
)

// Client errors
var (
	ErrClientNotRegistered             = errors.New("ws: client not registered in the hub")
	ErrClientAlreadyExists             = errors.New("ws: client already exists in the hub")
	ErrClientNotFound                  = errors.New("ws: client not found in the hub")
	ErrClientDisconnected              = errors.New("ws: client disconnected unexpectedly")
	ErrClientTimeout                   = errors.New("ws: client connection timed out")
	ErrClientSendFailed                = errors.New("ws: failed to send message to client")
	ErrClientReceiveFailed             = errors.New("ws: failed to receive message from client")
	ErrClientInvalidMessage            = errors.New("ws: received invalid message from client")
	ErrClientMatchNotFound             = errors.New("ws: match not found for client")
	ErrClientAlreadyInMatch            = errors.New("ws: client already in a match")
	ErrClientMatchFull                 = errors.New("ws: match is full, cannot join")
	ErrClientMatchInProgress           = errors.New("ws: match is already in progress")
	ErrClientMatchEnded                = errors.New("ws: match has already ended")
	ErrClientMatchNotReady             = errors.New("ws: match is not ready to start")
	ErrClientMatchOpponentDisconnected = errors.New("ws: opponent has disconnected")
)
