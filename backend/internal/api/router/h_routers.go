package router

import (
	"golbugames/backend/internal/api/handlers"
	"net/http"
)

func InitRoutesSudoku(mux *http.ServeMux) {
	mux.HandleFunc("/create_user", handlers.CreateUser)
	mux.HandleFunc("/delete_user", handlers.DeleteUser)
	mux.HandleFunc("/user", handlers.GetUser)
	mux.HandleFunc("/add_grid", handlers.AddGrid)
	mux.HandleFunc("/random_grid", handlers.GetGrid)
	mux.HandleFunc("/submit_game", handlers.SubmitGame)
	mux.HandleFunc("/user_stats", handlers.GetUserStats)
	mux.HandleFunc("/leaderboard", handlers.GetLeaderboard)
	mux.HandleFunc("/updateuser", handlers.UpdateUser)
	mux.HandleFunc("/user_history", handlers.GetUserHistory)
	mux.HandleFunc("/save_game", handlers.SaveGameProgress)

}
