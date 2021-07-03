package systems

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/state"
)

func CssToAttr(cssColor string) tl.Attr {
	c, err := ecs.ParseHexColorFast(cssColor)
	if err != nil {
		return tl.ColorDefault
	}
	return tl.RgbTo256Color(int(c.R), int(c.G), int(c.B))
}

func Render(engine *ecs.Engine, level *tl.BaseLevel) {
	renderable := []string{"position", "apparence"}
	for _, entity := range engine.GetEntities(renderable) {
		fmt.Println(entity)

		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)
		fmt.Println("pos", position)

		apparence, _ := entity.GetComponent(state.Apparence).(state.ApparenceComponent)
		fmt.Println("app", apparence)

		player := tl.NewEntity(position.X, position.Y, 1, 1)
		player.SetCell(0, 0, &tl.Cell{Fg: CssToAttr(apparence.Color), Ch: apparence.Char})
		level.AddEntity(player)
	}
}
