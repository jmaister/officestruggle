package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"jordiburgos.com/officestruggle/gamestate"
)

func ToTile(gs *gamestate.GameState, x int, y int) (int, int) {
	tx := x / gs.TileWidth
	ty := y / gs.TileHeight

	return tx, ty
}

func toPixelRect(gs *gamestate.GameState, x int, y int) (int, int, int, int) {
	x1 := gs.TileWidth * x
	y1 := gs.TileHeight * y

	x2 := x1 + gs.TileWidth - 1
	y2 := y1 + gs.TileHeight - 1

	return x1, y1, x2, y2
}

func DrawTile(screen *ebiten.Image, gs *gamestate.GameState, x int, y int, color color.Color) {
	x1, y1, x2, y2 := toPixelRect(gs, x, y)
	ebitenutil.DrawRect(screen, float64(x1), float64(y1), float64(x2), float64(y2), color)
}

func DrawChar(screen *ebiten.Image, gs *gamestate.GameState, x int, y int, font font.Face, str string, fgColor color.Color, bgColor color.Color) {
	x1, x2, y1, y2 := toPixelRect(gs, x, y)
	ebitenutil.DrawRect(screen, float64(x1), float64(y1), float64(x2), float64(y2), bgColor)

	// rect := text.BoundString(font, str)
	text.Draw(screen, str, font, x1, y1, fgColor)
}
