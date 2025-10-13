package main

import (
	"context"
	"golbugames/backend/config"
	"golbugames/backend/internal/api/router"
	"golbugames/backend/internal/database"
	"golbugames/backend/internal/websocket/multiplayer"
	"log"
	"net/http"
	"time"
)

// Middleware CORS global
// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Skip CORS headers pour WebSocket
// 		if strings.HasPrefix(r.URL.Path, "/ws/") {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// CORS pour REST API
// 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true") // Important pour les cookies

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func main() {
	ctx := context.Background()

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

	// Init du multijoueurs HubManager
	hubManager := multiplayer.NewHubManager()
	go hubManager.MatchmakingLoop(ctx)

	// Init router
	r := router.NewRouter(hubManager)
	r.InitRoutes()
	// handler := corsMiddleware(r)

	server := &http.Server{
		Addr:         ":3002",
		Handler:      r, // Le middleware wrappé
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Print("Listening on :3002...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
