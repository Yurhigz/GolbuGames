package game

import (
	"errors"
	"golbugames/backend/pkg/utils"
	"math/rand"
)

func isValid(grid *MainGrid, row, col, value int) bool {
	// Je vérifie dans les lignes et colonnes
	for i := 0; i < 9; i++ {

		if grid[row][i] == value || grid[i][col] == value {

			return false
		}
	}
	// Il faut que la case soit vide
	if grid[row][col] != 0 {
		return false
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

func sudokuSolver(grid *MainGrid) bool {
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
			if sudokuSolver(grid) {
				return true
			}
			grid[row][col] = 0
		}

	}
	return false
}

func validNbIndices(numbers int) (bool, error) {
	if numbers > 54 || numbers < 0 {
		return false, errors.New("")
	}
	return true, nil
}

func symmetryRandomRemoving(numbers int, grid *MainGrid) (bool, error) {
	if numbers%2 != 0 {
		return false, errors.New("numbers must be even for symmetrical requirements")
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

func basicRandomRemoving(numbers int, grid *MainGrid) (bool, error) {
	valid, err := validNbIndices(numbers)
	if !valid {
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

// 43-45 cases dissimulées => easy
// 46-51 cases dissimulées => medium
// 52-55 cases dissimulées => hard
// 56-60 cases dissimulées => expert

func GenerateGrid(difficulty string) (*MainGrid, error) {
	var emptyGrid MainGrid
	v, err := difficultySelection(difficulty)
	if err != nil {
		return &emptyGrid, err
	}
	sudokuSolver(&emptyGrid)
	switch v {
	case 0:
		indices := utils.RandRange(43, 46)
		symmetryRandomRemoving(indices, &emptyGrid)
	case 1:
		indices := utils.RandRange(46, 52)
		symmetryRandomRemoving(indices, &emptyGrid)
	case 2:
		indices := utils.RandRange(52, 56)
		basicRandomRemoving(indices, &emptyGrid)
	case 3:
		indices := utils.RandRange(56, 61)
		basicRandomRemoving(indices, &emptyGrid)
	}

	return &emptyGrid, nil
}
