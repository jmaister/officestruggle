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

func Render(engine *ecs.Engine, gameState *state.GameState, screen *tl.Screen) {
	layers := []string{state.Layer100, state.Layer300, state.Layer400}

	// Reset visibility
	visitables := engine.Entities.GetEntities([]string{state.Visitable})
	state.SetVisibleEntities(visitables, false)

	// Update visibility
	player := engine.Entities.GetEntity([]string{state.Player})
	position := state.GetPosition(player)
	gameState.Fov.RayCast(engine, position.X, position.Y, &gameState.Grid.Map)

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
		visitable, _ := entity.GetComponent(state.Visitable).(state.VisitableComponent)

		fg := apparence.Color
		if fg == "" || len(fg) == 0 {
			fg = "#FFF"
		}
		bg := apparence.Bg
		if bg == "" || len(bg) == 0 {
			bg = "#000"
		}
		ch := apparence.Char

		if visitable.Visible {
			screen.RenderCell(position.X, position.Y, &tl.Cell{Fg: CssToAttr(fg), Bg: CssToAttr(bg), Ch: ch})
		} else if visitable.Explored {
			screen.RenderCell(position.X, position.Y, &tl.Cell{Fg: CssToAttr("#CCC"), Bg: CssToAttr(bg), Ch: ch})
		} else {
			screen.RenderCell(position.X, position.Y, &tl.Cell{Fg: CssToAttr("#F00"), Bg: CssToAttr(bg), Ch: ch})
		}
	}
}
