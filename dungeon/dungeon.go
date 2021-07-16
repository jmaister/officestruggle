package dungeon

import (
	"math"
	"math/rand"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type DungeonOptions struct {
	MinRoomSize  int
	MaxRoomSize  int
	MaxRoomCount int
}

func randBetween(min int, max int) int {
	return rand.Intn(max-min) + min
}

func CreateDungeon(engine *ecs.Engine, m grid.Rect, opts DungeonOptions) grid.Rectangle {

	dungeon, dungeonTiles := grid.GetRectangle(m.X, m.Y, m.Width, m.Height, false, grid.RectangleOptions{
		Sprite: grid.Wall,
	})

	var tiles map[string]grid.Tile = make(map[string]grid.Tile)
	var rooms []grid.Rectangle

	for _, tile := range dungeonTiles {
		tiles[tile.GetKey()] = tile
	}

	for i := 0; i < opts.MaxRoomCount; i++ {
		rw := randBetween(opts.MinRoomSize, opts.MaxRoomSize)
		rh := randBetween(opts.MinRoomSize, opts.MaxRoomSize)
		rx := randBetween(m.X, m.Width+m.X-rw)
		ry := randBetween(m.Y, m.Height+m.Y-rh)

		// Create a candidate room
		candidate, candidateTiles := grid.GetRectangle(rx, ry, rw, rh, true, grid.RectangleOptions{
			Sprite: grid.Floor,
		})

		// test if candidate is overlapping with any existing rooms
		existIntersection := false
		for r := 0; r < len(rooms); r++ {
			thisIntersects := grid.RectsIntersect(rooms[0], candidate)
			if thisIntersects {
				existIntersection = true
				break
			}
		}
		if !existIntersection {
			rooms = append(rooms, candidate)
			for _, tile := range candidateTiles {
				tiles[tile.GetKey()] = tile
			}
		}
	}

	for r := 1; r < len(rooms); r++ {
		prev := rooms[r-1].Center
		curr := rooms[r].Center

		for _, tile := range digHorizontalPassage(prev.X, curr.X, curr.Y) {
			tiles[tile.GetKey()] = tile
		}
		for _, tile := range digVerticalPassage(prev.Y, curr.Y, prev.X) {
			tiles[tile.GetKey()] = tile
		}
	}

	// TODO: move tiles to entities.go
	for _, tile := range tiles {
		tileEntity := engine.NewEntity()
		tileEntity.AddComponent(state.Layer100, state.Layer100Component{})
		tileEntity.AddComponent(state.Position, state.PositionComponent{X: tile.X, Y: tile.Y})
		tileEntity.AddComponent(state.Visitable, state.VisitableComponent{Explored: false, Visible: false})
		if tile.Sprite == grid.Wall {
			tileEntity.AddComponent(state.Description, state.DescriptionComponent{Name: "Wall"})
			tileEntity.AddComponent(state.IsBlocking, state.IsBlockingComponent{})
			tileEntity.AddComponent(state.Apparence, state.ApparenceComponent{Color: "#1a1aff", Bg: "#00FFFF", Char: '#'})
		} else if tile.Sprite == grid.Floor {
			tileEntity.AddComponent(state.Description, state.DescriptionComponent{Name: "Floor"})
			tileEntity.AddComponent(state.IsFloor, state.IsFloorComponent{})
			tileEntity.AddComponent(state.Apparence, state.ApparenceComponent{Color: "#e3e3e3", Char: '•'})
		} else {
			tileEntity.AddComponent(state.Apparence, state.ApparenceComponent{Color: "#FF66FF", Char: '.'})
		}
	}

	dungeon.Center = rooms[0].Center

	return dungeon
}

func digHorizontalPassage(x1 int, x2 int, y int) []grid.Tile {
	var tiles []grid.Tile
	start := math.Min(float64(x1), float64(x2))
	end := math.Max(float64(x1), float64(x2)) + 1
	x := start

	for x < end {
		tile := grid.Tile{
			X:      int(x),
			Y:      y,
			Sprite: grid.Floor,
		}
		tiles = append(tiles, tile)
		x++
	}

	return tiles
}

func digVerticalPassage(y1 int, y2 int, x int) []grid.Tile {
	var tiles []grid.Tile
	start := math.Min(float64(y1), float64(y2))
	end := math.Max(float64(y1), float64(y2)) + 1
	y := start

	for y < end {
		tile := grid.Tile{
			X:      x,
			Y:      int(y),
			Sprite: grid.Floor,
		}
		tiles = append(tiles, tile)
		y++
	}

	return tiles
}
