package systems_test

import (
	"fmt"
	"testing"

	"jordiburgos.com/officestruggle/game"
)

func TestAi(t *testing.T) {
	gameState := game.NewGameState()
	// systems.AI(gameState.Engine, gameState)
	fmt.Println(gameState)
}
