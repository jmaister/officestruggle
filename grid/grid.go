package grid

import (
	"math"
	"math/rand"
	"strconv"
)

type Direction struct {
	X int
	Y int
}

var UP Direction = Direction{X: 0, Y: -1}
var RIGHT Direction = Direction{X: 1, Y: 0}
var DOWN Direction = Direction{X: 0, Y: 1}
var LEFT Direction = Direction{X: -1, Y: 0}

var UP_RIGHT Direction = Direction{X: 1, Y: -1}
var DOWN_RIGHT Direction = Direction{X: 1, Y: 1}
var DOWN_LEFT Direction = Direction{X: -1, Y: 1}
var UP_LEFT Direction = Direction{X: -1, Y: -1}

var Cardinal = []Direction{
	UP,    // N
	RIGHT, // E
	DOWN,  // S
	LEFT,  // W
}

var Diagonal = []Direction{
	UP_RIGHT,   // NE
	DOWN_RIGHT, // SE
	DOWN_LEFT,  // SW
	UP_LEFT,    // NW
}

var AllDirection = append(Cardinal, Diagonal...)

type Tile struct {
	X      int
	Y      int
	Z      int
	Sprite TileType
}

func (t Tile) GetKey() string {
	return strconv.Itoa(t.X) + "," + strconv.Itoa(t.Y) + "," + strconv.Itoa(t.Z)
}
func (t Tile) String() string {
	return t.GetKey()
}

func IsInsideCircle(center Tile, tile Tile, radius int) bool {
	dx := center.X - tile.X
	dy := center.Y - tile.Y
	distance_squared := dx*dx + dy*dy
	return distance_squared <= radius*radius
}

func GetCircle(center Tile, radius int) []Tile {
	diameter := radius * 2
	if radius%1 != 0 {
		diameter += 1
	}
	top := center.Y - radius
	bottom := center.Y + radius
	left := center.X - radius
	right := center.X + radius

	var tiles []Tile

	for y := top; y <= bottom; y++ {
		for x := left; x <= right; x++ {
			thisTile := Tile{X: x, Y: y}
			if IsInsideCircle(center, thisTile, radius) {
				tiles = append(tiles, thisTile)

			}
		}
	}
	return tiles

}

type TileType string

const (
	Wall       TileType = "Wall"
	Floor      TileType = "Floor"
	Upstairs   TileType = "Upstairs"
	Downstairs TileType = "Downstairs"
	Undefined  TileType = "Undef"
)

type RectangleOptions struct {
	Sprite TileType
	Z      int
}

func GetRectangle(x int, y int, width int, height int, hasWalls bool, opts RectangleOptions) (Rectangle, []Tile) {
	var tiles []Tile

	x1 := x
	x2 := x + width - 1
	y1 := y
	y2 := y + height - 1

	if hasWalls {
		for yi := y1 + 1; yi <= y2-1; yi++ {
			for xi := x1 + 1; xi <= x2-1; xi++ {
				thisTile := Tile{X: xi, Y: yi, Z: opts.Z, Sprite: opts.Sprite}
				tiles = append(tiles, thisTile)
			}
		}
	} else {
		for yi := y1; yi <= y2; yi++ {
			for xi := x1; xi <= x2; xi++ {
				thisTile := Tile{X: xi, Y: yi, Z: opts.Z, Sprite: opts.Sprite}
				tiles = append(tiles, thisTile)
			}
		}
	}

	center := Tile{
		X: int(math.Round(float64(x1+x2) / 2.0)),
		Y: int(math.Round(float64(y1+y2) / 2.0)),
		Z: opts.Z,
	}

	return Rectangle{x1, x2, y1, y2, center}, tiles
}

type Rectangle struct {
	X1     int
	X2     int
	Y1     int
	Y2     int
	Center Tile
}

func RectsIntersect(r1 Rectangle, r2 Rectangle) bool {
	return r1.X1 <= r2.X2 &&
		r1.X2 >= r2.X1 &&
		r1.Y1 <= r2.X2 &&
		r1.Y2 >= r2.Y1
}

func Distance(t1 Tile, t2 Tile) int {
	x := math.Pow(float64(t2.X-t1.X), 2)
	y := math.Pow(float64(t2.Y-t1.Y), 2)
	return int(math.Floor(math.Sqrt(x + y)))
}

type Rect struct {
	Width  int
	Height int
	X      int
	Y      int
}

type Grid struct {
	Width         int
	Height        int
	Levels        int
	Map           Rect
	Camera        Rect
	MessageLog    Rect
	PlayerHud     Rect
	InfoBar       Rect
	GameInventory Rect
	Inventory     Rect
	Equipment     Rect
	Money         Rect
}

func IsOnMapEdge(x int, y int, rect Rect) bool {

	if x == rect.X {
		return true
	}
	if y == rect.Y {
		return true
	}
	if x == rect.X+rect.Width-1 {
		return true
	}
	if y == rect.Y+rect.Height-1 {
		return true
	}
	return false
}

func GetNeighbors(x int, y int, z int, directions []Direction, grid Grid) []Tile {
	var tiles []Tile

	for _, dir := range directions {
		candidate := Tile{X: x + dir.X, Y: y + dir.Y, Z: z}
		if candidate.X >= 0 &&
			candidate.X < grid.Width &&
			candidate.Y >= 0 &&
			candidate.Y < grid.Height {
			tiles = append(tiles, candidate)
		}
	}
	return tiles
}

func IsNeighbor(a Tile, b Tile) bool {
	return (a.X-b.X == 1 && a.Y-b.Y == 0) ||
		(a.X-b.X == 0 && a.Y-b.Y == -1) ||
		(a.X-b.X == -1 && a.Y-b.Y == 0) ||
		(a.X-b.X == 0 && a.Y-b.Y == 1)

}

func RandomNeighbor(start Tile) Tile {
	size := len(Cardinal)
	dir := Cardinal[rand.Int31n(int32(size))]
	return Tile{
		X: start.X + dir.X,
		Y: start.Y + dir.Y,
		Z: start.Z,
	}
}
