package sudoku

import "math"

func EloCalculation(playerAElo, playerBElo int, result string) (int, int) {
	k := 30

	Ea := 1 / (1 + math.Pow(10, float64(playerBElo-playerAElo)/400))
	Eb := 1 / (1 + math.Pow(10, float64(playerAElo-playerBElo)/400))

	switch result {
	case "win":
		playerAElo += int(math.Round(float64(k) * (1 - Ea)))
		playerBElo += int(math.Round(float64(k) * (0 - Eb)))
	case "lose":
		playerAElo += int(math.Round(float64(k) * (0 - Ea)))
		playerBElo += int(math.Round(float64(k) * (1 - Eb)))
	case "draw":
		playerAElo += int(math.Round(float64(k) * (0.5 - Ea)))
		playerBElo += int(math.Round(float64(k) * (0.5 - Eb)))
	}
	return playerAElo, playerBElo
}
