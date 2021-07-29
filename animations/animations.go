package animations

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

type Animation interface {
	StartTime() time.Time
	Duration() time.Duration
	Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image)
	End(engine *ecs.Engine, gs *gamestate.GameState, entity *ecs.Entity)
}

type AnimatedComponent struct {
	Animation Animation
}

func (a AnimatedComponent) ComponentType() string {
	return state.Animated
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

func (a HealthPotionAnimation) StartTime() time.Time {
	return a.AnimationStart
}
func (a HealthPotionAnimation) Duration() time.Duration {
	return a.AnimationDuration
}
func (a HealthPotionAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	player := gs.Player
	apparence, _ := player.GetComponent(state.Apparence).(state.ApparenceComponent)
	newColor := ""
	if (percent > 0 && percent <= 0.25) || (percent > 0.5 && percent <= 0.75) {
		newColor = "#FF0000"
	} else {
		newColor = a.StartingApparence.Color
	}
	if newColor != apparence.Color {
		apparence.Color = newColor
		player.ReplaceComponent(state.Apparence, state.ApparenceComponent{
			Color: newColor,
			Char:  apparence.Char,
		})
	}
}
func (a HealthPotionAnimation) End(engine *ecs.Engine, gs *gamestate.GameState, entity *ecs.Entity) {
	gs.Player.ReplaceComponent(state.Apparence, a.StartingApparence)
}
