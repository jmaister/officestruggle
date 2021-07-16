package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

var (
	mplusFontCached map[float64]font.Face
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
			visitable.AddComponent(state.Visitable, vsComponent)
		}
	}

	for _, layer := range layers {
		renderable := []string{state.Position, state.Apparence, layer}
		entities := engine.Entities.GetEntities(renderable)

		renderEntities(entities, gameState, screen)
	}

	drawMessageLog(screen, gameState)
}

func showDebug(screen *ebiten.Image) {
	w, _ := screen.Size()
	// Draw info
	fnt := mplusFont(10)
	msg := fmt.Sprintf("TPS: %0.2f, FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msg, fnt, w-100, 20, color.White)
}

func renderEntities(entities []*ecs.Entity, gameState *gamestate.GameState, screen *ebiten.Image) {

	w := gameState.ScreenWidth
	h := gameState.ScreenHeight
	tw := w / gameState.Grid.Width
	th := h / gameState.Grid.Height

	font := mplusFont(float64(th))

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

				// Mark as explored
				vsComponent, _ := entity.RemoveComponent(state.Visitable).(state.VisitableComponent)
				vsComponent.Explored = true
				entity.AddComponent(state.Visitable, vsComponent)

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
			}
		}
	}
}

func drawChar(screen *ebiten.Image, str string, x int, y int, font font.Face, fgColor color.Color, bgColor color.Color) {
	drawBackground(screen, str, x, y, font, bgColor)
	text.Draw(screen, str, font, x, y, fgColor)
}

func drawBackground(screen *ebiten.Image, str string, x int, y int, face font.Face, bgColor color.Color) {
	rect := text.BoundString(face, str)
	ebitenutil.DrawRect(screen, float64(x+rect.Min.X), float64(y+rect.Min.Y), float64(rect.Max.X-rect.Min.X), float64(rect.Max.Y-rect.Min.Y), bgColor)
}

func mplusFont(size float64) font.Face {
	if mplusFontCached == nil {
		mplusFontCached = map[float64]font.Face{}
	}
	fnt, ok := mplusFontCached[size]
	if !ok {
		tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
		if err != nil {
			panic(err)
		}

		fnt, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			panic(err)
		}
		mplusFontCached[size] = fnt
	}
	return fnt
}

var messageLogColors = [5]color.RGBA{
	ParseHexColorFast("#A9A9A9"),
	ParseHexColorFast("#C0C0C0"),
	ParseHexColorFast("#D3D3D3"),
	ParseHexColorFast("#DCDCDC"),
	ParseHexColorFast("#FFFFFF"),
}

func drawMessageLog(screen *ebiten.Image, gs *gamestate.GameState) {

	fontSize := 15
	font := mplusFont(float64(fontSize))

	position := gs.Grid.MessageLog

	lines := gs.GetLog(5)
	n := len(lines)
	for i, line := range lines {
		fgColor := messageLogColors[5-n+i]
		text.Draw(screen, line, font, (position.X)*fontSize, (position.Y+i+1)*fontSize, fgColor)
	}

}
