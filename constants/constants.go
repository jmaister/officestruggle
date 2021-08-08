package constants

import (
	"image/color"

	"jordiburgos.com/officestruggle/palette"
)

const (
	Player        = "player"
	Apparence     = "apparence"
	Position      = "position"
	Move          = "move"
	IsBlocking    = "isBlocking"
	IsFloor       = "isFloor"
	Layer100      = "layer100" // Walls and floors
	Layer300      = "layer300" // Objects
	Layer400      = "layer400" // Player and enemies
	Layer500      = "layer500" // Animations
	Visitable     = "visitable"
	Description   = "description"
	AI            = "ai"
	Stats         = "stats"
	Consumable    = "consumable"
	IsPickup      = "isPickup"
	Dead          = "dead"
	Inventory     = "inventory"
	Equipable     = "equipable"
	Equipment     = "equipment"
	Animated      = "animated"
	ConsumeEffect = "consumeEffect"
	Paralize      = "paralize"
	Stairs        = "stairs"
)

type LogType string

const (
	Info   LogType = "i"
	Warn   LogType = "w"
	Bad    LogType = "b"
	Danger LogType = "d"
	Good   LogType = "g"
)

var LogColors = map[LogType]color.Color{
	Info:   color.White,
	Warn:   palette.PColor(palette.Orange, 0.6),
	Bad:    palette.PColor(palette.Red, 0.6),
	Danger: palette.PColor(palette.Red, 0.3),
	Good:   palette.PColor(palette.Green, 0.6),
}
