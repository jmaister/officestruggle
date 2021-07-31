package gamestate_test

import (
	"fmt"
	"testing"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

func TestNewGame(t *testing.T) {
	engine := ecs.NewEngine()
	g := gamestate.NewGameState(engine)

	pos := g.Player.GetComponent(constants.Position)
	fmt.Println(pos)
}
