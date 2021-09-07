package systems

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/palette"
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
				percent := 1.0 - (float64(remaining) / float64(animation.GetAnimationInfo().Duration.Milliseconds()))
				// TODO: update to use camera coordinates
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
	pos := a.Target
	camX, camY := gs.Camera.ToCameraCoordinates(pos.X, pos.Y)
	camX += gs.Grid.Camera.X
	camY += gs.Grid.Camera.Y
	x, y := toPixel(gs, camX+a.Direction.X, camY+a.Direction.Y)

	x = x + int(float64(3*gs.TileWidth)*(percent))*a.Direction.X
	y = y + int(float64(3*gs.TileHeight)*(percent))*a.Direction.Y

	fnt := assets.MplusFont(20)
	text.Draw(screen, a.Damage, fnt, x, y, palette.PColor(palette.Red, percent))
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

	newColor := palette.PColor(palette.Red, percent)
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
	Text      string
}

func (a FallingCharAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	a.StartTime = time.Now()
	if source != nil {
		srcPos := state.GetPosition(source)
		a.Source = interfaces.Point{
			X: srcPos.X,
			Y: srcPos.Y,
		}
	}
	if target != nil {
		targetPos := state.GetPosition(target)
		a.Target = interfaces.Point{
			X: targetPos.X,
			Y: targetPos.Y,
		}
	}
	return a
}
func (a FallingCharAnimation) NeedsInit() bool {
	return true
}

func (a FallingCharAnimation) GetAnimationInfo() interfaces.AnimationInfo {
	return a.AnimationInfo
}

func (a FallingCharAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	srcPos := a.Source
	tgtPos := a.Target
	line := BresenhamLine(srcPos.X, srcPos.Y, tgtPos.X, tgtPos.Y)

	current := len(line) * int((percent)*100) / 100
	if current < 0 {
		current = 0
	} else if current > len(line)-1 {
		current = len(line) - 1
	}

	camX, camY := gs.Camera.ToCameraCoordinates(line[current].X, line[current].Y)
	camX += gs.Grid.Camera.X
	camY += gs.Grid.Camera.Y
	x, y := toPixel(gs, camX, camY)

	x = x + randInt(-5, 5)
	y = y + randInt(-5, 5)

	cl := ColorBlend(a.Color, color.White, percent)
	fnt := assets.MplusFont(20)
	text.Draw(screen, a.Char, fnt, x, y, cl)
}
func (a FallingCharAnimation) End(engine *ecs.Engine, gs *gamestate.GameState, entity *ecs.Entity) {
	engine.DestroyEntity(entity)
	// Trigger damage animation
	CreateDamageAnimation(engine, a.Source, a.Target, a.Text)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// LevelUp animation

type LevelUpAnimation struct {
	interfaces.AnimationInfo
}

func (a LevelUpAnimation) Init(source *ecs.Entity, target *ecs.Entity) interfaces.Animation {
	return a
}
func (a LevelUpAnimation) NeedsInit() bool {
	return false
}

func (a LevelUpAnimation) GetAnimationInfo() interfaces.AnimationInfo {
	return a.AnimationInfo
}

func (a LevelUpAnimation) Update(percent float64, gs *gamestate.GameState, screen *ebiten.Image) {
	minRadius := 1.0
	maxRadius := 5.0
	radius := int(lerp(minRadius, maxRadius, percent))

	circle := grid.GetCircle(grid.Tile{
		X: a.Source.X,
		Y: a.Source.Y,
	}, radius)
	for _, tile := range circle {
		cl := palette.PColor(palette.Yellow, rand.Float64())

		camX, camY := gs.Camera.ToCameraCoordinates(tile.X, tile.Y)
		camX += gs.Grid.Camera.X
		camY += gs.Grid.Camera.Y
		DrawTile(screen, gs, camX, camY, cl)
	}
}
func (a LevelUpAnimation) End(engine *ecs.Engine, gs *gamestate.GameState, entity *ecs.Entity) {
}

func lerp(v0 float64, v1 float64, t float64) float64 {
	return (1-t)*v0 + t*v1
}
