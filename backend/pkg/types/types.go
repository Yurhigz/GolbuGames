package types

import (
	"time"
)

// MainGrid représente une grille de Sudoku 9x9
type MainGrid [9][9]int

// Coordinates représente une position dans la grille
type Coordinates [2]int

// // Game interface définit les méthodes communes à tous les jeux
// type Game interface {
// 	Initialize() error
// 	Validate() bool
// 	IsComplete() bool
// }

type Tournament struct {
	ID          int
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

// // SudokuGame implémente l'interface Game
type SudokuGrid struct {
	Board      []int  `json:"board"`
	Solution   []int  `json:"solution"`
	Difficulty string `json:"difficulty"`
}

// Ajout user API
type UserRegistration struct {
	Username    string `json:"username"`
	Accountname string `json:"accountname"`
	Password    string `json:"password"`
}

type User struct {
	Username    string `json:"username"`
	Accountname string `json:"accountname"`
	Password    string `json:"password"`
	ID          int    `json:"id"`
}

type AddFriendRequest struct {
	UserID         int    `json:"user_id"`
	FriendUsername string `json:"friend_username"`
}

type GridRequest struct {
	Difficulty string `json:"difficulty"`
}

// User stats API
type UserStats struct {
	ID           int `json:"id"`
	Total_games  int `json:"total_games"`
	Total_wins   int `json:"total_wins"`
	Total_losses int `json:"total_losses"`
	Total_draws  int `json:"total_draws"`
	Total_time   int `json:"total_time"`
	Average_time int `json:"average_time"`
}

type Game struct {
	UserID          int    `json:"user_id"`
	OpponentID      *int   `json:"opponent_id,omitempty"`
	GameMode        string `json:"game_mode"`
	Results         *int   `json:"results,omitempty"`
	Completion_time int    `json:"completion_time"`
	Difficulty      string `json:"difficulty"`
}

type PasswordUpdate struct {
	ID          int    `json:"id"`
	NewPassword string `json:"new_password"`
}
