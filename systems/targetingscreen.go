package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func RenderTargetingScreen(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	player := gameState.Player
	plPosition := player.GetComponent(state.Position).(state.PositionComponent)
	stats := player.GetComponent(state.Stats).(state.StatsComponent)

	fov := stats.Fov

	bg := color.RGBA{
		R: 255,
		G: 232,
		B: 199,
		A: 127, // 0.5
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Calculate line with tile positions
	targetX, targetY := ToTile(gameState, mouseX, mouseY)
	line := BresenhamLine(plPosition.X, plPosition.Y, targetX, targetY)
	for _, tile := range line {
		if CalcDistance(plPosition.X, plPosition.Y, tile.X, tile.Y) >= fov {
			break
		}
		_, blocked := engine.PosCache.GetOneByCoordAndComponents(tile.X, tile.Y, []string{state.IsBlocking})
		DrawTile(screen, gameState, tile.X, tile.Y, bg)
		if blocked {
			break
		}
	}
}

func TargetingMouseClick(engine *ecs.Engine, gameState *gamestate.GameState, mouseX int, mouseY int) {

	x, y := ToTile(gameState, mouseX, mouseY)

	targetEntities, ok := engine.PosCache.GetByCoord(x, y)
	if ok {
		fmt.Println("targets", targetEntities)
	} else {
		fmt.Println("No targets found")
	}
}

// Returns the list of points from (x0, y0) to (x1, y1).
// https://www.codeproject.com/Articles/15604/Ray-casting-in-a-2D-tile-based-environment
func BresenhamLine(x0 int, y0 int, x1 int, y1 int) []grid.Tile {
	// Optimization: it would be preferable to calculate in
	// advance the size of "result" and to use a fixed-size array
	// instead of a list.
	result := []grid.Tile{}

	steep := math.Abs(float64(y1-y0)) > math.Abs(float64(x1-x0))
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	deltax := x1 - x0
	deltay := int(math.Abs(float64(y1 - y0)))
	error := 0
	var ystep int
	y := y0

	if y0 < y1 {
		ystep = 1
	} else {
		ystep = -1
	}
	for x := x0; x <= x1; x++ {
		if steep {
			result = append(result, grid.Tile{X: y, Y: x})
		} else {
			result = append(result, grid.Tile{X: x, Y: y})
		}
		error += deltay
		if 2*error >= deltax {
			y += ystep
			error -= deltax
		}
	}

	return result
}
