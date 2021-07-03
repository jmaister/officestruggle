package main

import (
	tl "github.com/JoelOtter/termloop"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

type Player struct {
	*tl.Entity
}

func (player *Player) Tick(event tl.Event) {
	if event.Type == tl.EventKey { // Is it a keyboard event?
		x, y := player.Position()
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			player.SetPosition(x+1, y)
		case tl.KeyArrowLeft:
			player.SetPosition(x-1, y)
		case tl.KeyArrowUp:
			player.SetPosition(x, y-1)
		case tl.KeyArrowDown:
			player.SetPosition(x, y+1)
		}
	}
}

type Controller struct {
	*tl.Entity
	Engine ecs.Engine
}

func (ctl *Controller) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		// fmt.Println("Keyboard", event.Ch)
		player := ctl.Engine.GetEntities([]string{"player"})[0]
		pos := player.GetComponent(state.Position).(state.PositionComponent)
		x, y := pos.X, pos.Y
		switch event.Key { // If so, switch on the pressed key.
		case tl.KeyArrowRight:
			x = x + 1
		case tl.KeyArrowLeft:
			x = x - 1
		case tl.KeyArrowUp:
			y = y - 1
		case tl.KeyArrowDown:
			y = y + 1
		}

		player.RemoveComponent(state.Position)
		player.AddComponent(state.Position, state.PositionComponent{X: x, Y: y})

	}
}

func (ctl *Controller) Draw(screen *tl.Screen) {
	// fmt.Println("draw")
	renderable := []string{"position", "apparence"}
	for _, entity := range ctl.Engine.GetEntities(renderable) {
		position, _ := entity.GetComponent(state.Position).(state.PositionComponent)
		apparence, _ := entity.GetComponent(state.Apparence).(state.ApparenceComponent)

		screen.RenderCell(position.X, position.Y, &tl.Cell{Fg: systems.CssToAttr(apparence.Color), Ch: apparence.Char})
	}
}

func main() {

	engine := ecs.NewEngine()
	player := engine.NewEntity()
	player.AddComponent(state.Player, state.PlayerComponent{})
	player.AddComponent(state.Apparence, state.ApparenceComponent{
		Color: "#fff",
		Char:  '@',
	})
	player.AddComponent(state.Position, state.PositionComponent{X: 10, Y: 10})

	game := tl.NewGame()
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})

	ctl := Controller{
		Engine: *engine,
	}
	level.AddEntity(&ctl)

	// systems.Render(engine, level)

	game.Screen().SetLevel(level)
	game.Start()

}
