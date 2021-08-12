package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
)

func ToTile(gs *gamestate.GameState, x int, y int) (int, int) {
	tx := x / gs.TileWidth
	ty := y / gs.TileHeight

	return tx, ty
}

func ToPixel(gs *gamestate.GameState, x int, y int) (int, int) {
	x1 := gs.TileWidth * x
	y1 := gs.TileHeight * y
	return x1, y1
}

func ToPixelRect(gs *gamestate.GameState, x int, y int) (int, int, int, int) {
	x1 := gs.TileWidth * x
	y1 := gs.TileHeight * y

	x2 := x1 + gs.TileWidth - 1
	y2 := y1 + gs.TileHeight - 1

	return x1, y1, x2, y2
}

func DrawTile(screen *ebiten.Image, gs *gamestate.GameState, x int, y int, color color.Color) {
	x1, y1 := ToPixel(gs, x, y)
	ebitenutil.DrawRect(screen, float64(x1), float64(y1), float64(gs.TileWidth), float64(gs.TileHeight), color)
}

func DrawChar(screen *ebiten.Image, gs *gamestate.GameState, x int, y int, font font.Face, str string, fgColor color.Color, bgColor color.Color) {
	x1, y1 := ToPixel(gs, x, y)
	ebitenutil.DrawRect(screen, float64(x1), float64(y1), float64(gs.TileWidth), float64(gs.TileHeight), bgColor)

	// rect := text.BoundString(font, str)
	// rect.Dx(), rect.Dy()
	h1 := gs.TileHeight

	xx := x1
	yy := y1 + h1
	text.Draw(screen, str, font, xx, yy, fgColor)

}

func DrawText(screen *ebiten.Image, gs *gamestate.GameState, x int, y int, font font.Face, str string, fgColor color.Color, bgColor color.Color) {
	// TODO: handle line breaks
	for i := 0; i < len(str); i++ {
		s := string(str[i])
		DrawChar(screen, gs, x+i, y, font, s, fgColor, bgColor)
	}
}

func DrawGridRect(screen *ebiten.Image, gs *gamestate.GameState, r grid.Rect, color color.Color) {

	x1, y1 := ToPixel(gs, r.X, r.Y)
	x2, y2 := ToPixel(gs, r.X+r.Width, r.Y+r.Height)

	y1 = y1 - 1

	// x1, y1   -1-    x2, y1
	//   -2-             -4-
	// x1, y2   -3-    x2, y2

	ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x2), float64(y1), color)
	ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x1), float64(y2), color)
	ebitenutil.DrawLine(screen, float64(x1), float64(y2), float64(x2), float64(y2), color)
	ebitenutil.DrawLine(screen, float64(x2), float64(y1), float64(x2), float64(y2), color)
}
