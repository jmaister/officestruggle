package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
)

func DrawSelectionList(screen *ebiten.Image, listState *gamestate.ListState, items []string, position grid.Rect, title string) {
	fontSize := 12
	font := assets.MplusFont(float64(fontSize))

	y := position.Y
	text.Draw(screen, title, font, (position.X)*fontSize, y*fontSize, color.White)

	if len(items) > 0 {
		for i, itemStr := range items {
			selected := (i == listState.Selected)

			str := fmt.Sprintf("%2d - %s", i+1, itemStr)
			px := (position.X) * fontSize
			py := (y + 1) * fontSize
			if selected && listState.IsFocused {
				DrawTextRect(screen, str, px, py, font, color.White)
				text.Draw(screen, str, font, px, py, color.Black)
			} else {
				text.Draw(screen, str, font, px, py, color.White)
			}
			y++
		}
	} else {
		str := "- No items -"
		px := (position.X) * fontSize
		py := (y + 1) * fontSize
		if listState.IsFocused {
			DrawTextRect(screen, str, px, py, font, color.White)
			text.Draw(screen, str, font, px, py, color.Black)
		} else {
			text.Draw(screen, str, font, px, py, color.White)
		}
	}
}
