package systems

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

func SaveGame(engine *ecs.Engine, gs *gamestate.GameState) error {
	gob.Register(ecs.Engine{})
	gob.Register(gamestate.GameState{})

	cleanCycles(engine, gs)

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(engine)
	restoreCycles(engine, gs)
	if err != nil {
		fmt.Println("Failed to save ENGINE.")
		return err
	}
	gs.Engine = nil
	err = encoder.Encode(gs)
	if err != nil {
		fmt.Println("Failed to save GAMESTATE.")
		return err
	}

	home, err := os.UserHomeDir()
	outputFileName := path.Join(home, "save.save")

	err2 := os.WriteFile(outputFileName, buffer.Bytes(), 0666)
	if err != nil {
		return err2
	}
	return nil
}

func LoadSystem(engine *ecs.Engine, gs *gamestate.GameState) {

}

func cleanCycles(engine *ecs.Engine, gs *gamestate.GameState) {
	// Remove cycles
	for _, entity := range engine.Entities {
		entity.Engine = nil
	}

	gs.Fov = nil
}

func restoreCycles(engine *ecs.Engine, gs *gamestate.GameState) {
	for _, entity := range engine.Entities {
		entity.Engine = engine
	}
}
