package dungeon

import (
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
)

func CreateDungeon(engine *ecs.Engine, g grid.Grid) grid.Rectangle {
	m := g.Map
	dungeon, tiles := grid.GetRectangle(m.X, m.Y, m.Width, m.Height, false)

	for _, tile := range tiles {
		tileEntity := engine.NewEntity()
		tileEntity.AddComponent(state.Apparence, state.ApparenceComponent{Color: "#555", Char: '•'})
		tileEntity.AddComponent(state.Position, state.PositionComponent{X: tile.X, Y: tile.Y})
	}

	return dungeon
}
