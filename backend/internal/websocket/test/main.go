package main

import (
	"context"
	"fmt"
	"golbugames/backend/config"
	"golbugames/backend/internal/database"
	"golbugames/backend/internal/websocket/multiplayer"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	parentContext := context.Background()
	HubManager := multiplayer.NewHubManager()
	go HubManager.MatchmakingLoop(parentContext)

	// A debug
	// Init DB
	err := database.InitDB(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Init Sudoku Grids

	err = config.InitGridGeneration(ctx)
	if err != nil {
		log.Fatalf("Failed to generate Grids : %v", err)
	}
	// go HubManager.HubCleanupLoop()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		multiplayer.WebsocketHandler(w, r, HubManager)
	})
	fmt.Printf("Listening on port %v ...", 3005)

	http.ListenAndServe(":3005", nil)
	time.Sleep(100 * time.Millisecond)

}
