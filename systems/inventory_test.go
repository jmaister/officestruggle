package systems_test

import (
	"math/rand"
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

func TestPickupMoney(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	inventory := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	assert.Equal(t, 0, inventory.Coins)

	money := state.NewMoneyAmount(engine.NewEntity(), 7)
	systems.PickupEntity(gs, money)
	inventory2 := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	assert.Equal(t, 7, inventory2.Coins)

	randomMoney := rand.Intn(1000000)
	money2 := state.NewMoneyAmount(engine.NewEntity(), randomMoney)
	systems.PickupEntity(gs, money2)
	inventory3 := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	assert.Equal(t, 7+randomMoney, inventory3.Coins)

}

func TestPickupItem(t *testing.T) {
	engine := ecs.NewEngine()
	gs := game.NewGameState(engine)

	inventory := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	assert.Equal(t, 0, inventory.Coins)

	sword := state.NewSword(engine.NewEntity())
	systems.PickupEntity(gs, sword)
	inventory2 := gs.Player.GetComponent(constants.Inventory).(state.InventoryComponent)
	assert.Equal(t, 1, len(inventory2.Items))

}
