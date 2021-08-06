package systems

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/sqweek/dialog"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type Save struct {
	GameVersion  string
	Entities     ecs.EntityList
	Grid         grid.Grid
	ScreenWidth  int
	ScreenHeight int
	TileWidth    int
	TileHeight   int
	LogLines     []gamestate.LogLine
}

// YYYYMMDD_hhmmss
const dateTimeLayout = "20060102_150405"

func registerGob() {

	gob.Register(Save{})
	gob.Register(ecs.Entity{})

	gob.Register(state.AIComponent{})
	gob.Register(state.AnimatedComponent{})
	gob.Register(state.ApparenceComponent{})
	gob.Register(state.ConsumableComponent{})
	gob.Register(state.ConsumeEffectComponent{})
	gob.Register(state.DeadComponent{})
	gob.Register(state.DescriptionComponent{})
	gob.Register(state.EquipableComponent{})
	gob.Register(state.EquipmentComponent{})
	gob.Register(state.InventoryComponent{})
	gob.Register(state.IsBlockingComponent{})
	gob.Register(state.IsFloorComponent{})
	gob.Register(state.IsPickupComponent{})
	gob.Register(state.Layer100Component{})
	gob.Register(state.Layer300Component{})
	gob.Register(state.Layer400Component{})
	gob.Register(state.Layer500Component{})
	gob.Register(state.PlayerComponent{})
	gob.Register(state.StatsComponent{})
	gob.Register(state.PositionComponent{})
	gob.Register(state.VisitableComponent{})

	gob.Register(FallingCharAnimation{})
	gob.Register(HealthPotionAnimation{})
	gob.Register(DamageAnimation{})

	gob.Register(colorful.Color{})
}

func SaveGame(engine *ecs.Engine, gs *gamestate.GameState) error {
	registerGob()

	cleanCycles(engine, gs)

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)

	save := Save{
		GameVersion:  gamestate.GameVersion,
		Entities:     engine.Entities,
		Grid:         *gs.Grid,
		ScreenWidth:  gs.ScreenWidth,
		ScreenHeight: gs.ScreenHeight,
		TileWidth:    gs.TileWidth,
		TileHeight:   gs.TileHeight,
		LogLines:     gs.LogLines,
	}

	err := encoder.Encode(save)
	if err != nil {
		fmt.Println("Failed to save.")
		return err
	}

	home, err := os.UserHomeDir()
	filename := fmt.Sprintf("%s-%s.save", gamestate.SaveGamePrefix, time.Now().Format(dateTimeLayout))
	saveFileName := path.Join(home, filename)

	err2 := os.WriteFile(saveFileName, buffer.Bytes(), 0666)
	if err != nil {
		return err2
	}

	restoreCycles(engine, gs)
	return nil
}

func LoadGame(engine *ecs.Engine, gs *gamestate.GameState) error {
	registerGob()

	saveFileName, err := dialog.File().Filter(".save files", "save").Title("Load saved game").Load()
	if err != nil {
		return err
	}

	contents, err2 := os.ReadFile(saveFileName)
	if err2 != nil {
		return err2
	}

	buffer := bytes.NewBuffer(contents)
	decoder := gob.NewDecoder(buffer)

	save := Save{}
	decoder.Decode(&save)

	if gamestate.GameVersion != save.GameVersion {
		fmt.Printf("Warning: Game version [%s] differs from saved file version [%s]. You may have unexpected results.\n", gamestate.GameVersion, save.GameVersion)
	}

	engine.SetEntityList(save.Entities)

	gs.Engine = engine
	gs.Grid = &save.Grid
	gs.ScreenWidth = save.ScreenWidth
	gs.ScreenHeight = save.ScreenHeight
	gs.TileWidth = save.TileWidth
	gs.TileHeight = save.TileHeight
	gs.Player = engine.Entities.GetEntity([]string{constants.Player})
	gs.LogLines = save.LogLines

	return nil
}

func cleanCycles(engine *ecs.Engine, gs *gamestate.GameState) {
	// Remove cycles
	for _, entity := range engine.Entities {
		entity.Engine = nil
	}
}

func restoreCycles(engine *ecs.Engine, gs *gamestate.GameState) {
	for _, entity := range engine.Entities {
		entity.Engine = engine
	}
}