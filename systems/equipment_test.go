package systems_test

import (
	"strings"
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

func TestEquipEntitySimple(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	sword := state.NewSword(engine.NewEntity())

	systems.EquipEntity(gs, sword)

	inventory := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	equipment := gs.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)

	assert.Equal(t, 0, len(inventory.Items))
	assert.Equal(t, 1, len(equipment.Items))
}

func TestEquipEntityReplaceEquipment(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	// Equip new item
	sword := state.NewSword(engine.NewEntity())
	systems.EquipEntity(gs, sword)
	inventory := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	equipment := gs.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	assert.Equal(t, 0, len(inventory.Items))
	assert.Equal(t, 1, len(equipment.Items))

	// Equip the same item
	systems.EquipEntity(gs, sword)
	inventory2 := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	equipment2 := gs.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	assert.Equal(t, 0, len(inventory2.Items))
	assert.Equal(t, 1, len(equipment2.Items))

	// Equip new item in the same slot
	sword3 := state.NewSword(engine.NewEntity())
	systems.EquipEntity(gs, sword3)
	inventory3 := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	equipment3 := gs.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	assert.Equal(t, 1, len(inventory3.Items))
	assert.Equal(t, 1, len(equipment3.Items))

}

func TestEquipEntityReplaceEquipmentWithFullInventory(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	// Fill the inventory
	inv := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	for i := 0; i < inv.MaxItems; i++ {
		potion := state.NewHealthPotion(engine.NewEntity())
		systems.PickupEntity(gs, potion)
	}

	// Equip new item
	sword := state.NewSword(engine.NewEntity())
	systems.EquipEntity(gs, sword)
	inventory := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	equipment := gs.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	assert.Equal(t, 10, len(inventory.Items))
	assert.Equal(t, 1, len(equipment.Items))

	// Equip the same item
	systems.EquipEntity(gs, sword)
	inventory2 := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	equipment2 := gs.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	assert.Equal(t, 10, len(inventory2.Items))
	assert.Equal(t, 1, len(equipment2.Items))

	// Equip new item in the same slot
	sword3 := state.NewSword(engine.NewEntity())
	systems.EquipEntity(gs, sword3)
	inventory3 := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	equipment3 := gs.Player.GetComponent(constants.Equipment).(state.EquipmentComponent)
	assert.Equal(t, 10, len(inventory3.Items))
	assert.Equal(t, 1, len(equipment3.Items))

	lastlog := gs.GetLogLines(1)
	assert.Equal(t, 1, len(lastlog))
	assert.Assert(t, strings.Contains(lastlog[0].Msg, "Inventory is full"))
}
