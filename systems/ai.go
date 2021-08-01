package systems

import (
	"math"
	"math/rand"
	"strconv"

	"jordiburgos.com/officestruggle/astar"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
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
		visitable, ok := eng.PosCache.GetOneByCoordAndComponents(x, y, []string{constants.Visitable})
		if ok && !visitable.HasComponent(constants.IsBlocking) {
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

func (t *Tile) String() string {
	return strconv.Itoa(t.X) + "," + strconv.Itoa(t.Y)
}

func AI(engine *ecs.Engine, gameState *gamestate.GameState) {
	visitables := engine.Entities.GetEntities([]string{constants.Visitable})
	tiles := map[int]*Tile{}
	for _, visitable := range visitables {
		if !visitable.HasComponent(constants.IsBlocking) {
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
	player := gameState.Player
	toTileEntity := getTileOfEntity(player)
	to := tiles[toTileEntity.Id]

	aiEntities := engine.Entities.GetEntities([]string{constants.AI})
	for _, enemy := range aiEntities {
		fromTileEntity := getTileOfEntity(enemy)
		from := tiles[fromTileEntity.Id]

		distance := from.H(to)
		stats := enemy.GetComponent(constants.Stats).(state.StatsComponent)

		if enemy.HasComponent(constants.Paralize) {
			current, _ := enemy.GetComponent(constants.Paralize).(state.ParalizeComponent)
			if current.TurnsLeft > 1 {
				enemy.ReplaceComponent(state.ParalizeComponent{
					TurnsLeft: current.TurnsLeft - 1,
				})
			} else {
				enemy.RemoveComponent(constants.Paralize)
			}
		} else if distance == 1 {
			// Attack to the player
			Attack(engine, gameState, enemy, []*ecs.Entity{player})
		} else if distance > stats.Fov {
			// Wander
			wander(enemy)
		} else {
			// Follow the player
			path, found := astar.AStar(from, to)
			if found && len(path) > 0 {
				currStep := (*path[0]).(*Tile)
				nextStep := (*path[1]).(*Tile)
				dx := nextStep.X - currStep.X
				dy := nextStep.Y - currStep.Y
				enemy.ReplaceComponent(state.MoveComponent{X: dx, Y: dy})
			}
		}

	}

}

func getTileOfEntity(entity *ecs.Entity) *ecs.Entity {
	playerPos := state.GetPosition(entity)
	toTile, _ := entity.Engine.PosCache.GetOneByCoordAndComponents(playerPos.X, playerPos.Y, []string{constants.Visitable})
	return toTile
}

func SimpleAI(engine *ecs.Engine, gameState *gamestate.GameState) {
	aiEntities := engine.Entities.GetEntities([]string{constants.AI})
	for _, enemy := range aiEntities {
		wander(enemy)
	}
}

func wander(entity *ecs.Entity) {
	walkable := getWalkableNeighbors(entity)
	selected := walkable[rand.Intn(len(walkable))]

	entity.AddComponent(state.MoveComponent{X: selected.X, Y: selected.Y})
}

func getWalkableNeighbors(enemy *ecs.Entity) []Dir {
	fromTile := getTileOfEntity(enemy)
	fromPos := state.GetPosition(fromTile)
	points := []Dir{}
	for _, d := range DIRECTIONS {
		x := fromPos.X + d.X
		y := fromPos.Y + d.Y
		visitable, ok := enemy.Engine.PosCache.GetOneByCoordAndComponents(x, y, []string{constants.Visitable})
		if ok && !visitable.HasComponent(constants.IsBlocking) {
			point := Dir{d.X, d.Y}
			points = append(points, point)

		}
	}
	return points
}
