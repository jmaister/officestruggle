package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

func Render(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	showDebug(screen)

	layers := []string{state.Layer100, state.Layer300, state.Layer400}

	// Reset visibility
	visitables := engine.Entities.GetEntities([]string{state.Visitable})
	state.SetVisibleEntities(visitables, false)

	// Update visibility
	player := gameState.Player
	position := state.GetPosition(player)
	gameState.Fov.RayCast(engine, position.X, position.Y, &gameState.Grid.Map)

	for _, layer := range layers {
		renderable := []string{state.Position, state.Apparence, layer}
		entities := engine.Entities.GetEntities(renderable)

		renderEntities(entities, gameState, screen)
	}
}

func showDebug(screen *ebiten.Image) {
	// Draw info
	fnt := mplusFont(10)
	msg := fmt.Sprintf("TPS: %0.2f, FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msg, fnt, 20, 20, color.White)
}

func renderEntities(entities []*ecs.Entity, gameState *gamestate.GameState, screen *ebiten.Image) {

	font := mplusFont(16)

	for _, entity := range entities {
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)
		apparence, _ := entity.GetComponent(state.Apparence).(state.ApparenceComponent)
		visitable, isVisitable := entity.GetComponent(state.Visitable).(state.VisitableComponent)

		fg := apparence.Color
		if fg == "" || len(fg) == 0 {
			fg = "#FFF"
		}
		bg := apparence.Bg
		if bg == "" || len(bg) == 0 {
			bg = "#000"
		}
		ch := string(apparence.Char)

		// Pixel positions
		px := position.X * gameState.TileWidth
		py := position.Y * gameState.TileHeight

		if isVisitable {
			// Walls and floor
			if visitable.Visible {
				fgColor, _ := ecs.ParseHexColorFast(fg)
				text.Draw(screen, ch, font, px, py, fgColor)
			} else if visitable.Explored {
				fgColor := color.RGBA{
					R: 128,
					G: 128,
					B: 128,
					A: 128,
				}
				text.Draw(screen, ch, font, px, py, fgColor)
			}
		} else {
			fgColor, _ := ecs.ParseHexColorFast(fg)
			text.Draw(screen, ch, font, px, py, fgColor)
		}
	}
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
