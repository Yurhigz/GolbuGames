package sudoku

import "math"

func EloCalculation(playerAElo, playerBElo int, result string) (int, int) {
	// Formule Arpad Elo : Ea = 1 / (1 + 10 ^ ((Rb - Ra) / 400)) probabilité de victoire du joueur A
	// Ra = classement du joueur A (gagnant)	Rb = classement du joueur B (perdant)
	// Mise à jour de l'elo : R'a = Ra + K * (S - Ea)		S = 1 si le joueur A gagne, 0 si il perd, 0.5 si match nul
	// On part du principe que K est fixe peu importe le type de joueurs/match/Elo
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
