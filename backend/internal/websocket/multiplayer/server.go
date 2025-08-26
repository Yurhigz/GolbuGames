package multiplayer

import (
	"fmt"
	"golbugames/backend/internal/websocket"
	"log"
	"net/http"
	"time"
)

// WebsocketHandler multijoueurs
func WebsocketHandler(w http.ResponseWriter, r *http.Request, hubManager *HubManager) {

	conn, err := websocket.UpgradeConnection(w, r)

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
