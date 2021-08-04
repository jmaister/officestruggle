package systems_test

import (
	"testing"

	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/systems"
)

func TestSaveGame(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	err := systems.SaveGame(engine, gs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadGame(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)
	err := systems.LoadGame(engine, gs)

	if err != nil {
		t.Fatal(err)
	}

}
