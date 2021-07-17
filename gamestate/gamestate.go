package gamestate

import (
	"log"
	"math/rand"
	"os"

	"github.com/norendren/go-fov/fov"
	"jordiburgos.com/officestruggle/dungeon"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type GameState struct {
	Engine       *ecs.Engine
	Fov          *fov.View
	Grid         *grid.Grid
	Player       *ecs.Entity
	IsPlayerTurn bool
	L            *log.Logger
	ScreenWidth  int
	ScreenHeight int
	TileWidth    int
	TileHeight   int
	logLines     []string
}

func NewGameState(engine *ecs.Engine) *GameState {

	// Dungeon
	g := grid.Grid{
		Width:  80,
		Height: 34,
		Map: grid.Rect{
			X:      2,
			Y:      6,
			Width:  79,
			Height: 29,
		},
		MessageLog: grid.Rect{
			Width:  79,
			Height: 5,
			X:      0,
			Y:      0,
		},
	}
	dungeonRectangle := dungeon.CreateDungeon(engine, g.Map, dungeon.DungeonOptions{
		MinRoomSize:  6,
		MaxRoomSize:  12,
		MaxRoomCount: 8,
	})

	// Player
	player := state.NewPlayer(engine.NewEntity())
	state.ApplyPosition(player, dungeonRectangle.Center.X, dungeonRectangle.Center.Y)

	visitables := engine.Entities.GetEntities([]string{state.IsFloor})
	// Enemies
	for i := 0; i < 5; i++ {
		v := visitables[rand.Intn(len(visitables))]
		pos := state.GetPosition(v)
		goblin := state.NewGlobin(engine.NewEntity())
		state.ApplyPosition(goblin, pos.X, pos.Y)
	}
	// Health potions
	for i := 0; i < 5; i++ {
		v := visitables[rand.Intn(len(visitables))]
		pos := state.GetPosition(v)
		potion := state.NewHealthPotion(engine.NewEntity())
		state.ApplyPosition(potion, pos.X, pos.Y)
	}

	return &GameState{
		Engine:       engine,
		Fov:          fov.New(),
		Grid:         &g,
		Player:       player,
		IsPlayerTurn: true,
		L:            log.New(os.Stderr, "", 0),
		ScreenWidth:  1024,
		ScreenHeight: 576,
		TileWidth:    16,
		TileHeight:   16,
	}
}

func (gs *GameState) Log(s string) {
	gs.logLines = append(gs.logLines, s)
}

func (gs *GameState) GetLog(lineNumber int) []string {
	n := len(gs.logLines)
	if lineNumber <= n {
		return gs.logLines[n-lineNumber : n]
	}
	return gs.logLines
}

// Implement the GridMap interface for the Fov.

func (gs *GameState) InBounds(x int, y int) bool {
	gm := gs.Grid.Map
	return x >= gm.X || x <= gm.X+gm.Width || y >= gm.Y || y <= gm.Y+gm.Height
}

func (gs *GameState) IsOpaque(x int, y int) bool {
	visitableEntity, ok := gs.Engine.PosCache.GetOneByCoordAndComponents(x, y, []string{state.Visitable})
	if ok {
		_, ok2 := visitableEntity.GetComponent(state.IsBlocking).(state.IsBlockingComponent)
		return ok2
	}
	return false
}
