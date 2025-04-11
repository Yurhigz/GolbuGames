package router

import (
	"golbugames/backend/internal/api/handlers"
	"net/http"
)

func InitRoutesSudoku(mux *http.ServeMux) {
	mux.HandleFunc("POST /create_user", handlers.CreateUser)
	mux.HandleFunc("DELETE /delete_user", handlers.DeleteUser)
	mux.HandleFunc("GET /user/{id}", handlers.GetUser)
	mux.HandleFunc("POST /add_grid", handlers.AddGrid)
	mux.HandleFunc("GET /random_grid", handlers.GetGrid)
	mux.HandleFunc("POST /submit_game", handlers.SubmitSoloGame)
	mux.HandleFunc("GET /user_stats", handlers.GetUserStats)
	mux.HandleFunc("GET /leaderboard", handlers.GetLeaderboard)
	mux.HandleFunc("POST /updateuser", handlers.UpdateUser)
	mux.HandleFunc("GET /user_history", handlers.GetUserHistory)
	mux.HandleFunc("POST /save_game", handlers.SaveGameProgress)

}
