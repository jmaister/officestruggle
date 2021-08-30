package state

import (
	"fmt"
	"image/color"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/palette"
)

type EquipmentSet struct {
	Name          string
	ItemLevel     int
	ImprovedStats StatsValues
}

var equipmentSets = []EquipmentSet{
	{
		Name:      "Black Sheep",
		ItemLevel: 1,
		ImprovedStats: StatsValues{
			Health:    0,
			MaxHealth: 0,
			Defense:   2,
			Power:     0,
			Fov:       0,
		},
	},
	{
		Name:      "Red Eagle",
		ItemLevel: 2,
		ImprovedStats: StatsValues{
			Health:    0,
			MaxHealth: 0,
			Defense:   0,
			Power:     1,
			Fov:       1,
		},
	},
}

// Generates a set of items for an specific level. Level starts at 1
func GenerateEquipables(engine *ecs.Engine, level int) ecs.EntityList {

	generated := ecs.EntityList{}

	for _, equipmentSet := range equipmentSets {
		for _, slot := range constants.EquipmentSlots {
			name := fmt.Sprintf("%s %s", equipmentSet.Name, getSlotElementName(slot))
			char := getSlotChar(slot)
			clr := getSlotColor(slot)

			entity := engine.NewEntity()
			entity.AddComponent(IsPickupComponent{})
			entity.AddComponent(Layer300Component{})
			entity.AddComponent(EquipableComponent{
				EquipSlot: slot,
				Level:     equipmentSet.ItemLevel,
				SetName:   equipmentSet.Name,
				StatsValues: &StatsValues{
					// TODO: add SLOT modifier, ie. weapon more power, shield more defense
					Health:    (1 + equipmentSet.ImprovedStats.Health) * equipmentSet.ItemLevel,
					MaxHealth: (1 + equipmentSet.ImprovedStats.MaxHealth) * equipmentSet.ItemLevel,
					Defense:   (1 + equipmentSet.ImprovedStats.Defense) * equipmentSet.ItemLevel,
					Power:     (1 + equipmentSet.ImprovedStats.Power) * equipmentSet.ItemLevel,
					Fov:       (1 + equipmentSet.ImprovedStats.Fov) * equipmentSet.ItemLevel,
				},
			})
			entity.AddComponent(DescriptionComponent{Name: name})
			entity.AddComponent(ApparenceComponent{Color: clr, Char: char})

			generated = append(generated, entity)
		}
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
		return ']'
	case constants.EquipShield:
		return ']'
	case constants.EquipWeapon:
		return '/'
	case constants.EquipBoots:
		return ']'
	case constants.EquipArmor:
		return ']'
	default:
		return '¿'
	}
}

func getSlotColor(e constants.EquipSlot) color.Color {
	switch e {
	case constants.EquipHead:
		return palette.PColor(palette.Sepia, 0.4)
	case constants.EquipShield:
		return palette.PColor(palette.Sepia, 0.4)
	case constants.EquipWeapon:
		return palette.PColor(palette.Cyan, 0.6)
	case constants.EquipBoots:
		return palette.PColor(palette.Sepia, 0.4)
	case constants.EquipArmor:
		return palette.PColor(palette.Sepia, 0.4)
	default:
		return palette.PColor(palette.Pink, 0.4)
	}
}
