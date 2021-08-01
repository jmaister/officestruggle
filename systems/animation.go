package systems

import (
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

			now := time.Now()
			start := animation.GetAnimationInfo().StartTime
			end := start.Add(animation.GetAnimationInfo().Duration)

			remaining := end.Sub(now).Milliseconds()
			if remaining >= 0 {
				percent := float64(remaining) / float64(animation.GetAnimationInfo().Duration.Milliseconds())
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
	interfaces.AnimationInfo

	Direction grid.Direction
	Damage    string
}

func (a DamageAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	return a
}
func (a DamageAnimation) NeedsInit() bool {
	return false
}
func (a DamageAnimation) GetAnimationInfo() interfaces.AnimationInfo {
	return a.AnimationInfo
}
func (a DamageAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	pos := a.Target.GetComponent(constants.Position).(state.PositionComponent)
	x, y := toPixel(gs, pos.X+a.Direction.X, pos.Y+a.Direction.Y)

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
	interfaces.AnimationInfo

	StartingApparence state.ApparenceComponent
}

func (a HealthPotionAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	return a
}
func (a HealthPotionAnimation) NeedsInit() bool {
	return false
}
func (a HealthPotionAnimation) GetAnimationInfo() interfaces.AnimationInfo {
	return a.AnimationInfo
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
	interfaces.AnimationInfo

	Direction grid.Direction
	Char      string
	Color     color.Color
}

func (a FallingCharAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	a.StartTime = time.Now()
	a.Source = source
	a.Target = target
	return a
}
func (a FallingCharAnimation) NeedsInit() bool {
	return true
}

func (a FallingCharAnimation) GetAnimationInfo() interfaces.AnimationInfo {
	return a.AnimationInfo
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
