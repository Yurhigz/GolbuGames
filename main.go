package main

import (
	"fmt"
	"math/rand"
)

type mainGrid [9][9]int

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
	return false
}

func main() {
	// Sudoku := mainGrid{
	// 	{0, 8, 0, 0, 7, 0, 0, 0, 0},
	// 	{6, 3, 0, 1, 9, 5, 0, 0, 0},
	// 	{0, 0, 8, 0, 0, 0, 0, 6, 0},
	// 	{8, 0, 0, 0, 6, 0, 0, 0, 3},
	// 	{4, 0, 0, 8, 0, 3, 0, 0, 1},
	// 	{7, 0, 0, 0, 2, 0, 0, 0, 6},
	// 	{0, 6, 0, 0, 0, 0, 0, 8, 0},
	// 	{0, 0, 0, 4, 1, 9, 0, 0, 0},
	// 	{0, 0, 0, 0, 8, 0, 0, 7, 9},
	// }
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

}
