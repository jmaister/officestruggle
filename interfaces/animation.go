package interfaces

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

type Point struct {
	X int
	Y int
}

type Animation interface {
	GetAnimationInfo() AnimationInfo
	NeedsInit() bool
	Init(sourceEntity *ecs.Entity, targetEntity *ecs.Entity) Animation
	Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image)
	End(engine *ecs.Engine, gs *gamestate.GameState, animationEntity *ecs.Entity)
}

type AnimationInfo struct {
	StartTime time.Time
	Duration  time.Duration
	Source    Point
	Target    Point
}
