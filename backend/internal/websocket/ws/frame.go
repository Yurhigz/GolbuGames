package ws

// Ensemble des fonctions de lecture, écriture et parsing des frames websockets RFC6455

type Frame struct {
	FIN      bool
	Opcode   byte
	Masked   bool
	Mask     [4]byte
	Payload  []byte
	Length   int
	RawFrame []byte
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

func (f Frame) unmaskPayload() {
	for i := 0; i < len(f.Payload); i++ {
		f.Payload[i] ^= f.Mask[i%4]
	}
}

// Fonction de décodage des messages clients
func parseFrame(buf []byte) (Frame, int, error) {
	bufLength := len(buf)
	if bufLength < 2 {
		return Frame{}, 0, ErrIncompleteFrame
	}

	firstByte := buf[0]
	secondByte := buf[1]

	if !isMaskSet(secondByte) {
		return Frame{}, 0, ErrMissingMask
	}

	payloadLen := payloadLen(secondByte)

	// Check if len buf > headerlen + payloadlen (incorporer fonction pour déterminer au préalable longueur du payload)

	switch {
	case payloadLen < 126:
		if bufLength >= smallPayloadHeaderLen {
			mask := buf[smallIndicatorHeader:smallPayloadHeaderLen]
			payload := buf[smallPayloadHeaderLen : smallPayloadHeaderLen+payloadLen]
		}
		return Frame{}, 0, ErrIncompleteFrame

	case payloadLen == 126:
		if bufLength >= mediumPayloadHeaderLen {
			mask := buf[mediumIndicatorHeader:mediumPayloadHeaderLen]
			payload := buf[mediumPayloadHeaderLen : mediumPayloadHeaderLen+payloadLen]
		}
		return Frame{}, 0, ErrIncompleteFrame

	case payloadLen > 126:
		if bufLength >= largePayloadHeaderLen {
			mask := buf[largeIndicatorHeader:largePayloadHeaderLen]
			payload := buf[largePayloadHeaderLen : largePayloadHeaderLen+payloadLen]
		}

		return Frame{}, 0, ErrIncompleteFrame
	}

	frame := Frame{
		FIN:    isFinalFrame(firstByte),
		Opcode: getOpcode(firstByte),
		Masked: isMaskSet(secondByte),
	}

	return frame, frame.Length, nil
}

//  Fonction de construction des réponses côté serveur vers les clients
func BuildFrame(payload []byte, opcode byte, fin bool) []byte {

}
