package sudoku

import (
	"errors"
	"fmt"
	"golbugames/backend/pkg/types"
	"golbugames/backend/pkg/utils"
	"math/rand"
)

func isValid(grid *types.MainGrid, row, col, value int) bool {
	// Il faut que la case soit vide
	if grid[row][col] != 0 {
		return false
	}

	for i := 0; i < 9; i++ {

		if grid[row][i] == value || grid[i][col] == value {

			return false
		}
	}

	// je vérifie dans la subgrid
	subgridCol := (col / 3) * 3
	subgridRow := (row / 3) * 3
	for i := subgridRow; i < subgridRow+3; i++ {
		for j := subgridCol; j < subgridCol+3; j++ {
			if grid[i][j] == value {
				return false
			}
		}
	}
	return true
}

func SudokuSolver(grid *types.MainGrid) bool {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Shuffle(len(numbers), func(i, j int) { numbers[i], numbers[j] = numbers[j], numbers[i] })

	// Chercher les cases vides
	var row, col int
	emptyFound := false

	for i := 0; i < 9 && !emptyFound; i++ {
		for j := 0; j < 9; j++ {
			if grid[i][j] == 0 {
				row = i
				col = j
				emptyFound = true
			}
		}
	}

	if !emptyFound {
		return true
	}

	//  backtracking algorithm
	for _, number := range numbers {
		if isValid(grid, row, col, number) {
			grid[row][col] = number
			if SudokuSolver(grid) {
				return true
			}
			grid[row][col] = 0
		}

	}
	return false
}

func validNbIndices(numbers int) (bool, error) {
	if numbers > 64 || numbers < 0 {
		return false, errors.New("number of indices to be removed must be between 54 and 0 excluded")
	}
	return true, nil
}

func CountSolutions(grid *types.MainGrid, limit int) int {
	var count int

	var solve func(*types.MainGrid) bool
	solve = func(g *types.MainGrid) bool {
		// Chercher la première case vide
		var row, col int
		found := false
		for i := 0; i < 9 && !found; i++ {
			for j := 0; j < 9; j++ {
				if g[i][j] == 0 {
					row, col = i, j
					found = true
					break
				}
			}
		}

		// plus de case vide = une solution trouvée
		if !found {
			count++
			return count >= limit // stop si on a atteint la limite
		}

		for v := 1; v <= 9; v++ {
			if isValid(g, row, col, v) {
				g[row][col] = v
				if solve(g) { // si assez de solutions trouvées
					g[row][col] = 0
					return true
				}
				g[row][col] = 0
			}
		}
		return false
	}

	copyGrid := *grid
	solve(&copyGrid)
	return count
}

func symmetryRandomRemoving(numbers int, grid *types.MainGrid) (bool, error) {
	if numbers%2 != 0 {
		numbers++
	}
	valid, err := validNbIndices(numbers)
	if !valid {
		return valid, err
	}

	removed := 0
	attempts := 0
	for removed < numbers/2 && attempts < 1000 {
		attempts++
		row := rand.Intn(9)
		col := rand.Intn(9)

		if grid[row][col] != 0 && grid[col][row] != 0 {
			backup1 := grid[row][col]
			backup2 := grid[col][row]

			grid[row][col] = 0
			grid[col][row] = 0

			if CountSolutions(grid, 2) == 1 {
				removed++
			} else {
				grid[row][col] = backup1
				grid[col][row] = backup2
			}
		}
	}
	if removed < numbers {
		return false, fmt.Errorf("failed to remove %d numbers, only removed %d", numbers, removed)
	}
	return true, nil
}

func basicRandomRemoving(numbers int, grid *types.MainGrid) (bool, error) {
	valid, err := validNbIndices(numbers)
	if !valid || err != nil {
		return valid, err
	}

	removed := 0
	attempts := 0
	for removed < numbers && attempts < 1000 {
		attempts++
		row := rand.Intn(9)
		col := rand.Intn(9)

		if grid[row][col] != 0 {
			backup := grid[row][col]
			grid[row][col] = 0

			// Vérifie que la grille garde une seule solution
			if CountSolutions(grid, 2) == 1 {
				removed++
			} else {
				grid[row][col] = backup
			}
		}
	}
	if removed < numbers {
		return false, fmt.Errorf("failed to remove %d numbers, only removed %d", numbers, removed)
	}
	return true, nil
}

func difficultySelection(difficulty string) (int, error) {
	switch difficulty {
	case "easy":
		return 0, nil
	case "intermediate":
		return 1, nil
	case "advanced":
		return 2, nil
	case "expert":
		return 3, nil
	default:
		return -1, errors.New("difficulty must be a string such as easy")
	}
}
func GenerateSolvedGrid() (*types.MainGrid, error) {
	var emptyGrid types.MainGrid
	if !SudokuSolver(&emptyGrid) {
		return nil, errors.New("failed to generate valid grid")
	}
	solvedGrid := emptyGrid
	return &solvedGrid, nil
}

func GeneratePlayableGrid(solvedGrid *types.MainGrid, difficulty string) (*types.MainGrid, error) {

	var emptyGrid types.MainGrid
	v, err := difficultySelection(difficulty)
	if err != nil {
		return &emptyGrid, err
	}
	playableGrid := *solvedGrid

	var success bool

	switch v {
	case 0:
		// 43-45 cases dissimulées => easy
		indices := utils.RandRange(43, 46)
		success, err = symmetryRandomRemoving(indices, &playableGrid)
	case 1:
		// 46-51 cases dissimulées => medium
		indices := utils.RandRange(46, 52)
		success, err = symmetryRandomRemoving(indices, &playableGrid)
	case 2:
		// 52-55 cases dissimulées => hard
		indices := utils.RandRange(52, 56)
		success, err = basicRandomRemoving(indices, &playableGrid)
	case 3:
		// 56-60 cases dissimulées => expert
		indices := utils.RandRange(56, 61)
		success, err = basicRandomRemoving(indices, &playableGrid)
	}

	if err != nil || !success {
		return nil, fmt.Errorf("failed to remove numbers: %v", err)
	}

	return &playableGrid, nil
}
