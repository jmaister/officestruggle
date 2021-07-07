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
	layers := []string{state.Layer100, state.Layer300, state.Layer400}

	for _, layer := range layers {
		renderable := []string{state.Position, state.Apparence, layer}
		entities := engine.Entities.GetEntities(renderable)
		renderEntities(entities, screen)
	}
}

func renderEntities(entities []*ecs.Entity, screen *tl.Screen) {
	for _, entity := range entities {
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)
		apparence, _ := entity.GetComponent(state.Apparence).(state.ApparenceComponent)

		bg := apparence.Bg
		if bg == "" || len(bg) == 0 {
			bg = "#000"
		}

		screen.RenderCell(position.X, position.Y, &tl.Cell{Fg: CssToAttr(apparence.Color), Bg: CssToAttr(bg), Ch: apparence.Char})
	}
}
