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

func HashPassword(password string) (string, error) {
	// "Latacora, 2018: In order of preference, use scrypt, argon2, bcrypt, and then if nothing else is available PBKDF2."
	// Check les librairies scrypt et argon2
	// hachage =/= encryption
}

// comparaison du mot de passe utilisateur et du hash en BDD
func ValidatePassword(password, hash string) bool {
	// Récupération du hash dans la bdd en fonction de l'utilisateur et comparaison avec la proposition du user
}

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
