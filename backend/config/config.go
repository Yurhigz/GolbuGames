package config

import (
	"context"
	"golbugames/backend/internal/database"
	"golbugames/backend/internal/games"
	"golbugames/backend/pkg/utils"
)

// Ajouter une initialisation de la création des grilles complètes
// Envisager de mettre en place un système de contrôle pour générer des grilles selon le volume réalisé ou le nombre d'utilisateurs

func GenerateAndStoreGrid(ctx context.Context, difficulty string) error {
	grid, err := games.GenerateSolvedGrid()
	if err != nil {
		return err
	}

	var board, solution string

	solution = utils.GridTransformer(grid)

	boardGrid, err := games.GeneratePlayableGrid(grid, difficulty)
	if err != nil {
		return err
	}

	board = utils.GridTransformer(boardGrid)

	err = database.AddGrid(ctx, board, solution, difficulty)

	if err != nil {
		return err
	}
	return nil
}
