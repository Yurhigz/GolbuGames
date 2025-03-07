package main

import (
	"fmt"
	"golbugames/backend/internal/game"
)

func main() {

	grid, err := game.GenerateGrid("easy")

	if err != nil {
		return
	}

	fmt.Println(grid)

}
