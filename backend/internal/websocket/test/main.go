package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

func main() {

	http.HandleFunc("/websocket", websocketHandler)
	fmt.Printf("Listening on port %v ...", 3005)

	http.ListenAndServe(":3005", nil)
	time.Sleep(100 * time.Millisecond)

	// select {}
	// b := byte(0b00001011)

	// // fmt.Println(isFinalFrame(b))
	// // fmt.Println(getOpcode(b))
	// // fmt.Println(hasAnyRSVSet(b))
	// // fmt.Println(buildControlByte(true, false, false, false, 0x1))
}

func isFinalFrame(b byte) bool {
	if b&(1<<7) != 0 {
		return true
	}
	return false
}

func isMaskSet(b byte) bool {
	return b&(1<<7) != 0
}

func getOpcode(b byte) byte {
	return b & 0b00001111
}

func hasAnyRSVSet(b byte) bool {
	if (b&0b01110000)>>4 > 0 {
		return true
	}
	return false
}

// func decode() {

// }

func buildControlByte(fin bool, rsv1, rsv2, rsv3 bool, opcode byte) byte {
	b := byte(0b00000000)
	if fin {
		b |= (1 << 7)
	}
	if rsv1 {
		b |= (1 << 6)
	}
	if rsv2 {
		b |= (1 << 5)

	}
	if rsv3 {
		b |= (1 << 4)
	}
	return b | opcode
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Header.Get("Connection")) != "upgrade" || strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		http.Error(w, "Invalid upgrade request", http.StatusBadRequest)
		return
	}

	secretKey := r.Header.Get("Sec-WebSocket-Key")

	if secretKey == "" {
		http.Error(w, "missing Sec-WebSocket-Key", http.StatusBadRequest)
		return
	}

	// un GUID est fixé par le RFC , il est immuable , preuve de conformité lors des handshakes et upgrade vers websocket

	const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.Sum([]byte(secretKey + magicGUID))
	acceptKey := base64.StdEncoding.EncodeToString(hash[:])

	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", acceptKey)
	w.WriteHeader(http.StatusSwitchingProtocols)

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	fmt.Println("WebSocket connection established!")
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		if isFinalFrame(buf[0]) {
			fmt.Println("This is the end of the message")
		}

		if !isMaskSet(buf[1]) {
			fmt.Println("Mask is not set")
		}
		//  Vider le buf utilisé : buf = buf[n:] mais attention allocation mémoire

		payloadLenIndicator := buf[1] & 0b01111111
		fmt.Println(int(payloadLenIndicator))
		switch {
		case payloadLenIndicator < 125:
			payloadLen := int(payloadLenIndicator)
			if len(buf) < 6+payloadLen {
				continue
			}
			maskKey := buf[2:6]
			payload := buf[6 : payloadLen+6]
			for i := 0; i < len(payload); i++ {
				payload[i] ^= maskKey[i%4]
			}
			fmt.Println(string(payload))

		case payloadLenIndicator == 126:
			if len(buf) < 8 {
				continue
			}
			payloadLen := int(binary.BigEndian.Uint16(buf[2:4]))
			if len(buf) < 8+payloadLen {
				continue
			}

			maskKey := buf[4:8]
			payload := buf[8 : 8+payloadLen]
			for i := 0; i < len(payload); i++ {
				payload[i] ^= maskKey[i%4]
			}
			fmt.Println(string(payload))

		case payloadLenIndicator == 127:
			if len(buf) < 14 {
				continue
			}
			payloadLen64 := binary.BigEndian.Uint64(buf[2:10])
			if payloadLen64 > math.MaxInt32 {
				continue
			}
			payloadLen := int(payloadLen64)
			if len(buf) < 14+payloadLen {
				continue
			}
			maskKey := buf[10:14]
			payload := buf[14 : 14+payloadLen]
			for i := 0; i < len(payload); i++ {
				payload[i] ^= maskKey[i%4]
			}
			fmt.Println(string(payload))
		}

	}

}
