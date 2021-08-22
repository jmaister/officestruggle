package systems

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/palette"
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

	// Draw palette of colors
	y := 12
	for _, hue := range palette.Hues {
		for x := 0; x <= 10; x++ {
			cl := palette.PColor(hue, float64(x)/10.0)

			DrawTile(screen, gs, x, y, cl)
		}
		DrawText(screen, gs, 12, y, fnt20, strconv.Itoa(int(hue)), color.White, color.Black)
		y++
	}

}
