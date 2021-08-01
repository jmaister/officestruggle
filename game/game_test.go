package game_test

import (
	"fmt"
	"testing"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
)

func TestNewGame(t *testing.T) {
	engine := ecs.NewEngine()
	g := game.NewGameState(engine)

	pos := g.Player.GetComponent(constants.Position)
	fmt.Println(pos)
}
