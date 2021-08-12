package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
)

func DrawSelectionList(screen *ebiten.Image, gs *gamestate.GameState, listState *gamestate.ListState, items []string, position grid.Rect, title string) {
	fontSize := 18
	font := assets.MplusFont(float64(fontSize))

	y := position.Y
	DrawText(screen, gs, position.X, position.Y, font, title, color.White, color.Black)
	y++

	if len(items) > 0 {
		for i, itemStr := range items {
			selected := (i == listState.Selected)

			str := fmt.Sprintf("%2d - %s", i+1, itemStr)
			px := position.X
			py := y
			if selected && listState.IsFocused {
				DrawText(screen, gs, px, py, font, str, color.Black, color.White)
			} else {
				DrawText(screen, gs, px, py, font, str, color.White, color.Black)
			}
			y++
		}
	} else {
		str := "- No items -"
		px := position.X
		py := y
		if listState.IsFocused {
			DrawText(screen, gs, px, py, font, str, color.Black, color.White)
		} else {
			DrawText(screen, gs, px, py, font, str, color.White, color.Black)

		}
	}
}
