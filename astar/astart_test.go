package astar_test

import (
	"math"
	"strconv"
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/astar"
)

type Tile struct {
	X       int
	Y       int
	Blocked bool
	TheMap  *TheMap
}

type TheMap struct {
	Tiles []*Tile
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

func (t *Tile) GetNeighbors() []astar.Node {
	//	neighbors := make([]astar.Node, len(DIRECTIONS))
	neighbors := []astar.Node{}

	for _, d := range DIRECTIONS {
		x := t.X + d.X
		y := t.Y + d.Y
		for _, n := range t.TheMap.Tiles {
			if n.X == x && n.Y == y && !n.Blocked {
				neighbors = append(neighbors, n)
			}
		}
	}
	return neighbors
}

func (t *Tile) H(to astar.Node) int {
	toTile := to.(*Tile)
	cost := math.Abs(float64(t.X-toTile.X)) + math.Abs(float64(t.Y-toTile.Y))
	return int(cost)
}
func (t *Tile) D(neighbor astar.Node) int {
	return 1
}

func (t *Tile) String() string {
	return strconv.Itoa(t.X) + "," + strconv.Itoa(t.Y)
}

func TestAStar(t *testing.T) {

	theMap := TheMap{
		Tiles: []*Tile{},
	}

	var from Tile
	var to Tile

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			tile := Tile{x, y, false, &theMap}
			theMap.Tiles = append(theMap.Tiles, &tile)
			if x == 1 && y == 1 {
				from = tile
			}
			if x == 4 && y == 4 {
				to = tile
			}
		}
	}

	path, found := astar.AStar(&from, &to)
	assert.Equal(t, true, found)
	assert.Equal(t, 7, len(path))
}

func TestAStarWithBlock(t *testing.T) {

	theMap := TheMap{
		Tiles: []*Tile{},
	}

	var from Tile
	var to Tile

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			blocked := false
			if x == y && x > 2 && x < 4 {
				blocked = true
			}
			if x == 2 && y == 3 {
				blocked = true
			}
			tile := Tile{x, y, blocked, &theMap}
			theMap.Tiles = append(theMap.Tiles, &tile)
			if x == 1 && y == 1 {
				from = tile
			}
			if x == 4 && y == 4 {
				to = tile
			}
		}
	}

	path, found := astar.AStar(&from, &to)
	assert.Equal(t, true, found)
	assert.Equal(t, 7, len(path))
}

func constructDiagonalMap(side int) (TheMap, Tile, Tile) {
	theMap := TheMap{
		Tiles: []*Tile{},
	}

	for x := 0; x < side; x++ {
		for y := 0; y < side; y++ {
			blocked := false
			if x+y == side && x >= 1 && y > 1 && x < side-1 && y < side-1 {
				blocked = true
			}

			tile := Tile{x, y, blocked, &theMap}
			theMap.Tiles = append(theMap.Tiles, &tile)
		}
	}

	from := Tile{
		X:      0,
		Y:      0,
		TheMap: &theMap,
	}
	to := Tile{
		X:      side - 1,
		Y:      side - 1,
		TheMap: &theMap,
	}

	return theMap, from, to
}

func TestAStarBig(t *testing.T) {

	_, from, to := constructDiagonalMap(50)

	path, found := astar.AStar(&from, &to)
	assert.Equal(t, true, found)
	assert.Equal(t, 99, len(path))
}

func calculate() {
	_, from, to := constructDiagonalMap(50)
	astar.AStar(&from, &to)
}

func BenchmarkAstar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculate()
	}
}
