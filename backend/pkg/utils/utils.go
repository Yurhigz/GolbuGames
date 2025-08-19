package utils

import (
	"golbugames/backend/pkg/types"
	"log"
	"math/rand/v2"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

// Fonction de hashage pour les passwords

func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Println("password cannot be hashed : ", err)
		return "", err
	}

	return string(bytes), nil
}

// comparaison du mot de passe utilisateur et du hash en BDD
func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}
	return false
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
