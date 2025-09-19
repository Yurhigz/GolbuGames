package main

import (
	"context"
	"fmt"
	"golbugames/backend/internal/websocket/multiplayer"
	"net/http"
	"time"
)

func main() {
	parentContext := context.Background()
	HubManager := multiplayer.NewHubManager()
	go HubManager.MatchmakingLoop(parentContext)

	// A debug
	// go HubManager.HubCleanupLoop()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		multiplayer.WebsocketHandler(w, r, HubManager)
	})
	fmt.Printf("Listening on port %v ...", 3005)

	http.ListenAndServe(":3005", nil)
	time.Sleep(100 * time.Millisecond)

}
