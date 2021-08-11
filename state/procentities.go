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
	MinLevel      int
	MaxLevel      int
	ImprovedStats StatsValues
}

var equipmentSets = []EquipmentSet{
	{
		Name:      "Black Sheep",
		ItemLevel: 1,
		MinLevel:  1,
		MaxLevel:  5,
		ImprovedStats: StatsValues{
			Health:     0,
			MaxHealth:  0,
			Defense:    2,
			MaxDefense: 3,
			Power:      0,
			MaxPower:   0,
			Fov:        0,
			MaxFov:     0,
		},
	},
	{
		Name:      "Red Eagle",
		ItemLevel: 2,
		MinLevel:  2,
		MaxLevel:  5,
		ImprovedStats: StatsValues{
			Health:     0,
			MaxHealth:  0,
			Defense:    0,
			MaxDefense: 0,
			Power:      1,
			MaxPower:   2,
			Fov:        1,
			MaxFov:     2,
		},
	},
}

// Generates a set of items for an specific level. Level starts at 1
func GenerateEquipables(engine *ecs.Engine, level int) ecs.EntityList {

	generated := ecs.EntityList{}

	for _, equipmentSet := range equipmentSets {
		if level >= equipmentSet.MinLevel && level <= equipmentSet.MaxLevel {
			for _, slot := range constants.EquipmentSlots {
				name := fmt.Sprintf("%s %s", equipmentSet.Name, getSlotElementName(slot))
				char := getSlotChar(slot)
				clr := getSlotColor(slot)

				entity := engine.NewEntity()
				entity.AddComponent(IsPickupComponent{})
				entity.AddComponent(Layer300Component{})
				entity.AddComponent(EquipableComponent{
					EquipSlot: slot,
					MinLevel:  level,
					StatsValues: &StatsValues{
						Health:     (1 + equipmentSet.ImprovedStats.Health) * equipmentSet.ItemLevel,
						MaxHealth:  (1 + equipmentSet.ImprovedStats.MaxHealth) * equipmentSet.ItemLevel,
						Defense:    (1 + equipmentSet.ImprovedStats.Defense) * equipmentSet.ItemLevel,
						MaxDefense: (1 + equipmentSet.ImprovedStats.MaxDefense) * equipmentSet.ItemLevel,
						Power:      (1 + equipmentSet.ImprovedStats.Power) * equipmentSet.ItemLevel,
						MaxPower:   (1 + equipmentSet.ImprovedStats.MaxPower) * equipmentSet.ItemLevel,
						Fov:        (1 + equipmentSet.ImprovedStats.Fov) * equipmentSet.ItemLevel,
						MaxFov:     (1 + equipmentSet.ImprovedStats.MaxFov) * equipmentSet.ItemLevel,
					},
				})
				entity.AddComponent(DescriptionComponent{Name: name})
				entity.AddComponent(ApparenceComponent{Color: clr, Char: char})

				generated = append(generated, entity)
			}

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
