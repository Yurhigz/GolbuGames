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
