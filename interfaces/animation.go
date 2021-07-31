package interfaces

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

type Animation interface {
	NeedsInit() bool
	Init(sourceEntity *ecs.Entity, targetEntity *ecs.Entity) Animation
	StartTime() time.Time
	Duration() time.Duration
	Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image)
	End(engine *ecs.Engine, gs *gamestate.GameState, animationEntity *ecs.Entity)
}
