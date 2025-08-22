package router

import (
	"golbugames/backend/internal/api/handlers"
	"net/http"
)

func InitRoutesSudoku(mux *http.ServeMux) {
	// Users API
	mux.HandleFunc("POST /create_user", handlers.CreateUser)
	mux.HandleFunc("DELETE /delete_user/{id}", handlers.DeleteUser)
	mux.HandleFunc("GET /user/{id}", handlers.GetUser)

	mux.HandleFunc("POST /updateuser", handlers.UpdateUserPassword)
	
	mux.HandleFunc("GET /user_stats/{id}", handlers.GetUserStats)
	mux.HandleFunc("GET /user_id", handlers.GetUserId)
	mux.HandleFunc("POST /login", handlers.UserLogin)

	// Sudoku API
	mux.HandleFunc("POST /add_grid", handlers.AddGrid)
	mux.HandleFunc("GET /grid", handlers.GetGrid)

	// Game API
	mux.HandleFunc("POST /submit_solo_game", handlers.SubmitSoloGame)
	mux.HandleFunc("POST /submit_multi_game", handlers.SubmitMultiGame)

	mux.HandleFunc("GET /leaderboard", handlers.GetLeaderboard)
	// mux.HandleFunc("GET /user_history", handlers.GetUserHistory)
	// mux.HandleFunc("POST /save_game", handlers.SaveGameProgress)

	mux.HandleFunc("GET /friends/{id}", handlers.GetUserFriends)
	mux.HandleFunc("DELETE /delete_friend/{id}/{f_id}", handlers.RemoveFriend)
	mux.HandleFunc("POST /add_friend", handlers.AddFriend)

	mux.HandleFunc("GET /tournaments", handlers.GetAllTournaments)
	mux.HandleFunc("POST /add_tournament", handlers.AddTournament)

}
