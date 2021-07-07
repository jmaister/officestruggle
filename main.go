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
	player := engine.NewEntity()
	player.AddComponent(state.Player, state.PlayerComponent{})
	player.AddComponent(state.Description, state.DescriptionComponent{Name: "You"})
	player.AddComponent(state.Apparence, state.ApparenceComponent{Color: "#0000FF", Char: '@'})
	player.AddComponent(state.Position, state.PositionComponent{X: dungeonRectangle.Center.X, Y: dungeonRectangle.Center.Y})
	player.AddComponent(state.Layer400, state.Layer400Component{})

	// Game state
	gameState := state.NewGameState(&g, player)

	// Enemies
	visitables := engine.Entities.GetEntities([]string{state.IsFloor})
	for i := 0; i < 5; i++ {
		v := visitables[rand.Intn(len(visitables))]
		pos := state.GetPosition(v)
		goblin := engine.NewEntity()
		goblin.AddComponent(state.Description, state.DescriptionComponent{Name: "Goblin"})
		goblin.AddComponent(state.IsBlocking, state.IsBlockingComponent{})
		goblin.AddComponent(state.Apparence, state.ApparenceComponent{Color: "#00FC00", Char: 'g'})
		goblin.AddComponent(state.Position, state.PositionComponent{X: pos.X, Y: pos.Y})
		goblin.AddComponent(state.Layer400, state.Layer400Component{})
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
