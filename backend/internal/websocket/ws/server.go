package ws

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const magicGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

// Fonction de gestion de la clef secrete

func secretKeyVerification(clientKey string) (string, error) {
	if clientKey == "" {
		return "", errors.New("missing Sec-WebSocket-Key")
	}
	hash := sha1.Sum([]byte(clientKey + magicGUID))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// Fonction de coordination/service
func WebsocketHandler(w http.ResponseWriter, r *http.Request, hubManager *HubManager) {

	conn, err := upgradeConnection(w, r)

	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur d'upgrade WebSocket: %v", err), http.StatusInternalServerError)
		log.Printf("[ERR] <websocket handler> Erreur d'upgrade connection")
		return
	}

	log.Println("[INFO] WebSocket connection established!")

	client := newClient(conn, hubManager)
	client.queueTime = time.Now()
	go client.writePump()
	go client.readPump()
	hubManager.mu.Lock()
	hubManager.ClientQueue = append(hubManager.ClientQueue, client)
	hubManager.mu.Unlock()

	log.Printf("[INFO] Nouveau client connecté au hubManager")

}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (net.Conn, error) {

	if strings.ToLower(r.Header.Get("Connection")) != "upgrade" || strings.ToLower(r.Header.Get("Upgrade")) != "websocket" {
		return nil, errors.New("Invalid upgrade request")
	}

	acceptKey, err := secretKeyVerification(r.Header.Get("Sec-WebSocket-Key"))

	if err != nil {
		return nil, err
	}
	w.Header().Set("Sec-WebSocket-Version", "13")
	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	// w.Header().Set("Sec-WebSocket-Extensions", "")
	w.Header().Set("Sec-WebSocket-Accept", acceptKey)

	w.WriteHeader(http.StatusSwitchingProtocols)

	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("webserver doesn't support hijacking")
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		return nil, err
	}

	log.Printf("=== HANDSHAKE COMPLETE ===")
	return conn, nil
}

func handleWebSocketConnection(conn net.Conn, incoming chan<- Frame) {
	defer conn.Close()
	tmpbuf := make([]byte, 1024)
	var buf []byte

	for {
		n, err := conn.Read(tmpbuf)
		if err != nil {
			log.Println("[ERR] <handleWebSocketConnection> Error reading:", err)
			break
		}
		buf = append(buf, tmpbuf[:n]...)
		for {
			frame, frameLength, err := parseFrame(buf)
			if err != nil {
				if errors.Is(err, ErrIncompleteFrame) {
					break // attendre plus de données
				}
				log.Println("[ERR] <handleWebSocketConnection> parse error:", err)
				return
			}
			incoming <- frame
			buf = buf[frameLength:]

		}

	}

}
