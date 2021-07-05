package systems

import (
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

func Render(engine *ecs.Engine, screen *tl.Screen) {
	renderable := []string{"position", "apparence"}
	for _, entity := range engine.GetEntities(renderable) {
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)
		apparence, _ := entity.GetComponent(state.Apparence).(state.ApparenceComponent)

		bg := apparence.Bg
		if bg == "" {
			bg = "#000"
		}

		screen.RenderCell(position.X, position.Y, &tl.Cell{Fg: CssToAttr(apparence.Color), Bg: CssToAttr(bg), Ch: apparence.Char})
	}

}
