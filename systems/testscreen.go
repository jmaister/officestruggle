package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

func RenderTestScreen(engine *ecs.Engine, gs *gamestate.GameState, screen *ebiten.Image) {

	text.Draw(screen, "Office Struggle", fnt20, 10, 10, color.White)

	bgColor := color.White
	fgColor := color.Black

	chars := []string{".", "!", "@", "#", "(", ")", "!", "@", "#", ".", "a", "A"}

	xmax := gs.ScreenWidth / gs.TileWidth
	ymax := gs.ScreenHeight / gs.TileHeight

	for x := 0; x < xmax; x++ {
		for y := 0; y < ymax; y++ {
			ch := chars[y%len(chars)]
			//ch := fmt.Sprintf("%d", x%10)
			DrawChar(screen, gs, x, y, fnt20, ch, fgColor, bgColor)
		}
	}

}
