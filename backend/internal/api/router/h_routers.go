package router

import (
	"golbugames/internal/api/handlers"
	"golbugames/internal/api/middleware"
	"golbugames/internal/websocket/multiplayer"
	"golbugames/internal/websocket/solo"
	"net/http"
)

func InitRoutesSudoku(mux *http.ServeMux, hubManager *multiplayer.HubManager, jwtManager *middleware.JWTManager) {
	// Users API
	mux.HandleFunc("POST /create_user", handlers.CreateUserWithJWT(jwtManager))
	// 	mux.HandleFunc("POST /create_user", handlers.UserSignin)
	mux.Handle("DELETE /delete_user/{id}", middleware.RequireJWT(jwtManager, http.HandlerFunc(handlers.DeleteUser)))
	mux.HandleFunc("GET /user/{id}", handlers.GetUser)

	mux.Handle("POST /updateuser", middleware.RequireJWT(jwtManager, http.HandlerFunc(handlers.UpdateUserPassword)))

	mux.HandleFunc("GET /user_stats/{id}", handlers.GetUserStats)
	mux.HandleFunc("GET /user_id", handlers.GetUserId)
	mux.HandleFunc("POST /login", handlers.UserLoginWithJWT(jwtManager))

	// Sudoku API
	mux.Handle("POST /add_grid", middleware.RequireJWT(jwtManager, http.HandlerFunc(handlers.AddGrid)))
	// mux.HandleFunc("GET /grid", handlers.GetGrid)
	mux.HandleFunc("GET /grid", solo.WebsocketHandlerSolo)

	// Game API
	mux.Handle("POST /submit_solo_game", middleware.RequireJWT(jwtManager, http.HandlerFunc(handlers.SubmitSoloGame)))

	mux.HandleFunc("GET /leaderboard", handlers.GetLeaderboard)
	mux.HandleFunc("GET /metrics/multiplayer", handlers.GetMultiplayerMetrics(hubManager))
	// mux.HandleFunc("GET /user_history", handlers.GetUserHistory)
	// mux.HandleFunc("POST /save_game", handlers.SaveGameProgress)

	mux.HandleFunc("GET /friends/{id}", handlers.GetUserFriends)
	mux.Handle("DELETE /delete_friend/{id}/{f_id}", middleware.RequireJWT(jwtManager, http.HandlerFunc(handlers.RemoveFriend)))
	mux.Handle("POST /add_friend", middleware.RequireJWT(jwtManager, http.HandlerFunc(handlers.AddFriend)))

	mux.HandleFunc("GET /tournaments", handlers.GetAllTournaments)
	mux.Handle("POST /add_tournament", middleware.RequireJWT(jwtManager, http.HandlerFunc(handlers.AddTournament)))

	// ========== WEBSOCKET ENDPOINTS ==========

	// WebSocket Solo
	mux.HandleFunc("GET /ws/solo", solo.WebsocketHandlerSolo)

	// WebSocket Multiplayer
	mux.HandleFunc("GET /ws/multi", func(w http.ResponseWriter, r *http.Request) {
		multiplayer.WebsocketHandler(w, r, hubManager, jwtManager)
	})

}
