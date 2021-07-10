package systems

import (
	"fmt"
	"math"

	"github.com/beefsack/go-astar"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/state"
)

type Tile struct {
	Ent   *ecs.Entity
	X     int
	Y     int
	tiles *map[int]Tile
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
	neighbors := []astar.Pather{}
	eng := t.Ent.Engine

	for _, d := range DIRECTIONS {
		x := t.X + d.X
		y := t.Y + d.Y
		visitable, ok := eng.PosCache.GetOneByCoordAndComponents(x, y, []string{state.Visitable})
		if ok {
			if !visitable.HasComponent(state.IsBlocking) {
				n := (*t.tiles)[visitable.Id]
				neighbors = append(neighbors, &n)
			}
		}
	}
	return neighbors
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1.0
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	posFrom := state.GetPosition(t.Ent)

	toTile := to.(*Tile)
	posTo := state.GetPosition(toTile.Ent)

	cost := math.Abs(float64(posTo.X-posFrom.X)) + math.Abs(float64(posTo.Y-posFrom.Y))
	return cost
}

func AI(engine *ecs.Engine, gameState *game.GameState) {

	visitables := engine.Entities.GetEntities([]string{state.Visitable})
	tiles := map[int]Tile{}
	for _, visitable := range visitables {
		if !visitable.HasComponent(state.IsBlocking) {
			pos := state.GetPosition(visitable)
			t := Tile{
				Ent:   visitable,
				X:     pos.X,
				Y:     pos.Y,
				tiles: &tiles,
			}
			tiles[visitable.Id] = t
		}
	}

	// Go to the tile where the Player is located
	player := engine.Entities.GetEntity([]string{state.Player})
	toTile := getTileOfEntity(player)
	to := tiles[toTile.Id]

	aiEntities := engine.Entities.GetEntities([]string{state.AI})
	for i, enemy := range aiEntities {
		if i > 0 {
			break
		}
		fromTile := getTileOfEntity(enemy)
		from := tiles[fromTile.Id]

		path, distance, found := astar.Path(&from, &to)
		fmt.Println(path, distance, found)
	}

}

func getTileOfEntity(entity *ecs.Entity) *ecs.Entity {
	playerPos := state.GetPosition(entity)
	toTile, _ := entity.Engine.PosCache.GetOneByCoordAndComponents(playerPos.X, playerPos.Y, []string{state.Visitable})
	return toTile
}
