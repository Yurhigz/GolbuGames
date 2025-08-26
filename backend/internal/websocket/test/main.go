package main

import (
	"fmt"
	"golbugames/backend/internal/websocket/multiplayer"
	"net/http"
	"time"
)

const (
	smallPayloadHeaderLen  = 6
	mediumPayloadHeaderLen = 8
	largePayloadHeaderLen  = 14
)

func main() {
	HubManager := multiplayer.NewHubManager()
	go HubManager.MatchmakingLoop()
	// A debug
	// go HubManager.HubCleanupLoop()
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		multiplayer.WebsocketHandler(w, r, HubManager)
	})
	fmt.Printf("Listening on port %v ...", 3005)

	http.ListenAndServe(":3005", nil)
	time.Sleep(100 * time.Millisecond)

}
