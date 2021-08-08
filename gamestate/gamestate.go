package gamestate

import (
	"image/color"
	"log"

	"github.com/norendren/go-fov/fov"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
)

const GameVersion = "1.0.0"
const SaveGamePrefix = "officestruggle"

type ScreenState string

const (
	WelcomeScreen   ScreenState = "welcome"
	GameScreen      ScreenState = "game"
	TargetingScreen ScreenState = "target"
	InventoryScreen ScreenState = "inventory"
	LoadingScreen   ScreenState = "loading"
	GameoverScreen  ScreenState = "gameover"
	ActionDialog    ScreenState = "actiondialog"
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

type ActionScreenState struct {
	Actions ListState
	Items   ecs.EntityList
}

type GameState struct {
	Engine               *ecs.Engine
	Fov                  *fov.View
	Grid                 *grid.Grid
	Player               *ecs.Entity
	CurrentZ             int
	ScreenState          ScreenState
	InventoryScreenState InventoryScreenState
	ActionScreenState    ActionScreenState
	IsPlayerTurn         bool
	L                    *log.Logger
	ScreenWidth          int
	ScreenHeight         int
	TileWidth            int
	TileHeight           int
	LogLines             []LogLine
}

type LogLine struct {
	Msg   string
	Count int
	Type  constants.LogType
}

func (gs *GameState) Log(t constants.LogType, s string) {
	n := len(gs.LogLines)
	if n > 0 {
		if gs.LogLines[n-1].Msg == s {
			gs.LogLines[n-1].Count++
			return
		}
	}
	gs.LogLines = append(gs.LogLines, LogLine{Msg: s, Count: 1, Type: t})
}

func (gs *GameState) GetLog(lineNumber int) []LogLine {
	n := len(gs.LogLines)
	if lineNumber <= n {
		return gs.LogLines[n-lineNumber : n]
	}
	return gs.LogLines
}

// Implement the GridMap interface for the Fov.

func (gs *GameState) InBounds(x int, y int) bool {
	gm := gs.Grid.Map
	return x >= gm.X || x <= gm.X+gm.Width || y >= gm.Y || y <= gm.Y+gm.Height
}

func (gs *GameState) IsOpaque(x int, y int) bool {
	_, ok := gs.Engine.PosCache.GetOneByCoordAndComponents(x, y, gs.CurrentZ, []string{constants.Visitable, constants.IsBlocking})
	return ok
}

// Effects with targets

type TargetingType string

const RandomAcquisitionType TargetingType = "random"
const ManualAcquisitionType TargetingType = "manual"
const AreaAcquisitionType TargetingType = "area"

type EffectFunction func(engine *ecs.Engine, gs *GameState, item *ecs.Entity, source *ecs.Entity, target *ecs.Entity)

type EffectInfo interface {
	EffectInfo() string
	EffectInfoColor() color.Color
}
