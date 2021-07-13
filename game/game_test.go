package game_test

import (
	"fmt"
	"testing"

	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/state"
)

func TestNewGame(t *testing.T) {
	g := game.NewGameState()

	pos := g.Player.GetComponent(state.Position)
	fmt.Println(pos)
}
