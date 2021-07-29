package systems

import (
	"fmt"
	"image/color"
	"math"
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

var fnt20 = assets.LoadFontCached(float64(20))
var fnt40 = assets.LoadFontCached(float64(40))

func setVisibleEntities(entities ecs.EntityList, isVisible bool) {
	for _, e := range entities {
		visitable, _ := e.RemoveComponent(state.Visitable).(state.VisitableComponent)
		visitable.Visible = isVisible
		e.AddComponent(state.Visitable, visitable)
	}
}

func Render(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	showDebug(screen, gameState)

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
	DrawGridRect(screen, gameState, gameState.Grid.Map, color.White)

	drawMessageLog(screen, gameState)
	drawPlayerHud(screen, gameState)
	drawInfo(screen, gameState, visibleEntities)
	drawGameInventory(screen, gameState)
}

func showDebug(screen *ebiten.Image, gs *gamestate.GameState) {
	w, _ := screen.Size()
	// Draw info
	fnt := assets.MplusFont(10)
	msg := fmt.Sprintf("TPS: %0.2f, FPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msg, fnt, w-150, 20, color.White)

	// Mouse info
	x, y := ebiten.CursorPosition()
	cursorStr := strconv.Itoa(x) + " " + strconv.Itoa(y)
	text.Draw(screen, cursorStr, fnt, w-150, 35, color.White)
	tx, ty := ToTile(gs, x, y)
	tileStr := fmt.Sprintf("Tile: %d,%d", tx, ty)
	text.Draw(screen, tileStr, fnt, w-150, 45, color.White)

}

func renderEntities(entities []*ecs.Entity, gameState *gamestate.GameState, screen *ebiten.Image) []*ecs.Entity {

	// font := assets.LoadFontCached(float64(18))
	font := assets.MplusFont(float64(18))

	visibleEntities := []*ecs.Entity{}

	pp := gameState.Player.GetComponent(state.Position).(state.PositionComponent)
	pStats := gameState.Player.GetComponent(state.Stats).(state.StatsComponent)
	lightColor := color.RGBA{
		R: 250,
		G: 250,
		B: 30,
		A: 0,
	}

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

		if isVisitable {
			// Walls and floor
			if visitable.Visible {
				bgColor := ParseHexColorFast(bg)
				if entity.HasComponent(state.IsBlocking) {
					distance := CalcDistance(position.X, position.Y, pp.X, pp.Y)
					mix := (float64(pStats.Fov) - float64(distance)) / float64(pStats.Fov)
					bgColor = ColorBlend(lightColor, bgColor, mix)
				}

				fgColor := ParseHexColorFast(fg)
				DrawChar(screen, gameState, position.X, position.Y, font, ch, fgColor, bgColor)
			} else if visitable.Explored {
				bgColor := ParseHexColorFast("#000000")
				fgColor := ParseHexColorFast("#555555")
				DrawChar(screen, gameState, position.X, position.Y, font, ch, fgColor, bgColor)
			}
		} else {
			if gameState.Fov.IsVisible(position.X, position.Y) {
				bgColor := ParseHexColorFast(bg)
				fgColor := ParseHexColorFast(fg)
				DrawChar(screen, gameState, position.X, position.Y, font, ch, fgColor, bgColor)

				visibleEntities = append(visibleEntities, entity)
			}
		}
	}
	return visibleEntities
}

func DrawTextRect(screen *ebiten.Image, str string, x int, y int, font font.Face, bgColor color.Color) {
	rect := text.BoundString(font, str)
	padL := 0
	padR := padL * 2
	ebitenutil.DrawRect(screen, float64(x+rect.Min.X-padL), float64(y+rect.Min.Y-padL), float64(rect.Max.X-rect.Min.X+padR), float64(rect.Max.Y-rect.Min.Y+padR), bgColor)

}

var messageLogColors = [5]color.Color{
	ParseHexColorFast("#333333"),
	ParseHexColorFast("#555555"),
	ParseHexColorFast("#777777"),
	ParseHexColorFast("#AAAAAA"),
	ParseHexColorFast("#FFFFFF"),
}

func drawMessageLog(screen *ebiten.Image, gs *gamestate.GameState) {

	fontSize := 14
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
		DrawText(screen, gs, position.X, position.Y+i, font, logStr, fgColor, color.Black)
	}
	DrawGridRect(screen, gs, position, color.White)

}

func drawPlayerHud(screen *ebiten.Image, gs *gamestate.GameState) {
	fontSize := 18
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.PlayerHud

	player := gs.Player
	stats, ok := player.GetComponent(state.Stats).(state.StatsComponent)
	if ok {
		DrawText(screen, gs, position.X, position.Y, font, stats.String(), ParseHexColorFast("#00AA00"), color.Black)
	}
	DrawGridRect(screen, gs, position, color.White)
}

func drawInfo(screen *ebiten.Image, gs *gamestate.GameState, visibleEntities []*ecs.Entity) {
	fontSize := 12
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.InfoBar

	y := position.Y
	for _, entity := range visibleEntities {
		if !entity.HasComponent(state.Player) {
			str := state.GetLongDescription(entity)
			DrawText(screen, gs, position.X, y, font, str, color.White, color.Black)
			y++
		}
	}
	DrawGridRect(screen, gs, position, color.White)
}

func drawGameInventory(screen *ebiten.Image, gs *gamestate.GameState) {
	fontSize := 12
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.GameInventory
	inventory, _ := gs.Player.GetComponent(state.Inventory).(state.InventoryComponent)

	y := position.Y
	status := fmt.Sprintf("Inventory %2d/%2d", len(inventory.Items), inventory.MaxItems)
	DrawText(screen, gs, position.X, y, font, status, color.White, color.Black)

	if len(inventory.Items) > 0 {
		for i, entity := range inventory.Items {
			str := fmt.Sprintf("%2d - %s", i+1, state.GetDescription(entity))
			DrawText(screen, gs, position.X, y+1, font, str, color.White, color.Black)
			y++
		}
	} else {
		DrawText(screen, gs, position.X, y+1, font, "- No items -", color.White, color.Black)
	}
	DrawGridRect(screen, gs, position, color.White)
}

var distances = map[string]int{}

func CalcDistance(x1 int, y1 int, x2 int, y2 int) int {
	key := strconv.Itoa(x1) + "-" + strconv.Itoa(y1) + "-" + strconv.Itoa(x2) + "-" + strconv.Itoa(y2)
	distance, ok := distances[key]
	if ok {
		return distance
	}

	x := float64(x2 - x1)
	y := float64(y2 - y1)
	dfloat := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
	distance = int(dfloat)
	distances[key] = distance
	return distance
}
