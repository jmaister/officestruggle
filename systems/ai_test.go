package systems_test

import (
	"testing"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/systems"
)

func TestAi(t *testing.T) {
	engine := ecs.NewEngine()
	gameState := game.NewGameState(engine)
	systems.AI(gameState.Engine, gameState)
}
