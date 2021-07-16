package main

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/game"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Game state
	game := game.NewGame()

	ebiten.SetWindowSize(game.GameState.ScreenWidth, game.GameState.ScreenHeight)
	ebiten.SetWindowTitle("Office Struggle")
	ebiten.SetWindowResizable(true)

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

}
