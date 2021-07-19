package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
)

func RenderWelcomesScreen(engine *ecs.Engine, gameState *gamestate.GameState, screen *ebiten.Image) {

	text.Draw(screen, "Office Struggle", fnt40, 30, 35, color.White)

	text.Draw(screen, "Story...", fnt20, 35, 200, color.White)
	text.Draw(screen, "Objective...", fnt20, 35, 300, color.White)

	text.Draw(screen, "Press [ENTER] to start", fnt20, 35, 500, color.White)

}
