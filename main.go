package main

import (
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
	"jordiburgos.com/officestruggle/dungeon"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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

	// Game state
	gameState := state.NewGameState(&g, player)

	// Enemies
	visitables := engine.Entities.GetEntities([]string{state.IsFloor})
	for i := 0; i < 5; i++ {
		v := visitables[rand.Intn(len(visitables))]
		pos := state.GetPosition(v)
		goblin := state.NewGlobin(engine.NewEntity())
		state.ApplyPosition(goblin, pos.X, pos.Y)
	}

	game := tl.NewGame()
	game.Screen().SetFps(30)
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})

	ctl := systems.Controller{
		Engine:    engine,
		GameState: gameState,
		Grid:      &g,
	}
	level.AddEntity(&ctl)

	game.Screen().SetLevel(level)
	game.Start()
}
