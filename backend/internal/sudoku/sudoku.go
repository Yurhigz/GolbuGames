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

func symmetryRandomRemoving(numbers int, grid *types.MainGrid) (bool, error) {
	if numbers%2 != 0 {
		numbers += 1
	}
	valid, err := validNbIndices(numbers)
	if !valid {
		return valid, err
	}

	removed := 0

	for removed < numbers/2 {
		row := rand.Intn(9)
		col := rand.Intn(9)

		if grid[row][col] != 0 && grid[col][row] != 0 {
			grid[row][col] = 0
			grid[col][row] = 0
			removed++
		}

	}
	return true, nil
}

func basicRandomRemoving(numbers int, grid *types.MainGrid) (bool, error) {
	valid, err := validNbIndices(numbers)
	if !valid || err != nil {
		return valid, err
	}
	removed := 0

	for removed < numbers {
		row := rand.Intn(9)
		col := rand.Intn(9)

		if grid[row][col] != 0 {
			grid[row][col] = 0
			removed++
		}
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
