package main

import (
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/systems"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Game state
	gameState := game.NewGameState()

	tlGame := tl.NewGame()
	tlGame.Screen().SetFps(30)
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})

	ctl := systems.NewController(gameState)
	level.AddEntity(ctl)

	tlGame.Screen().SetLevel(level)
	tlGame.Start()
}
