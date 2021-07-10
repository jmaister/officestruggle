package astar_test

import (
	"fmt"
	"math"
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/astar"
)

type Tile struct {
	X       int
	Y       int
	Blocked bool
	TheMap  TheMap
}

type TheMap struct {
	Tiles []Tile
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

	neighbors := []astar.Node{}

	for _, d := range DIRECTIONS {
		x := t.X + d.X
		y := t.Y + d.Y
		fmt.Println(x, ",", y)

		for _, n := range t.TheMap.Tiles {
			fmt.Println("trying tile", x, n.X, y, n.Y, n.X == x && n.Y == y)

			if n.X == x && n.Y == y && !n.Blocked {
				fmt.Println("found")
				neighbors = append(neighbors, &n)
			}
		}
	}

	fmt.Println("neighbors", t, len(neighbors), neighbors)
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

func TestAStar(t *testing.T) {

	theMap := TheMap{
		Tiles: []Tile{},
	}

	var from Tile
	var to Tile

	for x := 0; x < 5; x++ {
		for y := 0; y < 5; y++ {
			tile := Tile{x, y, false, theMap}
			fmt.Println("tile", tile.X, tile.Y)
			theMap.Tiles = append(theMap.Tiles, tile)
			if x == 1 && y == 1 {
				from = tile
			}
			if x == 4 && y == 4 {
				to = tile
			}
		}
	}
	fmt.Println("map", len(theMap.Tiles))
	fmt.Println("from", from.X, from.Y)
	fmt.Println("to", to.X, to.Y)

	path, found := astar.AStar(&from, &to)
	assert.Equal(t, true, found)
	assert.Equal(t, 4, len(path))
}
