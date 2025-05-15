package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"strings"
)

func main() {

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
	hash := sha1.Sum([]byte(secKey + magicGUID))
	acceptKey := base64.StdEncoding.EncodeToString(hash[:])

}
