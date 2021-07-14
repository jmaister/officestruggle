package gamestate_test

import (
	"fmt"
	"testing"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func TestNewGame(t *testing.T) {
	engine := ecs.NewEngine()
	g := gamestate.NewGameState(engine)

	pos := g.Player.GetComponent(state.Position)
	fmt.Println(pos)
}
