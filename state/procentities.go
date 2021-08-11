package state

import (
	"fmt"
	"image/color"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/palette"
)

var setNames = map[int]string{
	1: "Black Sheep",
	2: "Red Eagle",
}

// Generates a set of items for an specific level. Level starts at 1
func GenerateEquipables(engine *ecs.Engine, level int) ecs.EntityList {

	generated := ecs.EntityList{}

	setName := setNames[level]
	for _, slot := range constants.EquipmentSlots {
		name := fmt.Sprintf("%s %s", setName, getSlotElementName(slot))
		char := getSlotChar(slot)
		clr := getSlotColor(slot)

		entity := engine.NewEntity()
		entity.AddComponent(IsPickupComponent{})
		entity.AddComponent(Layer300Component{})
		entity.AddComponent(EquipableComponent{
			EquipSlot: slot,
			MinLevel:  level,
			StatsValues: &StatsValues{
				Health:     0,
				MaxHealth:  1,
				Defense:    0,
				MaxDefense: 0,
				Power:      5,
				MaxPower:   5,
				Fov:        1,
				MaxFov:     1,
			},
		})
		entity.AddComponent(DescriptionComponent{Name: name})
		entity.AddComponent(ApparenceComponent{Color: clr, Char: char})

		generated = append(generated, entity)
	}

	return generated
}

// http://www.roguebasin.com/index.php?title=Items
func getSlotElementName(e constants.EquipSlot) string {
	switch e {
	case constants.EquipHead:
		return "Helmet"
	case constants.EquipShield:
		return "Shield"
	case constants.EquipWeapon:
		return "Sword"
	case constants.EquipBoots:
		return "Boots"
	case constants.EquipArmor:
		return "Armor"
	default:
		return "?UnKnOwN?"
	}
}

func getSlotChar(e constants.EquipSlot) rune {
	switch e {
	case constants.EquipHead:
		return '^'
	case constants.EquipShield:
		return ')'
	case constants.EquipWeapon:
		return '/'
	case constants.EquipBoots:
		return ')'
	case constants.EquipArmor:
		return ')'
	default:
		return '¿'
	}
}

func getSlotColor(e constants.EquipSlot) color.Color {
	switch e {
	case constants.EquipHead:
		return palette.PColor(palette.Green, 0.6)
	case constants.EquipShield:
		return palette.PColor(palette.Green, 0.6)
	case constants.EquipWeapon:
		return palette.PColor(palette.Blue, 0.6)
	case constants.EquipBoots:
		return palette.PColor(palette.Sepia, 0.6)
	case constants.EquipArmor:
		return palette.PColor(palette.Green, 0.6)
	default:
		return palette.PColor(palette.Pink, 0.6)
	}
}
