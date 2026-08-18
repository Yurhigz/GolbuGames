package websocket

import (
	"encoding/binary"
	"time"
)

type Frame struct {
	FIN      bool
	Opcode   byte
	Masked   bool
	Mask     [4]byte
	Payload  []byte
	Length   int
	RawFrame []byte
}

type BroadcastFrame struct {
	Frame  *Frame
	Sender string
	SentAt time.Time
}

const (
	OpcodeContinuation = 0x0
	OpcodeText         = 0x1
	OpcodeBinary       = 0x2
	OpcodeClose        = 0x8
	OpcodePing         = 0x9
	OpcodePong         = 0xA
)

const (
	smallIndicatorHeader   = 2
	smallPayloadHeaderLen  = 6
	mediumIndicatorHeader  = 4
	mediumPayloadHeaderLen = 8
	largeIndicatorHeader   = 10
	largePayloadHeaderLen  = 14
	MaxPayloadSize         = 1024 * 1024
)

func isFinalFrame(b byte) bool {
	return b&(1<<7) != 0
}

func getOpcode(b byte) byte {
	return b & 0b00001111
}

func isControleFrame(opcode byte) bool {
	return opcode == OpcodePing || opcode == OpcodePong
}

func isCloseFrame(opcode byte) bool {
	return opcode == OpcodeClose
}

func isContinuationFrame(opcode byte) bool {
	return opcode == OpcodeContinuation
}

func isMaskSet(b byte) bool {
	return b&(1<<7) != 0
}

func hasAnyRSVSet(b byte) bool {
	return (b&0b01110000)>>4 > 0
}

func payloadLen(b byte) byte {
	return b & 0b01111111
}

func unmaskPayload(payload []byte, mask []byte) {
    for i := 0; i < len(payload); i++ {
        payload[i] ^= mask[i%4]
    }
}

func Ping(payload []byte) []byte {
	return BuildFrame(payload, OpcodePing, true)
}

func Pong(payload []byte) []byte {
	return BuildFrame(payload, OpcodePong, true)
}

func OpcodeToString(opcode byte) string {
	switch opcode {
	case OpcodeContinuation:
		return "Continuation"
	case OpcodeText:
		return "Text"
	case OpcodeBinary:
		return "Binary"
	case OpcodeClose:
		return "Close"
	case OpcodePing:
		return "Ping"
	case OpcodePong:
		return "Pong"
	default:
		return "Unknown"
	}
}

func extractPayload(buf []byte) []byte {
	if len(buf) < 2 {
		return nil
	}

	payloadLenIndicator := payloadLen(buf[1])
	var payloadStart int

	switch {
	case payloadLenIndicator < 126:
		payloadStart = 2 // 2 bytes for header
	case payloadLenIndicator == 126:
		payloadStart = 4 // 2 bytes for extended length
	case payloadLenIndicator > 126:
		payloadStart = 10 // 8 bytes for extended length
	}

	if len(buf) < payloadStart {
		return nil
	}

	return buf[payloadStart:]
}

// Fonction de décodage des messages clients
func ParseFrame(buf []byte) (Frame, int, error) {
	bufLength := len(buf)
	frame := Frame{}
	if bufLength < 2 {
		return frame, 0, ErrIncompleteFrame
	}

	firstByte := buf[0]
	secondByte := buf[1]

	frame.FIN = isFinalFrame(firstByte)
	frame.Opcode = getOpcode(firstByte)
	frame.Masked = isMaskSet(secondByte)

	if !isMaskSet(secondByte) {
		return frame, 0, ErrMissingMask
	}

	payloadLenIndicator := payloadLen(secondByte)
	var payloadLen int
	var headerLen int
	// Check if len buf > headerlen + payloadlen (incorporer fonction pour déterminer au préalable longueur du payload)

	switch {
	case payloadLenIndicator < 126:
		payloadLen = int(payloadLenIndicator)
		headerLen = 6 // 2 (base) + 4 (mask)

	case payloadLenIndicator == 126:
		if len(buf) < 4 {
			return frame, 0, ErrIncompleteFrame
		}
		payloadLen = int(binary.BigEndian.Uint16(buf[2:4]))
		headerLen = 8 // 2 (base) + 2 (extended) + 4 (mask)

	case payloadLenIndicator > 126:
		if len(buf) < 10 {
			return frame, 0, ErrIncompleteFrame
		}
		payloadLen64 := binary.BigEndian.Uint64(buf[2:10])
		if payloadLen64 > uint64(MaxPayloadSize) {
			return frame, 0, ErrPayloadTooLarge
		}
		payloadLen = int(payloadLen64)
		headerLen = 14 // 2 + 8 + 4
	}

	totalLen := headerLen + payloadLen
	if len(buf) < totalLen {
		return frame, 0, ErrIncompleteFrame
	}

	maskStart := headerLen - 4
	mask := buf[maskStart:headerLen]
	frame.Mask = [4]byte{mask[0], mask[1], mask[2], mask[3]}

	frame.Payload = make([]byte, payloadLen)
	copy(frame.Payload, buf[headerLen:totalLen])

	unmaskPayload(frame.Payload, mask)

	frame.Length = totalLen
	frame.RawFrame = make([]byte, totalLen)
	copy(frame.RawFrame, buf[:totalLen])

	return frame, frame.Length, nil
}

// Fonction de construction des réponses côté serveur vers les clients
func BuildFrame(payload []byte, opcode byte, fin bool) []byte {
	payloadLength := len(payload)
    var headerLen int

	switch {
    case payloadLength < 126:
        headerLen = 2
    case payloadLength <= 65535:
        headerLen = 4
    default:
        headerLen = 10
    }

	frame := make([]byte, headerLen+payloadLength)

	// 1er octet : FIN + Opcode
    var firstByte byte = 0
    if fin {
        firstByte |= 1 << 7
    }
    firstByte |= opcode
    frame[0] = firstByte

	// 2ème octet + longueur étendue
    switch {
    case payloadLength < 126:
        frame[1] = byte(payloadLength)
    case payloadLength <= 65535:
        frame[1] = 126
        binary.BigEndian.PutUint16(frame[2:4], uint16(payloadLength))
    default:
        frame[1] = 127
        binary.BigEndian.PutUint64(frame[2:10], uint64(payloadLength))
    }

    copy(frame[headerLen:], payload)

    return frame

}

func (f *Frame) ToBytes() []byte {
	return BuildFrame(f.Payload, f.Opcode, f.FIN)
}

func CloseFrame(code uint16, reason string) *Frame {
	payload := make([]byte, 2+len(reason))
	binary.BigEndian.PutUint16(payload, code)
	copy(payload[2:], reason)
	return &Frame{
		FIN:     true,
		Opcode:  OpcodeClose,
		Payload: payload,
	}
}
