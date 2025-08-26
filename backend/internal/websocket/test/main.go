package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	smallPayloadHeaderLen  = 6
	mediumPayloadHeaderLen = 8
	largePayloadHeaderLen  = 14
)

func main() {
	HubManager := ws.NewHubManager()
	go HubManager.MatchmakingLoop()
	// A debug
	// go HubManager.HubCleanupLoop()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		ws.WebsocketHandler(w, r, HubManager)
	})
	fmt.Printf("Listening on port %v ...", 3005)

	http.ListenAndServe(":3005", nil)
	time.Sleep(100 * time.Millisecond)

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
	tmpbuf := make([]byte, 1024)
	var buf []byte
	for {
		n, err := conn.Read(tmpbuf)
		if err != nil {
			fmt.Println("Error reading:", err)
			break
		}
		buf = append(buf, tmpbuf[:n]...)

		//  Lire les deux premiers bytes du buf pour savoir s'il contient une frame complète

		payloadLenIndicator := buf[1] & 0b01111111

		switch {
		case payloadLenIndicator < 125:
			payloadLen := int(payloadLenIndicator)
			frameLen := int(smallPayloadHeaderLen + payloadLen)
			if len(buf) >= frameLen {
				maskKey := buf[2:6]
				payload := buf[smallPayloadHeaderLen : payloadLen+smallPayloadHeaderLen]
				unmaskedPayload := unmaskPayload(maskKey, payload)
				fmt.Println(string(unmaskedPayload))
				buf = buf[frameLen:]
			}
			continue

		case payloadLenIndicator == 126:
			if len(buf) < 8 {
				continue
			}
			payloadLen := int(binary.BigEndian.Uint16(buf[2:4]))
			frameLen := int(mediumPayloadHeaderLen + payloadLen)
			if len(buf) >= frameLen {
				maskKey := buf[4:8]
				payload := buf[mediumPayloadHeaderLen : mediumPayloadHeaderLen+payloadLen]
				unmaskedPayload := unmaskPayload(maskKey, payload)
				fmt.Println(string(unmaskedPayload))
				buf = buf[frameLen:]
			}

		case payloadLenIndicator == 127:
			if len(buf) < 14 {
				continue
			}
			payloadLen64 := int(binary.BigEndian.Uint64(buf[2:10]))
			frameLen := int(largePayloadHeaderLen + payloadLen64)
			if len(buf) >= frameLen {
				payloadLen := int(payloadLen64)
				maskKey := buf[10:14]
				payload := buf[largePayloadHeaderLen : largePayloadHeaderLen+payloadLen]
				unmaskedPayload := unmaskPayload(maskKey, payload)
				fmt.Println(string(unmaskedPayload))
				buf = buf[frameLen:]
			}
			continue

		}

	}

}
