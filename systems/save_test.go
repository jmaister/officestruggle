package systems_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/systems"
)

const testSuffix = "TeSt"

func getTestSaveFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s-%s.save", gamestate.SaveGamePrefix, testSuffix)
	return path.Join(home, filename), nil

}

func TestSaveGame(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	saveFileName, err := getTestSaveFile()
	if err != nil {
		fmt.Printf("Error saving the game: %s\n", err)
		t.Fatal(err)
		return
	}

	err2 := systems.SaveGame(gs.Engine, gs, saveFileName)
	if err2 != nil {
		t.Fatal(err)
	}

}

func TestLoadGame(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	saveFileName, err := getTestSaveFile()
	if err != nil {
		fmt.Printf("Error loading the game: %s\n", err)
		t.Fatal(err)
		return
	}

	err2 := systems.LoadGame(engine, gs, saveFileName)
	if err2 != nil {
		t.Fatal(err)
	}

}
