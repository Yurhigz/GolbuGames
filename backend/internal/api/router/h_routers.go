package router

import (
	"golbugames/backend/internal/api/handlers"
	"net/http"
	"golbugames/backend/internal/api/middleware"
)

func InitRoutesSudoku(mux *http.ServeMux) {
	// Users API
	mux.HandleFunc("POST /create_user", handlers.CreateUser)
// 	mux.HandleFunc("POST /create_user", handlers.UserSignin)
	mux.HandleFunc("DELETE /delete_user/{id}", handlers.DeleteUser)
	mux.HandleFunc("GET /user/{id}", handlers.GetUser)

	mux.HandleFunc("POST /refresh_token", handlers.RefreshTokenHandler)
	mux.HandleFunc("POST /logout", handlers.LogoutHandler)

	mux.HandleFunc("POST /updateuser", middleware.JWTMiddleware("", handlers.UpdateUserPassword))
	
	mux.HandleFunc("GET /user_stats/{id}", handlers.GetUserStats)
	mux.HandleFunc("GET /user_id", handlers.GetUserId)
	mux.HandleFunc("POST /login", handlers.UserLogin)

	// Sudoku API
	mux.HandleFunc("POST /add_grid", handlers.AddGrid)
	mux.HandleFunc("GET /grid", handlers.GetGrid)

	// Game API*ws.HubManager
	mux.HandleFunc("POST /submit_solo_game", handlers.SubmitSoloGame)
	mux.HandleFunc("POST /submit_multi_game", middleware.JWTMiddleware("", handlers.SubmitMultiGame))

	mux.HandleFunc("GET /leaderboard", handlers.GetLeaderboard)
	// mux.HandleFunc("GET /user_history", handlers.GetUserHistory)
	// mux.HandleFunc("POST /save_game", handlers.SaveGameProgress)

	mux.HandleFunc("GET /friends", middleware.JWTMiddleware("", handlers.GetUserFriends))
	mux.HandleFunc("DELETE /delete_friend/{f_id}", middleware.JWTMiddleware("", handlers.RemoveFriend))
	mux.HandleFunc("POST /add_friend", middleware.JWTMiddleware("", handlers.AddFriend))

	mux.HandleFunc("GET /tournaments", middleware.JWTMiddleware("", handlers.GetAllTournaments))
	mux.HandleFunc("POST /add_tournament", middleware.JWTMiddleware("", handlers.AddTournament))
// 	mux.HandleFunc("GET /tournament/{id}", handlers.GetTournamentHandler)
//     mux.HandleFunc("POST /tournament/{id}/join", handlers.JoinTournamentHandler)
//     mux.HandleFunc("GET /tournament/{id}/participants", handlers.GetParticipantsHandler)
//     mux.HandleFunc("POST /tournament/{id}/start", handlers.StartTournamentHandler)
//     mux.HandleFunc("POST /tournament/{id}/finish", handlers.StartTournamentHandler)
//     mux.HandleFunc("GET /tournament/{id}/ranking", handlers.GetTournamentHandler)


}
