package main

import tl "github.com/JoelOtter/termloop"

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

func main() {
	game := tl.NewGame()
	game.Screen().SetFps(30)
	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorGreen,
		Fg: tl.ColorBlack,
		Ch: 'v',
	})

	player := Player{tl.NewEntity(10, 10, 1, 1)}
	player.SetCell(0, 0, &tl.Cell{Fg: tl.ColorWhite, Ch: '@'})
	level.AddEntity(&player)

	//level.AddEntity(tl.NewText(10, 10, "@", tl.ColorWhite, tl.ColorBlack))
	game.Screen().SetLevel(level)
	game.Start()
}
