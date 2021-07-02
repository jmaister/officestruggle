package systems

import (
	"fmt"

	"jordiburgos.com/officestruggle/ecs"
)

func Render(engine *ecs.Engine) {
	renderable := []string{"position", "apparence"}
	for _, entity := range engine.GetEntities(renderable) {
		fmt.Println(entity)
	}
}
