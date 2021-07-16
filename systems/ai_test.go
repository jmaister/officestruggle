package systems_test

import (
	"testing"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/systems"
)

func TestAi(t *testing.T) {
	engine := ecs.NewEngine()
	gameState := gamestate.NewGameState(engine)
	systems.AI(gameState.Engine, gameState)
}
