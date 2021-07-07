package systems

import (
	"fmt"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/state"
)

func AI(engine *ecs.Engine, gameState *state.GameState) {
	aiEntities := engine.Entities.GetEntities([]string{state.AI})
	for _, e := range aiEntities {
		fmt.Println(state.GetDescription(e) + " ponders it's existence.")
	}
}
