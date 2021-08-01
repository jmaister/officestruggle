package systems

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/state"
)

func AnimationSystem(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {
	entities := engine.Entities.GetEntities([]string{constants.Animated})

	for _, entity := range entities {
		animatedCmp, ok := entity.GetComponent(constants.Animated).(state.AnimatedComponent)
		if ok {
			animation := animatedCmp.Animation
			fmt.Println("anim", animation)

			now := time.Now()
			start := animation.StartTime()
			end := start.Add(animation.Duration())

			remaining := end.Sub(now).Milliseconds()
			if remaining >= 0 {
				percent := float64(remaining) / float64(animation.Duration().Milliseconds())
				animation.Update(percent, gameState, screen)
			} else {
				entity.RemoveComponent(constants.Animated)
				animation.End(engine, gameState, entity)
			}

		}
	}

}

// Damage Animation

type DamageAnimation struct {
	X                 int
	Y                 int
	Direction         grid.Direction
	Damage            string
	AnimationStart    time.Time
	AnimationDuration time.Duration
}

func (a DamageAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	return a
}
func (a DamageAnimation) NeedsInit() bool {
	return false
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
func (a DamageAnimation) End(engine *ecs.Engine, gs *gamestate.GameState, entity *ecs.Entity) {
	engine.DestroyEntity(entity)
}

func toPixel(gs *gamestate.GameState, x int, y int) (int, int) {
	x1 := gs.TileWidth * x
	y1 := gs.TileHeight * y
	return x1, y1
}

// Health Potion Animation

type HealthPotionAnimation struct {
	AnimationStart    time.Time
	AnimationDuration time.Duration
	StartingApparence state.ApparenceComponent
}

func (a HealthPotionAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	return a
}
func (a HealthPotionAnimation) NeedsInit() bool {
	return false
}
func (a HealthPotionAnimation) StartTime() time.Time {
	return a.AnimationStart
}
func (a HealthPotionAnimation) Duration() time.Duration {
	return a.AnimationDuration
}
func (a HealthPotionAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	player := gs.Player
	apparence, _ := player.GetComponent(constants.Apparence).(state.ApparenceComponent)
	newColor := ""

	// TODO: use palette.PColor(..., percent)
	if (percent > 0 && percent <= 0.25) || (percent > 0.5 && percent <= 0.75) {
		newColor = "#FF0000"
	} else {
		newColor = a.StartingApparence.Color
	}
	if newColor != apparence.Color {
		apparence.Color = newColor
		player.ReplaceComponent(state.ApparenceComponent{
			Color: newColor,
			Char:  apparence.Char,
		})
	}
}
func (a HealthPotionAnimation) End(engine *ecs.Engine, gs *gamestate.GameState, entity *ecs.Entity) {
	gs.Player.ReplaceComponent(a.StartingApparence)
}

// Falling char animation

type FallingCharAnimation struct {
	Direction         grid.Direction
	Char              string
	Color             color.Color
	AnimationStart    time.Time
	AnimationDuration time.Duration
	Source            *ecs.Entity
	Target            *ecs.Entity
}

func (a FallingCharAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	a.AnimationStart = time.Now()
	a.Source = source
	a.Target = target
	return a
}
func (a FallingCharAnimation) NeedsInit() bool {
	return true
}
func (a FallingCharAnimation) StartTime() time.Time {
	return a.AnimationStart
}
func (a FallingCharAnimation) Duration() time.Duration {
	return a.AnimationDuration
}
func (a FallingCharAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	pos := a.Target.GetComponent(constants.Position).(state.PositionComponent)
	x, y := toPixel(gs, pos.X+a.Direction.X, pos.Y+a.Direction.Y)

	x = x - int(float64(3*gs.TileWidth)*(1-percent))*a.Direction.X
	y = y - int(float64(3*gs.TileHeight)*(1-percent))*a.Direction.Y

	cl := ColorBlend(a.Color, color.Black, percent)
	fnt := assets.MplusFont(20)
	text.Draw(screen, a.Char, fnt, x, y, cl)
}
func (a FallingCharAnimation) End(engine *ecs.Engine, gs *gamestate.GameState, entity *ecs.Entity) {
	engine.DestroyEntity(entity)
}
func (a *FallingCharAnimation) SetStartTime(t time.Time) {
	a.AnimationStart = t
}
func (a *FallingCharAnimation) SetTarget(entity *ecs.Entity) {
	a.Target = entity
}
