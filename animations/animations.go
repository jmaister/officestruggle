package animations

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type Animation interface {
	StartTime() time.Time
	Duration() time.Duration
	Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image)
}

type AnimatedComponent struct {
	Animation Animation
}

func (a AnimatedComponent) ComponentType() string {
	return state.Animated
}

type DamageAnimation struct {
	X                 int
	Y                 int
	Direction         grid.Direction
	Damage            string
	AnimationStart    time.Time
	AnimationDuration time.Duration
}

func (a DamageAnimation) StartTime() time.Time {
	return a.AnimationStart
}
func (a DamageAnimation) Duration() time.Duration {
	return a.AnimationDuration
}
func (a DamageAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	x, y := toPixel(gs, a.X+a.Direction.X, a.Y+a.Direction.Y)

	x = x + int(float64(3*gs.TileWidth)*(1-percent))*a.Direction.X
	y = y + int(float64(3*gs.TileHeight)*(1-percent))*a.Direction.Y

	fnt := assets.MplusFont(20)
	text.Draw(screen, a.Damage, fnt, x, y, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})
}

func toPixel(gs *gamestate.GameState, x int, y int) (int, int) {
	x1 := gs.TileWidth * x
	y1 := gs.TileHeight * y
	return x1, y1
}
