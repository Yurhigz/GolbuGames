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
		return "", errors.New("[ERR] missing Sec-WebSocket-Key")
	}
	hash := sha1.Sum([]byte(clientKey + magicGUID))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// WebsocketHandler Solo
func WebsocketHandlerSolo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgradeConnection(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur d'upgrade WebSocket: %v", err), http.StatusInternalServerError)
		log.Printf("[ERR] <websocket handler> Erreur d'upgrade connection")
		return
	}

	log.Println("[INFO] WebSocket connection established!")

	client := newSoloClient(conn)

	go client.writePump()
	go client.readPump()

	log.Printf("[INFO] Nouveau client crée")

}

// WebsocketHandler multijoueurs
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
