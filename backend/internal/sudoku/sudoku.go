package sudoku

import (
	"errors"
	"fmt"
	"golbugames/pkg/types"
	"golbugames/pkg/utils"
	"math/rand"
)

func isValid(grid *types.MainGrid, row, col, value int) bool {
	if grid[row][col] != 0 {
		return false
	}

	for i := 0; i < 9; i++ {

		if grid[row][i] == value || grid[i][col] == value {

			return false
		}
	}

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

		if !found {
			count++
			return count >= limit
		}

		for v := 1; v <= 9; v++ {
			if isValid(g, row, col, v) {
				g[row][col] = v
				if solve(g) {
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

func GenerateSolvedGrid() (*types.MainGrid, error) {
	var emptyGrid types.MainGrid
	if !SudokuSolver(&emptyGrid) {
		return nil, errors.New("failed to generate valid grid")
	}
	solvedGrid := emptyGrid
	return &solvedGrid, nil
}

func CloneGrid(grid *types.MainGrid) *types.MainGrid {
	var clone types.MainGrid
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			clone[i][j] = grid[i][j]
		}
	}
	return &clone
}

func SwapLines(grid *types.MainGrid, block int) {
	line1 := rand.Intn(3) + block*3
	line2 := rand.Intn(3) + block*3
	if line1 == line2 {
		line2 = (line1+1)%3 + block*3
	}
	grid[line1], grid[line2] = grid[line2], grid[line1]
}

func SwapColumns(grid *types.MainGrid, block int) {
	col1 := rand.Intn(3) + block*3
	col2 := rand.Intn(3) + block*3
	if col1 == col2 {
		col2 = (col1+1)%3 + block*3
	}
	for i := 0; i < 9; i++ {
		grid[i][col1], grid[i][col2] = grid[i][col2], grid[i][col1]
	}
}

func MirrorHorizontal(grid *types.MainGrid) {
	for i := 0; i < 4; i++ {
		grid[i], grid[8-i] = grid[8-i], grid[i]
	}
}

func MirrorVertical(grid *types.MainGrid) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 4; j++ {
			grid[i][j], grid[i][8-j] = grid[i][8-j], grid[i][j]
		}
	}
}

func Rotate90(grid *types.MainGrid) {
	var tmp types.MainGrid
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			tmp[j][8-i] = grid[i][j]
		}
	}
	*grid = tmp
}

func GenerateTransformedGrids(base *types.MainGrid, n int) []*types.MainGrid {
	grids := make([]*types.MainGrid, 0, n)
	unique := make(map[string]bool)

	for len(grids) < n {
		g := CloneGrid(base)

		for b := 0; b < 3; b++ {
			SwapLines(g, b)
			SwapColumns(g, b)
		}

		if rand.Intn(2) == 0 {
			MirrorHorizontal(g)
		}
		if rand.Intn(2) == 0 {
			MirrorVertical(g)
		}
		if rand.Intn(2) == 0 {
			Rotate90(g)
		}

		key := utils.GridTransformer(g)
		if !unique[key] {
			unique[key] = true
			grids = append(grids, g)
		}
	}
	return grids
}

func RemoveNumbersControlled(grid *types.MainGrid, count int) (*types.MainGrid, error) {

	g := CloneGrid(grid)

	if count <= 0 || count > 64 {
		return nil, fmt.Errorf("count doit être entre 1 et 64")
	}

	removed := 0
	attempts := 0
	maxAttempts := count * 10

	for removed < count && attempts < maxAttempts {
		attempts++
		row := rand.Intn(9)
		col := rand.Intn(9)

		if g[row][col] != 0 {
			backup := g[row][col]
			g[row][col] = 0

			if CountSolutions(g, 2) == 1 {
				removed++
			} else {
				g[row][col] = backup
			}
		}
	}

	if removed < count {
		return nil, fmt.Errorf("impossible de supprimer %d cases sans perdre l'unicité", count)
	}

	return g, nil
}
func GeneratePlayableGrid(baseGrid *types.MainGrid, difficulty string) (*types.MainGrid, error) {
	var count int
	switch difficulty {
	case "easy":
		count = utils.RandRange(43, 46)
	case "intermediate":
		count = utils.RandRange(46, 52)
	case "hard":
		count = utils.RandRange(52, 56)
	case "expert":
		count = utils.RandRange(56, 61)
	default:
		return nil, fmt.Errorf("difficulty inconnue")
	}

	return RemoveNumbersControlled(baseGrid, count)
}
