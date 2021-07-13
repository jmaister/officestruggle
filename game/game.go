package game

import (
	"log"
	"math/rand"
	"os"

	"jordiburgos.com/officestruggle/dungeon"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type GameState struct {
	Engine       *ecs.Engine
	Fov          *state.FieldOfVision
	Grid         *grid.Grid
	Player       *ecs.Entity
	IsPlayerTurn bool
	L            *log.Logger
}

func NewGameState() *GameState {
	// ECS engine
	engine := ecs.NewEngine()

	// Dungeon
	g := grid.Grid{
		Width:  100,
		Height: 34,
		Map: grid.Map{
			X:      21,
			Y:      3,
			Width:  79,
			Height: 29,
		},
	}
	dungeonRectangle := dungeon.CreateDungeon(engine, g.Map, dungeon.DungeonOptions{
		MinRoomSize:  6,
		MaxRoomSize:  12,
		MaxRoomCount: 7,
	})

	// Player
	player := state.NewPlayer(engine.NewEntity())
	state.ApplyPosition(player, dungeonRectangle.Center.X, dungeonRectangle.Center.Y)

	// Enemies
	visitables := engine.Entities.GetEntities([]string{state.IsFloor})
	for i := 0; i < 5; i++ {
		v := visitables[rand.Intn(len(visitables))]
		pos := state.GetPosition(v)
		goblin := state.NewGlobin(engine.NewEntity())
		state.ApplyPosition(goblin, pos.X, pos.Y)
	}

	fov := state.FieldOfVision{}
	fov.SetTorchRadius(6)

	return &GameState{
		Engine:       engine,
		Fov:          &fov,
		Grid:         &g,
		Player:       player,
		IsPlayerTurn: true,
		L:            log.New(os.Stderr, "", 0),
	}
}
