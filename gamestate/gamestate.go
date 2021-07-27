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

type ScreenState string

const (
	WelcomeScreen   ScreenState = "welcome"
	GameScreen      ScreenState = "game"
	TargetingScreen ScreenState = "target"
	InventoryScreen ScreenState = "inventory"
	TestScreen      ScreenState = "test"
)

type InventoryScreenFocus string

const (
	InventoryFocus InventoryScreenFocus = "i"
	EquipmentFocus InventoryScreenFocus = "e"
)

type ListState struct {
	Selected  int
	IsFocused bool
}

type InventoryScreenState struct {
	InventoryState ListState
	EquipmentState ListState
}

type GameState struct {
	Engine               *ecs.Engine
	Fov                  *fov.View
	Grid                 *grid.Grid
	Player               *ecs.Entity
	ScreenState          ScreenState
	InventoryScreenState InventoryScreenState
	IsPlayerTurn         bool
	L                    *log.Logger
	ScreenWidth          int
	ScreenHeight         int
	TileWidth            int
	TileHeight           int
	logLines             []LogLine
}

type LogType string

const (
	Info   LogType = "i"
	Warn   LogType = "w"
	Bad    LogType = "b"
	Danger LogType = "d"
	Good   LogType = "g"
)

type LogLine struct {
	Msg   string
	Count int
	Type  LogType
}

func NewGameState(engine *ecs.Engine) *GameState {

	// Dungeon
	g := grid.Grid{
		Width:  80,
		Height: 34,
		Map: grid.Rect{
			X:      16,
			Y:      8,
			Width:  74,
			Height: 30,
		},
		MessageLog: grid.Rect{
			X:      0,
			Y:      0,
			Width:  79,
			Height: 5,
		},
		PlayerHud: grid.Rect{
			X:      0,
			Y:      6,
			Width:  79,
			Height: 1,
		},
		InfoBar: grid.Rect{
			X:      16,
			Y:      39,
			Width:  79,
			Height: 10,
		},
		GameInventory: grid.Rect{
			X:      0,
			Y:      8,
			Width:  15,
			Height: 40,
		},
		Inventory: grid.Rect{
			X:      2,
			Y:      5,
			Width:  79,
			Height: 29,
		},
		Equipment: grid.Rect{
			X:      50,
			Y:      5,
			Width:  20,
			Height: 29,
		},
	}
	dungeonRectangle := dungeon.CreateDungeon(engine, g.Map, dungeon.DungeonOptions{
		MinRoomSize:  6,
		MaxRoomSize:  12,
		MaxRoomCount: 40,
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
	for i := 0; i < 10; i++ {
		v := visitables[rand.Intn(len(visitables))]
		pos := state.GetPosition(v)
		potion := state.NewHealthPotion(engine.NewEntity())
		state.ApplyPosition(potion, pos.X, pos.Y)
	}
	// Swords
	for i := 0; i < 10; i++ {
		v := visitables[rand.Intn(len(visitables))]
		pos := state.GetPosition(v)
		potion := state.NewSword(engine.NewEntity())
		state.ApplyPosition(potion, pos.X, pos.Y)
	}

	return &GameState{
		Engine:      engine,
		Fov:         fov.New(),
		Grid:        &g,
		Player:      player,
		ScreenState: WelcomeScreen,
		InventoryScreenState: InventoryScreenState{
			InventoryState: ListState{
				Selected:  0,
				IsFocused: true,
			},
			EquipmentState: ListState{
				Selected:  0,
				IsFocused: false,
			},
		},
		IsPlayerTurn: true,
		L:            log.New(os.Stderr, "", 0),
		ScreenWidth:  1920,
		ScreenHeight: 1080,
		TileWidth:    20,
		TileHeight:   20,
	}
}

func (gs *GameState) Log(t LogType, s string) {
	n := len(gs.logLines)
	if n > 0 {
		if gs.logLines[n-1].Msg == s {
			gs.logLines[n-1].Count++
			return
		}
	}
	gs.logLines = append(gs.logLines, LogLine{Msg: s, Count: 1, Type: t})
}

func (gs *GameState) GetLog(lineNumber int) []LogLine {
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
