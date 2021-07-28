package animations

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/gamestate"
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
