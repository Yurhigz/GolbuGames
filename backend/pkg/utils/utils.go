package utils

import (
	"golbugames/backend/pkg/types"
	"math/rand/v2"
	"strconv"
)

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

// Fonction de hashage pour les passwords

// Fonction de transformation des données de type sudokuGrid en mode string

func GridTransformer(sudokuGrid *types.MainGrid) string {
	var FlattenedSudoku string
	for _, row := range sudokuGrid {
		for i, col := range row {
			if i < len(row)-1 {
				FlattenedSudoku += strconv.Itoa(col) + ","
			} else {
				FlattenedSudoku += strconv.Itoa(col) + ";"
			}
		}
	}
	return FlattenedSudoku
}
