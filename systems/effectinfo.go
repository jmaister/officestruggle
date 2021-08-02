package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func EffectInfoSystem(engine *ecs.Engine, gs *gamestate.GameState, screen *ebiten.Image) {

	fontSize := 10
	font := assets.MplusFont(float64(fontSize))

	entities := engine.Entities.GetEntities([]string{constants.Paralize})
	for _, entity := range entities {
		effect := entity.GetComponent(constants.Paralize)
		effectInfo, ok := effect.(gamestate.EffectInfo)
		if ok {
			pos := state.GetPosition(entity)
			x, y := ToPixel(gs, pos.X, pos.Y)
			text.Draw(screen, effectInfo.EffectInfo(), font, x, y+gs.TileHeight-fontSize, effectInfo.EffectInfoColor())
		}
	}
}
