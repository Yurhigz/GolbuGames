package main

import (
	"context"
	"golbugames/backend/internal/api/router"
	"golbugames/backend/internal/database"
	"log"
	"net/http"
)

// Middleware CORS global
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init DB
	err := database.InitDB(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Init router
	r := router.NewRouter()
	r.InitRoutes()

	// ✅ Wrap du router avec le middleware CORS
	handler := corsMiddleware(r)

	log.Print("Listening on :3001...")
	if err := http.ListenAndServe(":3001", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
