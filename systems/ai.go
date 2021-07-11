package systems

import (
	"math"
	"math/rand"

	"jordiburgos.com/officestruggle/astar"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/state"
)

type Tile struct {
	Ent   *ecs.Entity
	X     int
	Y     int
	tiles map[int]*Tile
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
func (t *Tile) GetNeighbors() []astar.Node {
	neighbors := []astar.Node{}
	eng := t.Ent.Engine

	for _, d := range DIRECTIONS {
		x := t.X + d.X
		y := t.Y + d.Y
		visitable, ok := eng.PosCache.GetOneByCoordAndComponents(x, y, []string{state.Visitable})
		if ok && !visitable.HasComponent(state.IsBlocking) {
			n := (t.tiles)[visitable.Id]
			neighbors = append(neighbors, n)

		}
	}
	return neighbors
}

func (t *Tile) D(to astar.Node) int {
	return 1
}

func (t *Tile) H(to astar.Node) int {
	posFrom := state.GetPosition(t.Ent)

	toTile := to.(*Tile)
	posTo := state.GetPosition(toTile.Ent)

	cost := math.Abs(float64(posTo.X-posFrom.X)) + math.Abs(float64(posTo.Y-posFrom.Y))
	return int(cost)
}

func AI(engine *ecs.Engine, gameState *game.GameState) {

	visitables := engine.Entities.GetEntities([]string{state.Visitable})
	tiles := map[int]*Tile{}
	for _, visitable := range visitables {
		if !visitable.HasComponent(state.IsBlocking) {
			pos := state.GetPosition(visitable)
			t := Tile{
				Ent:   visitable,
				X:     pos.X,
				Y:     pos.Y,
				tiles: tiles,
			}
			tiles[visitable.Id] = &t
		}
	}

	// Go to the tile where the Player is located
	player := engine.Entities.GetEntity([]string{state.Player})
	toTileEntity := getTileOfEntity(player)
	to := tiles[toTileEntity.Id]

	aiEntities := engine.Entities.GetEntities([]string{state.AI})
	for i, enemy := range aiEntities {
		if i > 0 {
			return
		}
		fromTileEntity := getTileOfEntity(enemy)
		from := tiles[fromTileEntity.Id]

		path, found := astar.AStar(from, to)
		if found {
			nextStep := (*path[1]).(*Tile)
			enemy.ReplaceComponent(state.Move, state.MoveComponent{X: nextStep.X, Y: nextStep.Y})
		}
	}

}

func getTileOfEntity(entity *ecs.Entity) *ecs.Entity {
	playerPos := state.GetPosition(entity)
	toTile, _ := entity.Engine.PosCache.GetOneByCoordAndComponents(playerPos.X, playerPos.Y, []string{state.Visitable})
	return toTile
}

func SimpleAI(engine *ecs.Engine, gameState *game.GameState) {
	aiEntities := engine.Entities.GetEntities([]string{state.AI})
	for _, enemy := range aiEntities {
		walkable := getWalkableNeighbor(enemy)
		selected := walkable[rand.Intn(len(walkable))]

		enemy.AddComponent(state.Move, state.MoveComponent{X: selected.X, Y: selected.Y})
	}
}

type Point struct {
	X int
	Y int
}

func getWalkableNeighbor(enemy *ecs.Entity) []Point {
	fromTile := getTileOfEntity(enemy)
	fromPos := state.GetPosition(fromTile)
	points := []Point{}
	for _, d := range DIRECTIONS {
		x := fromPos.X + d.X
		y := fromPos.Y + d.Y
		visitable, ok := enemy.Engine.PosCache.GetOneByCoordAndComponents(x, y, []string{state.Visitable})
		if ok && !visitable.HasComponent(state.IsBlocking) {
			point := Point{x, y}
			points = append(points, point)

		}
	}
	return points
}
