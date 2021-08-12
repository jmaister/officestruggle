package interfaces

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

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
	// TODO: change source and target to positions to avoid nullpointers if the entity is consumed/destroyed and the animation still active
	Source *ecs.Entity
	Target *ecs.Entity
}
