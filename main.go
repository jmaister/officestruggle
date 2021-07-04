package main

import (
	"fmt"

	tl "github.com/JoelOtter/termloop"
	"jordiburgos.com/officestruggle/dungeon"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

type Controller struct {
	*tl.Entity
	Engine *ecs.Engine
}

func (ctl *Controller) Tick(event tl.Event) {
	if event.Type == tl.EventKey {

		var move state.MoveComponent
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			move = state.MoveComponent{
				X: 1, Y: 0,
			}
		case tl.KeyArrowLeft:
			move = state.MoveComponent{
				X: -1, Y: 0,
			}
		case tl.KeyArrowUp:
			move = state.MoveComponent{
				X: 0, Y: -1,
			}
		case tl.KeyArrowDown:
			move = state.MoveComponent{
				X: 0, Y: 1,
			}
		}

		player := ctl.Engine.GetEntities([]string{"player"})[0]
		player.AddComponent(state.Move, move)
	}

	// This is what defines a turn step
	// systems.Render not needed, done in Draw(...) func
	systems.Movement(ctl.Engine)
}

func (ctl *Controller) Draw(screen *tl.Screen) {
	systems.Render(ctl.Engine, screen)
}

func main() {

	// ECS engine
	engine := ecs.NewEngine()

	// Dungeon
	g := grid.Grid{
		Width:  100,
		Height: 34,
		Map: grid.Map{
			Width:  79,
			Height: 29,
			X:      21,
			Y:      3,
		},
	}
	dungeonRectangle := dungeon.CreateDungeon(engine, g)
	fmt.Println("Grid", g)

	// Player
	player := engine.NewEntity()
	player.AddComponent(state.Player, state.PlayerComponent{})
	player.AddComponent(state.Apparence, state.ApparenceComponent{Color: "#fff", Char: '@'})
	player.AddComponent(state.Position, state.PositionComponent{X: dungeonRectangle.Center.X, Y: dungeonRectangle.Center.Y})

	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})

	ctl := Controller{
		Engine: engine,
	}
	level.AddEntity(&ctl)

	game.Screen().SetLevel(level)
	game.Start()

}
