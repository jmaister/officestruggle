package game_test

import (
	"testing"

	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/gamestate"
)

func TestNewGame(t *testing.T) {
	g := game.NewGame()

	g.GameState.ScreenState = gamestate.GameScreen

	g.Update()
}
