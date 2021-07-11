package systems_test

import (
	"testing"

	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/systems"
)

func TestAi(t *testing.T) {
	gameState := game.NewGameState()
	systems.AI(gameState.Engine, gameState)
}
