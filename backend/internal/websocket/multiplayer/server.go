package multiplayer

import (
	"fmt"
	"golbugames/internal/websocket"
	"log"
	"net/http"
	"time"

	"golbugames/internal/api/middleware"
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request, hubManager *HubManager, jwtManager *middleware.JWTManager) {
	claims, err := middleware.ClaimsFromWebSocketProtocols(r, jwtManager)
	if err != nil {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.UpgradeConnection(w, r, "golbugames")

	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur d'upgrade WebSocket: %v", err), http.StatusBadRequest)
		log.Printf("[ERR] <websocket handler> Erreur d'upgrade connection")
		return
	}

	log.Println("[INFO] WebSocket connection established!")

	client := newClient(conn, hubManager, claims.UserID, claims.Username)
	client.queueTime = time.Now()
	hubManager.mu.Lock()
	hubManager.ClientQueue = append(hubManager.ClientQueue, client)
	hubManager.mu.Unlock()
	go client.writePump()
	go client.readPump()

	log.Printf("[INFO] Nouveau client connecté au hubManager")

}
