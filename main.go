package main

import "fmt"

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
	fmt.Println(subgridCol, subgridRow)
	for i := subgridRow; i < subgridRow+3; i++ {
		for j := subgridCol; j < subgridCol+3; j++ {
			if grid[i][j] == value {
				return false
			}
		}
	}
	return true
}

// func sudokuSolver(grid *mainGrid) bool {
// 	intermediateGrid := grid
// 	//  backtracking algorithm
// 	for row := 0; row < 9; row++ {
// 		for col := 0; col < 9; col++ {

// 		}
// 	}
// }

func main() {
	Sudoku := mainGrid{
		{0, 8, 0, 0, 7, 0, 0, 0, 0},
		{6, 3, 0, 1, 9, 5, 0, 0, 0},
		{0, 0, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 0, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 0},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	// fmt.Println(isValid(&Sudoku, 6, 6, 5))
	// fmt.Println(isValid(&Sudoku, 0, 0, 3))
	fmt.Println(isValid(&Sudoku, 5, 1, 9))
}
