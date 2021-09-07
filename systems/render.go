package systems

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/palette"
	"jordiburgos.com/officestruggle/state"
)

var fnt20 = assets.LoadFontCached(float64(20))
var fnt40 = assets.LoadFontCached(float64(40))
var fnt = assets.MplusFont(18)

func setVisibleEntities(entities ecs.EntityList, isVisible bool) {
	for _, e := range entities {
		visitable, _ := e.GetComponent(constants.Visitable).(state.VisitableComponent)
		visitable.Visible = isVisible
		e.ReplaceComponent(visitable)
	}
}

func Render(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	showDebug(screen, gameState)

	// Order of drawing the layers
	layers := []string{constants.Layer100, constants.Layer300, constants.Layer400, constants.Layer500}

	// Reset visibility
	visitables := engine.Entities.GetEntitiesWithFilter(func(e *ecs.Entity) bool {
		if e.HasComponent(constants.Visitable) {
			pos, ok := e.GetComponent(constants.Position).(state.PositionComponent)
			return ok && pos.Z == gameState.CurrentZ
		}
		return false
	})
	setVisibleEntities(visitables, false)

	// Update visibility based on player's FOV
	for _, visitable := range visitables {
		pos := state.GetPosition(visitable)
		if gameState.Fov.IsVisible(pos.X, pos.Y) {
			vsComponent, _ := visitable.GetComponent(constants.Visitable).(state.VisitableComponent)
			vsComponent.Visible = true
			vsComponent.Explored = true
			visitable.ReplaceComponent(vsComponent)
		}
	}

	visibleEntities := []*ecs.Entity{}
	for _, layer := range layers {
		renderable := []string{constants.Position, constants.Apparence, layer}
		entities := engine.Entities.GetEntities(renderable)
		entities = FilterZ(entities, gameState.CurrentZ)

		v := renderEntities(entities, gameState, screen)

		visibleEntities = append(visibleEntities, v...)
	}
	DrawGridRect(screen, gameState, gameState.Grid.Camera, color.White)

	drawMessageLog(screen, gameState)
	drawPlayerHud(screen, gameState)
	drawInfo(screen, engine, gameState, visibleEntities)
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

	visibleEntities := []*ecs.Entity{}

	plPos := gameState.Player.GetComponent(constants.Position).(state.PositionComponent)
	plStats := gameState.Player.GetComponent(constants.Stats).(state.StatsComponent)
	lightColor := palette.PColor(palette.Yellow, 0.5)

	for _, entity := range entities {
		position, _ := entity.GetComponent(constants.Position).(state.PositionComponent)
		apparence, _ := entity.GetComponent(constants.Apparence).(state.ApparenceComponent)
		visitable, isVisitable := entity.GetComponent(constants.Visitable).(state.VisitableComponent)

		fg := apparence.Color
		if fg == nil {
			fg = color.White
		}
		bg := apparence.Bg
		if bg == nil {
			bg = color.Black
		}
		ch := string(apparence.Char)

		shouldDraw := false
		if isVisitable {
			// Walls and floor
			if visitable.Visible {
				if entity.HasComponent(constants.IsBlocking) {
					distance := CalcDistance(position.X, position.Y, plPos.X, plPos.Y)
					mix := (float64(plStats.Fov) - float64(distance)) / float64(plStats.Fov)
					bg = ColorBlend(lightColor, bg, mix)
				}
				shouldDraw = true
			} else if visitable.Explored {
				bg = color.Black
				fg = palette.PColor(palette.Gray, 0.3)
				shouldDraw = true
			}
		} else {
			if gameState.Fov.IsVisible(position.X, position.Y) {
				shouldDraw = true

				visibleEntities = append(visibleEntities, entity)
			}
		}

		if shouldDraw {
			cameraX, cameraY := gameState.Camera.ToCameraCoordinates(position.X, position.Y)
			cameraX += gameState.Grid.Camera.X
			cameraY += gameState.Grid.Camera.Y
			DrawChar(screen, gameState, cameraX, cameraY, fnt, ch, fg, bg)
		}
	}
	return visibleEntities
}

func drawMessageLog(screen *ebiten.Image, gs *gamestate.GameState) {

	position := gs.Grid.MessageLog

	lines := gs.GetLogLines(position.Height)
	for i, line := range lines {
		fgColor := constants.LogColors[line.Type]
		logStr := line.Msg
		if line.Count > 1 {
			logStr = strconv.Itoa(line.Count) + "x " + line.Msg
		}
		DrawText(screen, gs, position.X, position.Y+i, fnt, logStr, fgColor, color.Black)
	}
	DrawGridRect(screen, gs, position, color.White)

}

func drawPlayerHud(screen *ebiten.Image, gs *gamestate.GameState) {

	msgColor := palette.PColor(palette.Lime, 0.6)

	position := gs.Grid.PlayerHud

	player := gs.Player
	stats, _ := player.GetComponent(constants.Stats).(state.StatsComponent)
	floorStr := fmt.Sprintf("%d of %d", gs.CurrentZ+1, gs.Grid.Levels)

	strings := []ColoredText{
		{"Player:", color.White, color.Black},
		{stats.String(), msgColor, color.Black},
		{"Floor: ", color.White, color.Black},
		{floorStr, msgColor, color.Black},
	}
	DrawStrings(screen, gs, position.X, position.Y, fnt, strings)

	lvl, _ := player.GetComponent(constants.Leveling).(state.LevelingComponent)
	strings2 := []ColoredText{
		{"Lvl:", color.White, color.Black},
		{strconv.Itoa(lvl.CurrentLevel), msgColor, color.Black},
		{" XP:", color.White, color.Black},
		{strconv.Itoa(lvl.CurrentXP), msgColor, color.Black},
		{" of ", color.White, color.Black},
		{strconv.Itoa(lvl.GetNextLevelXP()), msgColor, color.Black},
	}
	DrawStrings(screen, gs, position.X, position.Y+1, fnt, strings2)

	DrawGridRect(screen, gs, position, color.White)
}

func drawInfo(screen *ebiten.Image, engine *ecs.Engine, gs *gamestate.GameState, visibleEntities ecs.EntityList) {
	fontSize := 12
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.InfoBar

	plPos := state.GetPosition(gs.Player)
	playerId := gs.Player.Id

	y := position.Y
	for _, entity := range visibleEntities {
		if entity.Id == playerId {
			continue
		}
		entPos := state.GetPosition(entity)
		isOnPlayerTile := plPos.X == entPos.X && plPos.Y == entPos.Y
		var cl color.Color = color.White
		onSameTile := " "
		if isOnPlayerTile {
			cl = palette.PColor(palette.Green, 0.5)
			onSameTile = "*"
		}

		// TODO: change colors for Equipable based on tiers
		// TODO: print entity char and color
		str := fmt.Sprintf("%s%s", onSameTile, state.GetLongDescription(entity))
		DrawText(screen, gs, position.X, y, font, str, cl, color.Black)
		y++
	}
	DrawGridRect(screen, gs, position, color.White)
}

func drawGameInventory(screen *ebiten.Image, gs *gamestate.GameState) {
	fontSize := 12
	font := assets.MplusFont(float64(fontSize))

	position := gs.Grid.GameInventory
	inventory, _ := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)

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
