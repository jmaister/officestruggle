package systems

import (
	"fmt"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/state"
)

func Render(engine *ecs.Engine) {
	renderable := []string{"position", "apparence"}
	for _, entity := range engine.GetEntities(renderable) {
		fmt.Println(entity)

		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)
		fmt.Println("pos", position)

		apparence, _ := entity.GetComponent(state.Apparence).(state.ApparenceComponent)
		fmt.Println("app", apparence)
	}
}
