package interfaces

import (
	"context"
	"golbugames/backend/pkg/types"
)

type SudokuRepository interface {
	// Grid operations
	AddGridDB(ctx context.Context, board, solution, difficulty string) error
	GetGridDB(ctx context.Context, id int) (string, string, error)
	GetRandomGridDB(ctx context.Context, difficulty string) (*types.SudokuGrid, error)
}

type UserRepository interface {
	// User operations
	AddUserDB(ctx context.Context, username, accountname, password string) error
	DeleteUserDB(ctx context.Context, userID int) error
	GetUserDB(ctx context.Context, userID int) (*types.User, error)
	UpdateUserPasswordDB(ctx context.Context, userID int, newPassword string) error
	GetUserStatsDB(ctx context.Context, userID int) (*types.UserStats, error)
}

type GameRepository interface {
	// Game operations
	SubmitSoloGameDB(ctx context.Context, userID, completionTime int) error
	SubmitMultiGameDB(ctx context.Context, user1, user2, results, completionTime int) error
}
