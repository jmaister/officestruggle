package systems

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"

	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func setVisibleEntities(entities ecs.EntityList, isVisible bool) {
	for _, e := range entities {
		visitable, _ := e.RemoveComponent(state.Visitable).(state.VisitableComponent)
		visitable.Visible = isVisible
		e.AddComponent(state.Visitable, visitable)
	}
}

func Render(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	showDebug(screen)

	layers := []string{state.Layer100, state.Layer300, state.Layer400, state.Layer500}

	// Reset visibility
	visitables := engine.Entities.GetEntities([]string{state.Visitable})
	setVisibleEntities(visitables, false)

	// Update visibility
	for _, visitable := range visitables {
		pos := state.GetPosition(visitable)
		if gameState.Fov.IsVisible(pos.X, pos.Y) {
			vsComponent, _ := visitable.RemoveComponent(state.Visitable).(state.VisitableComponent)
			vsComponent.Visible = true
			vsComponent.Explored = true
			visitable.AddComponent(state.Visitable, vsComponent)
		}
	}

	visibleEntities := []*ecs.Entity{}
	for _, layer := range layers {
		renderable := []string{state.Position, state.Apparence, layer}
		entities := engine.Entities.GetEntities(renderable)

		v := renderEntities(entities, gameState, screen)

		visibleEntities = append(visibleEntities, v...)
	}

	drawMessageLog(screen, gameState)
	drawPlayerHud(screen, gameState)
	drawInfo(screen, gameState, visibleEntities)
}

func showDebug(screen *ebiten.Image) {
	w, _ := screen.Size()
	// Draw info
	fnt := assets.MplusFont(10)
	msg := fmt.Sprintf("TPS: %0.2f, FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msg, fnt, w-150, 20, color.White)

	// Mouse info
	x, y := ebiten.CursorPosition()
	cursorStr := strconv.Itoa(x) + " " + strconv.Itoa(y)
	text.Draw(screen, cursorStr, fnt, w-150, 35, color.White)

}

func renderEntities(entities []*ecs.Entity, gameState *gamestate.GameState, screen *ebiten.Image) []*ecs.Entity {

	w := gameState.ScreenWidth
	h := gameState.ScreenHeight
	tw := w / gameState.Grid.Width
	th := h / gameState.Grid.Height

	font := assets.LoadFontCached(float64(20))

	visibleEntities := []*ecs.Entity{}

	for _, entity := range entities {
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)
		apparence, _ := entity.GetComponent(state.Apparence).(state.ApparenceComponent)
		visitable, isVisitable := entity.GetComponent(state.Visitable).(state.VisitableComponent)

		fg := apparence.Color
		if fg == "" || len(fg) == 0 {
			fg = "#FFFFFF"
		}
		bg := apparence.Bg
		if bg == "" || len(bg) == 0 {
			bg = "#000000"
		}
		ch := string(apparence.Char)

		// Pixel positions
		px := position.X * tw
		py := position.Y * th

		if isVisitable {
			// Walls and floor
			if visitable.Visible {
				bgColor := ParseHexColorFast(bg)
				fgColor := ParseHexColorFast(fg)
				drawChar(screen, ch, px, py, font, fgColor, bgColor)

			} else if visitable.Explored {
				bgColor := ParseHexColorFast("#000000")
				fgColor := ParseHexColorFast("#555555")
				drawChar(screen, ch, px, py, font, fgColor, bgColor)
			}
		} else {
			if gameState.Fov.IsVisible(position.X, position.Y) {
				bgColor := ParseHexColorFast(bg)
				fgColor := ParseHexColorFast(fg)
				drawChar(screen, ch, px, py, font, fgColor, bgColor)

				visibleEntities = append(visibleEntities, entity)
			}
		}
	}
	return visibleEntities
}

func drawChar(screen *ebiten.Image, str string, x int, y int, font font.Face, fgColor color.Color, bgColor color.Color) {
	drawBackground(screen, str, x, y, font, bgColor)
	text.Draw(screen, str, font, x, y, fgColor)
}

func drawBackground(screen *ebiten.Image, str string, x int, y int, face font.Face, bgColor color.Color) {
	rect := text.BoundString(face, str)
	pad := 0
	ebitenutil.DrawRect(screen, float64(x+rect.Min.X-pad), float64(y+rect.Min.Y-pad), float64(rect.Max.X-rect.Min.X+pad), float64(rect.Max.Y-rect.Min.Y+pad), bgColor)
}

var messageLogColors = [5]color.RGBA{
	ParseHexColorFast("#333333"),
	ParseHexColorFast("#555555"),
	ParseHexColorFast("#777777"),
	ParseHexColorFast("#AAAAAA"),
	ParseHexColorFast("#FFFFFF"),
}

func drawMessageLog(screen *ebiten.Image, gs *gamestate.GameState) {

	fontSize := 15
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.MessageLog

	lines := gs.GetLog(position.Height)
	n := len(lines)
	for i, line := range lines {
		fgColor := messageLogColors[5-n+i]
		logStr := line.Msg
		if line.Count > 1 {
			logStr = strconv.Itoa(line.Count) + "x " + line.Msg
		}
		text.Draw(screen, logStr, font, (position.X)*fontSize, (position.Y+i+1)*fontSize, fgColor)
	}

}

func drawPlayerHud(screen *ebiten.Image, gs *gamestate.GameState) {
	fontSize := 15
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.PlayerHud

	player := gs.Player
	stats := player.GetComponent(state.Stats).(state.StatsComponent)
	text.Draw(screen, stats.String(), font, (position.X)*fontSize, (position.Y+1)*fontSize, ParseHexColorFast("#00AA00"))
}

func drawInfo(screen *ebiten.Image, gs *gamestate.GameState, visibleEntities []*ecs.Entity) {
	fontSize := 15
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.InfoBar

	y := position.Y
	for _, entity := range visibleEntities {
		if !entity.HasComponent(state.Player) {
			str := state.GetLongDescription(entity)
			text.Draw(screen, str, font, (position.X)*fontSize, (y+1)*fontSize, ParseHexColorFast("#FFFFFF"))
			y++
		}
	}

}
