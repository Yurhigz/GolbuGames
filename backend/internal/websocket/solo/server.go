package solo

import (
	"fmt"
	"golbugames/backend/internal/websocket"
	"log"
	"net/http"
)

// WebsocketHandler Solo
func WebsocketHandlerSolo(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.UpgradeConnection(w, r)
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
