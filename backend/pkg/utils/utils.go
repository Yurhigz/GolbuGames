package utils

import "math/rand/v2"

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

// Fonction de hashage pour les passwords

// Fonction de transformation des données de type sudokuGrid en mode string
