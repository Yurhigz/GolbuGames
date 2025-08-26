package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

// Point d'entrée HTTP->WS
// Convertit HTTP en WebSocket (upgrade)
// Initialise les nouvelles connexions
// Valide les connexions entrantes

const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// Fonction de gestion de la clef secrete

func secretKeyVerification(clientKey string) (string, error) {
	if clientKey == "" {
		return "", errors.New("[ERR] missing Sec-WebSocket-Key")
	}
	hash := sha1.Sum([]byte(clientKey + magicGUID))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

func UpgradeConnection(w http.ResponseWriter, r *http.Request) (net.Conn, error) {

	if !strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") ||
		strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		return nil, errors.New("Invalid upgrade request")
	}

	// Vérification version
	if r.Header.Get("Sec-WebSocket-Version") != "13" {
		return nil, errors.New("unsupported websocket version")
	}

	acceptKey, err := secretKeyVerification(r.Header.Get("Sec-WebSocket-Key"))
	if err != nil {
		return nil, err
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("webserver doesn't support hijacking")
	}
	conn, buf, err := hj.Hijack()
	if err != nil {
		return nil, err
	}

	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n\r\n",
		acceptKey,
	)

	if _, err := buf.WriteString(response); err != nil {
		conn.Close()
		return nil, err
	}
	if err := buf.Flush(); err != nil {
		conn.Close()
		return nil, err
	}

	log.Printf("=== HANDSHAKE COMPLETE ===")
	return conn, nil
}

// func HandleWebSocketConnection(conn net.Conn, incoming chan<- Frame) {
// 	defer conn.Close()
// 	tmpbuf := make([]byte, 1024)
// 	var buf []byte

// 	for {
// 		n, err := conn.Read(tmpbuf)
// 		if err != nil {
// 			log.Println("[ERR] <handleWebSocketConnection> Error reading:", err)
// 			break
// 		}
// 		buf = append(buf, tmpbuf[:n]...)
// 		for {
// 			frame, frameLength, err := ParseFrame(buf)
// 			if err != nil {
// 				if errors.Is(err, ErrIncompleteFrame) {
// 					break // attendre plus de données
// 				}
// 				log.Println("[ERR] <handleWebSocketConnection> parse error:", err)
// 				return
// 			}
// 			incoming <- frame
// 			buf = buf[frameLength:]

// 		}

// 	}

// }
