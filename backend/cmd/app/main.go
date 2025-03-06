package main

import (
	"errors"
	"fmt"
	"math/rand"
)

type mainGrid [9][9]int

type Coordinates [2]int

func isValid(grid *mainGrid, row, col, value int) bool {
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

func sudokuSolver(grid *mainGrid) bool {
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

func symmetryRandomRemoving(numbers int, grid *mainGrid) (bool, error) {
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

func basicRandomRemoving(numbers int, grid *mainGrid) (bool, error) {
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

func difficulty(difficulty string) (int, error) {
	for _, c := range difficulty {
		if (c < 48 || c > 57) && (c < 65 || c > 90) && (c < 97 || c > 122) {
			return -1, errors.New("difficulty must be a string such as easy")
		}
	}
	switch difficulty {
	case "easy":
		return 0, nil
	case "medium":
		return 1, nil
	case "hard":
		return 2, nil
	default:
		return -1, nil
	}
}

func generateGrid(difficulty string) bool {
	emptyGrid := mainGrid{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	sudokuSolver(&emptyGrid)

	return false
}

func main() {
	Sudoku2 := mainGrid{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	// fmt.Println(isValid(&Sudoku, 6, 6, 5))
	// fmt.Println(isValid(&Sudoku, 0, 0, 3))
	// fmt.Println(isValid(&Sudoku, 5, 1, 9))
	sudokuSolver(&Sudoku2)
	fmt.Println(Sudoku2)
	for row := 0; row < 9; row += 2 {
		Sudoku2[row][0] = 0
	}

	fmt.Println(Sudoku2)

	sudokuSolver(&Sudoku2)

	fmt.Println(Sudoku2)

}
