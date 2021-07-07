package state

import (
	"math"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
)

type FieldOfVision struct {
	sinTable    map[int]float64
	cosTable    map[int]float64
	torchRadius int
}

func (f *FieldOfVision) initialize() {

	f.cosTable = make(map[int]float64)
	f.sinTable = make(map[int]float64)

	for i := 0; i < 360; i++ {
		ax := math.Sin(float64(i) / (float64(180) / math.Pi))
		ay := math.Cos(float64(i) / (float64(180) / math.Pi))

		f.sinTable[i] = ax
		f.cosTable[i] = ay
	}
}

func (f *FieldOfVision) SetTorchRadius(radius int) {
	if radius > 1 {
		if f.torchRadius != radius {
			f.torchRadius = radius
			f.initialize()
		}
	}
}

func SetVisibleComponent(engine *ecs.Engine, x int, y int, isVisible bool) (*ecs.Entity, bool) {
	visitableEntity, ok := engine.PosCache.GetOneByCoordAndComponents(x, y, []string{Visitable})
	if ok {
		visitable, _ := visitableEntity.RemoveComponent(Visitable).(VisitableComponent)
		visitable.Visible = isVisible
		visitableEntity.AddComponent(Visitable, visitable)
	}
	return visitableEntity, ok
}

func SetExploredComponent(engine *ecs.Engine, x int, y int, isExplored bool) (*ecs.Entity, bool) {
	visitableEntity, ok := engine.PosCache.GetOneByCoordAndComponents(x, y, []string{Visitable})
	if ok {
		visitable, _ := visitableEntity.RemoveComponent(Visitable).(VisitableComponent)
		visitable.Explored = isExplored
		visitableEntity.AddComponent(Visitable, visitable)
	}
	return visitableEntity, ok
}

func SetVisibleExploredComponent(engine *ecs.Engine, x int, y int, isVisible bool, isExplored bool) (*ecs.Entity, bool) {
	visitableEntity, ok := engine.PosCache.GetOneByCoordAndComponents(x, y, []string{Visitable})
	if ok {
		visitable, _ := visitableEntity.RemoveComponent(Visitable).(VisitableComponent)
		visitable.Visible = isVisible
		visitable.Explored = isExplored
		visitableEntity.AddComponent(Visitable, visitable)
	}
	return visitableEntity, ok
}

func SetVisibleEntities(entities ecs.EntityList, isVisible bool) {
	for _, e := range entities {
		visitable, _ := e.RemoveComponent(Visitable).(VisitableComponent)
		visitable.Visible = isVisible
		e.AddComponent(Visitable, visitable)
	}
}

func (f *FieldOfVision) RayCast(engine *ecs.Engine, playerX int, playerY int, gameMap *grid.Map) {
	// Cast out rays each degree in a 360 circle from the player. If a ray passes over a floor (does not block sight)
	// tile, keep going, up to the maximum torch radius (view radius) of the player. If the ray intersects a wall
	// (blocks sight), stop, as the player will not be able to see past that. Every visible tile will get the Visible
	// and Explored properties set to true.

	for i := 0; i < 360; i++ {

		ax := f.sinTable[i]
		ay := f.cosTable[i]

		x := float64(playerX)
		y := float64(playerY)

		// Mark the players current position as explored
		// gameMap.Tiles[playerX][playerY].Explored = true
		SetExploredComponent(engine, playerX, playerY, true)

		for j := 0; j < f.torchRadius; j++ {
			x -= ax
			y -= ay

			roundedX := int(Round(x))
			roundedY := int(Round(y))

			if x < float64(gameMap.X) || x > float64(gameMap.X+gameMap.Width) || y < float64(gameMap.Y) || y > float64(gameMap.Y+gameMap.Height) {
				// If the ray is cast outside of the map, stop
				break
			}

			//gameMap.Tiles[roundedX][roundedY].Explored = true
			//gameMap.Tiles[roundedX][roundedY].Visible = true
			found, ok := SetVisibleExploredComponent(engine, roundedX, roundedY, true, true)

			//if gameMap.Tiles[roundedX][roundedY].Blocks_sight == true {
			//	// The ray hit a wall, go no further
			//	break
			//}
			if ok {
				_, ok2 := found.GetComponent(IsBlocking).(IsBlockingComponent)
				if ok2 {
					break
				}
			}
		}
	}
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}
