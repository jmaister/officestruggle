package systems

import (
	"fmt"
	"math"

	"github.com/beefsack/go-astar"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/state"
)

type Tile struct {
	E *ecs.Entity
	X int
	Y int
}

type Dir struct {
	X int
	Y int
}

var DIRECTIONS = []Dir{
	{0, 1},  // right
	{0, -1}, // left
	{-1, 0}, // up
	{1, 0},  // dow

}

// Entity implementation for a-star
func (t *Tile) PathNeighbors() []astar.Pather {
	eng := t.E.Engine
	neighbors := []astar.Pather{}

	for _, d := range DIRECTIONS {
		x := t.X + d.X
		y := t.Y + d.Y
		visitable, ok := eng.PosCache.GetOneByCoordAndComponents(x, y, []string{state.Visitable})
		fmt.Println("visitable", visitable, ok)
		if ok && !visitable.HasComponent(state.IsBlocking) {
			n := &Tile{
				E: visitable,
				X: x,
				Y: y,
			}
			neighbors = append(neighbors, n)
		}
	}
	fmt.Println("returing neigbors", len(neighbors))
	return neighbors
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	fmt.Println()
	return 1
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	posFrom := state.GetPosition(t.E)

	toTile := to.(*Tile)
	posTo := state.GetPosition(toTile.E)

	return math.Abs(float64(posTo.X-posFrom.X)) + math.Abs(float64(posTo.Y-posFrom.Y))
}

func AI(engine *ecs.Engine, gameState *state.GameState) {

	aiEntities := engine.Entities.GetEntities([]string{state.AI})
	fmt.Println("enemies", len(aiEntities))
	for i, enemy := range aiEntities {
		if i > 0 {
			return
		}
		enemyPos := state.GetPosition(enemy)
		from := Tile{
			enemy, enemyPos.X, enemyPos.Y,
		}

		player := engine.Entities.GetEntity([]string{state.Player})
		playerPos := state.GetPosition(player)
		to := Tile{
			player, playerPos.X, playerPos.Y,
		}

		path, distance, found := astar.Path(&from, &to)
		fmt.Println(path, distance, found)
	}

}
