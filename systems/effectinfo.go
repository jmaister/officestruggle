package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"jordiburgos.com/officestruggle/assets"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

const fontSize = 10

var effectInfoFont font.Face = assets.MplusFont(float64(fontSize))

func EffectInfoSystem(engine *ecs.Engine, gs *gamestate.GameState, screen *ebiten.Image) {

	entities := engine.Entities.GetEntities([]string{constants.Paralize})
	entities = FilterZ(entities, gs.CurrentZ)

	for _, entity := range entities {
		effect := entity.GetComponent(constants.Paralize)
		effectInfo, ok := effect.(gamestate.EffectInfo)
		if ok {
			pos := state.GetPosition(entity)
			x, y := ToPixel(gs, pos.X, pos.Y)
			text.Draw(screen, effectInfo.EffectInfo(), effectInfoFont, x, y+gs.TileHeight-fontSize, effectInfo.EffectInfoColor())
		}
	}
}
