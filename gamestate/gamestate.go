package gamestate

import (
	"log"

	"github.com/norendren/go-fov/fov"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
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
	_, ok := gs.Engine.PosCache.GetOneByCoordAndComponents(x, y, []string{constants.Visitable, constants.IsBlocking})
	return ok
}

// Effects with targets

type TargetingType string

const RandomAcquisitionType TargetingType = "random"
const SelectedAcquisitionType TargetingType = "selected"
const AreaAcquisitionType TargetingType = "area"

type DamageType string

const DamageSharedType DamageType = "shared"
const DamageEachType DamageType = "each"

type EffectFunction func(engine *ecs.Engine, gs *GameState, item *ecs.Entity, itemUser *ecs.Entity, targets ecs.EntityList)
