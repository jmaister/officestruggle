package state

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/palette"
)

type EquipmentSet struct {
	Name          string
	ItemLevel     int
	ImprovedStats StatsValues
}

var slotDistribution = createSlotDistribution()
var tierDistribution = createTierDistribution()
var itemTypeDistribution = createItemTypeDistribution()

func createSlotDistribution() DistributedRandom {
	slotDistribution := NewDistributedRandom()
	prob := 1.0 / float64(len(constants.EquipmentSlots))
	for i := range constants.EquipmentSlots {
		slotDistribution.SetNumber(i, prob)
	}
	return slotDistribution
}

func createTierDistribution() DistributedRandom {
	tierDistribution := NewDistributedRandom()
	tierDistribution.SetNumber(1, 0.9)  // Common
	tierDistribution.SetNumber(2, 0.5)  // Uncommon
	tierDistribution.SetNumber(3, 0.1)  // Rare
	tierDistribution.SetNumber(4, 0.01) // Mythic
	return tierDistribution
}

func createItemTypeDistribution() DistributedRandom {
	typeDistribution := NewDistributedRandom()
	typeDistribution.SetNumber(1, 0.63) // Item
	typeDistribution.SetNumber(2, 0.1)  // Health
	typeDistribution.SetNumber(3, 0.09) // Lightning Scroll
	typeDistribution.SetNumber(4, 0.09) // Paralize Scroll
	typeDistribution.SetNumber(5, 0.09) // Money
	return typeDistribution
}

func CreateRandomItem(engine *ecs.Engine, level int) *ecs.Entity {

	switch itemTypeDistribution.GetDistributedRandom() {
	case 1:
		slot := constants.EquipmentSlots[slotDistribution.GetDistributedRandom()]
		equipmentSet := createEquipmentSet(level)
		tier := tierDistribution.GetDistributedRandom()
		return createEquipable(engine, equipmentSet, level, slot, tier)
	case 2:
		return NewHealthPotion(engine.NewEntity())
	case 3:
		// TODO: return NewLightningScroll(engine.NewEntity())
		return NewHealthPotion(engine.NewEntity())
	case 4:
		// TODO: return NewParalizeScroll(engine.NewEntity())
		return NewHealthPotion(engine.NewEntity())
	case 5:
		return NewMoneyAmount(engine.NewEntity(), 1000+rand.Intn(1000))
	default:
		return nil
	}
}

// Create EquipmentSet for one level
func createEquipmentSet(level int) EquipmentSet {
	equipmentLevel := int(math.Ceil(float64(level) / 3.0))
	return EquipmentSet{
		Name:      fmt.Sprintf("Level %d %s", equipmentLevel, itemLevelSet(equipmentLevel)),
		ItemLevel: equipmentLevel,
		ImprovedStats: StatsValues{
			Health:    int(float64(0.25) * float64(equipmentLevel)),
			MaxHealth: int(float64(0.5) * float64(equipmentLevel)),
			Defense:   int(float64(1.5) * float64(equipmentLevel)),
			Power:     2 * equipmentLevel,
			Fov:       int(float64(equipmentLevel)),
		},
	}

}

// Generates a set of items for an specific level. Level starts at 1
func GenerateEquipables(engine *ecs.Engine, level int) ecs.EntityList {

	generated := ecs.EntityList{}

	equipmentSet := createEquipmentSet(level)
	for _, slot := range constants.EquipmentSlots {
		entity := createEquipable(engine, equipmentSet, level, slot, tierDistribution.GetDistributedRandom())
		generated = append(generated, entity)
	}

	return generated
}

func createEquipable(engine *ecs.Engine, equipmentSet EquipmentSet, level int, slot constants.EquipSlot, tier int) *ecs.Entity {
	name := fmt.Sprintf("%s %s T%d", equipmentSet.Name, getSlotElementName(slot), tier)
	char := getSlotChar(slot)
	clr := getSlotColor(slot)

	entity := engine.NewEntity()
	entity.AddComponent(IsPickupComponent{})
	entity.AddComponent(Layer300Component{})

	statValues := StatsValues{
		Health:    (1 + equipmentSet.ImprovedStats.Health) * equipmentSet.ItemLevel,
		MaxHealth: (1 + equipmentSet.ImprovedStats.MaxHealth) * equipmentSet.ItemLevel,
		Defense:   (1 + equipmentSet.ImprovedStats.Defense) * equipmentSet.ItemLevel,
		Power:     (1 + equipmentSet.ImprovedStats.Power) * equipmentSet.ItemLevel,
		Fov:       (1 + equipmentSet.ImprovedStats.Fov) * equipmentSet.ItemLevel,
	}
	// Apply SLOT modifier, ie. weapon more power, armor more defense
	statValues = statValues.ApplyMultiplier(getSlotMultiplier(slot))

	statValues = statValues.ApplyMultiplier(getTierMultiplier(tier))

	entity.AddComponent(EquipableComponent{
		EquipSlot:   slot,
		Level:       equipmentSet.ItemLevel,
		Tier:        tier,
		SetName:     equipmentSet.Name,
		StatsValues: &statValues,
	})
	entity.AddComponent(DescriptionComponent{Name: name})
	entity.AddComponent(ApparenceComponent{Color: clr, Char: char})

	return entity
}

// http://www.roguebasin.com/index.php?title=Items
func getSlotElementName(e constants.EquipSlot) string {
	switch e {
	case constants.EquipHead:
		return "Helmet"
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

// Multipliers by slot. 100 means keep the same: (value * 1.00)
func getSlotMultiplier(e constants.EquipSlot) StatsValues {
	switch e {
	case constants.EquipHead:
		return StatsValues{
			Health:    100,
			MaxHealth: 100,
			Defense:   150,
			Power:     100,
			Fov:       200,
		}
	case constants.EquipWeapon:
		return StatsValues{
			Health:    100,
			MaxHealth: 100,
			Defense:   100,
			Power:     300,
			Fov:       100,
		}
	case constants.EquipBoots:
		return StatsValues{
			Health:    100,
			MaxHealth: 100,
			Defense:   200,
			Power:     100,
			Fov:       100,
		}
	case constants.EquipArmor:
		return StatsValues{
			Health:    200,
			MaxHealth: 200,
			Defense:   300,
			Power:     100,
			Fov:       100,
		}
	default:
		return StatsValues{
			Health:    100,
			MaxHealth: 100,
			Defense:   100,
			Power:     100,
			Fov:       100,
		}
	}
}

func getTierMultiplier(tier int) StatsValues {
	return StatsValues{
		Health:    100 + (100 * (tier - 1)),
		MaxHealth: 100 + (100 * (tier - 1)),
		Defense:   100 + (100 * (tier - 1)),
		Power:     100 + (100 * (tier - 1)),
		Fov:       100 + (25 * (tier - 1)),
	}
}

func itemLevelSet(equipmentLevel int) string {
	switch equipmentLevel {
	case 1:
		return "Wood"
	case 2:
		return "Bone"
	case 3:
		return "Iron"
	case 4:
		return "Steel"
	case 5:
		return "Diamond"
	default:
		return "¿UnKnOwN¿"
	}

}
