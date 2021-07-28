package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

func RenderTestScreen(engine *ecs.Engine, gs *gamestate.GameState, screen *ebiten.Image) {

	font := assets.MplusFont(float64(18))
	text.Draw(screen, "Office Struggle", font, 400, 400, color.White)

	bgColor := color.White
	fgColor := color.Black

	chars := []string{".", "!", "@", "#", "(", ")", "!", "@", "#", ".", "a", "A"}

	xmax := 10
	ymax := 10

	for x := 0; x < xmax; x++ {
		for y := 0; y < ymax; y++ {
			ch := chars[y%len(chars)]
			//ch := fmt.Sprintf("%d", x%10)
			DrawChar(screen, gs, x, y, font, ch, fgColor, bgColor)
		}
	}

}
