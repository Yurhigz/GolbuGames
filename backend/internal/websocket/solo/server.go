package solo

import (
	"encoding/json"
	"fmt"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/internal/websocket"
	"log"
	"net/http"
)

// WebsocketHandler Solo
func WebsocketHandlerSolo(w http.ResponseWriter, r *http.Request) {
	difficulty := r.URL.Query().Get("difficulty")

	if difficulty == "" {
		difficulty = "easy"
	}

	conn, err := websocket.UpgradeConnection(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur d'upgrade WebSocket: %v", err), http.StatusInternalServerError)
		log.Printf("[ERR] <websocket handler> Erreur d'upgrade connection")
		return
	}

	log.Println("[INFO] WebSocket connection established!")

	sudokuGrid, err := repository.GetRandomGridDB(r.Context(), difficulty)
	if err != nil {
		log.Printf("[ERR] Impossible to retrieve the grid: %v", err)
		conn.Close()
		return
	}

	client := newSoloClient(conn)
	client.playable = sudokuGrid.Board
	client.solution = sudokuGrid.Solution
	log.Printf("[INFO] New client created")

	go client.writePump()
	go client.readPump()

	payload, _ := json.Marshal(map[string]interface{}{
		"type":       "init",
		"grid":       client.playable,
		"difficulty": difficulty,
	})

	grid := &websocket.Frame{
		Opcode:  websocket.OpcodeText,
		FIN:     true,
		Payload: payload,
	}

	client.send <- grid

}
