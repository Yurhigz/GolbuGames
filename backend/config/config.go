package config

import (
	"context"
	"fmt"
	"golbugames/backend/internal/sudoku"
	"golbugames/backend/internal/sudoku/repository"
	"golbugames/backend/pkg/types"
	"golbugames/backend/pkg/utils"
	"runtime"
	"sync"
)

// Initialisation des grilles dans la DB , génération d'une centaine de grille de chaque difficulté

func InitGridGeneration(ctx context.Context) error {
	// Fonction de génération : GenerateSolvedGrid
	// Fonction de vérification de l'unicité : GeneratePlayableGrid
	// Fonction de transformation en string de la grid : GridTransformer()
	// Les grids sont stockées dans la DB dans leur forme plate string
	target := 100
	var bulkGrid []*types.SudokuGrid
	results := make(chan *types.SudokuGrid, 1000)
	difficulties := []string{"easy", "intermediate", "hard", "expert"}
	jobs := make(chan string, 1000) // buffer par défaut à 1000 just in case

	var wg sync.WaitGroup
	// var mu sync.Mutex
	go func() {
		for i := 0; i <= target; i++ {
			for _, difficulty := range difficulties {
				jobs <- difficulty
			}
		}
		close(jobs)
	}()

	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for difficulty := range jobs {

				solved, err := sudoku.GenerateSolvedGrid()
				if err != nil {
					fmt.Printf("[ERROR] <GenerateSolvedGrid> %v", err)
					return
				}

				playable, err := sudoku.GeneratePlayableGrid(solved, difficulty)
				if err != nil {
					fmt.Printf("[ERROR] <GeneratePlayableGridWithDifficulty> %v", err)
					return
				}

				grid := &types.SudokuGrid{
					Board:      utils.GridTransformer(playable),
					Solution:   utils.GridTransformer(solved),
					Difficulty: difficulty,
				}

				results <- grid

			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for value := range results {
		bulkGrid = append(bulkGrid, value)
	}

	err := repository.BulkAddingGrids(ctx, bulkGrid)

	if err != nil {
		fmt.Printf("[ERROR] %v", err)
		return err
	}

	return nil
}
