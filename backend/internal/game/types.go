package game

// MainGrid représente une grille de Sudoku 9x9
type MainGrid [9][9]int

// Coordinates représente une position dans la grille
type Coordinates [2]int

// // Game interface définit les méthodes communes à tous les jeux
// type Game interface {
// 	Initialize() error
// 	Validate() bool
// 	IsComplete() bool
// }

// // SudokuGame implémente l'interface Game
// type SudokuGame struct {
// 	Grid     MainGrid
// 	Level    string
// 	Solution MainGrid
// }
