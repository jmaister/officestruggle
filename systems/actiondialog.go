package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/state"
)

func DrawActionDialog(engine *ecs.Engine, gs *gamestate.GameState, screen *ebiten.Image) {

	if gs.ActionScreenState.Items == nil {
		findActionableEntities(engine, gs)
	}

	if len(gs.ActionScreenState.Items) > 0 {
		title := "Available actions"
		invStrItems := []string{}
		for _, item := range gs.ActionScreenState.Items {
			invStrItems = append(invStrItems, state.GetLongDescription(item))
		}
		DrawSelectionList(screen, gs, &gs.ActionScreenState.Actions, invStrItems, gs.Grid.Inventory, title)
	} else {
		noActionsAvailable(gs)
	}
}

func findActionableEntities(engine *ecs.Engine, gs *gamestate.GameState) {
	pos := state.GetPosition(gs.Player)
	entities, ok := engine.PosCache.GetByCoord(pos.X, pos.Y, pos.Z)

	actionableEntities := ecs.EntityList{}
	if ok {
		for _, item := range entities {
			if item.HasComponent(constants.Equipable) ||
				item.HasComponent(constants.Consumable) ||
				item.HasComponent(constants.ConsumeEffect) ||
				item.HasComponent(constants.Money) ||
				item.HasComponent(constants.Stairs) {
				actionableEntities = append(actionableEntities, item)
			}
		}
	}
	if len(actionableEntities) > 0 {
		gs.ActionScreenState.Items = actionableEntities
	} else {
		gs.ActionScreenState.Items = nil
	}
}

func noActionsAvailable(gs *gamestate.GameState) {
	gs.ScreenState = gamestate.GameScreen
	gs.Log(constants.Warn, "No actions available here.")
}

func ActionDialogKeyUp(gs *gamestate.GameState) {
	updateActionSelection(gs, -1)
}

func ActionDialogKeyDown(gs *gamestate.GameState) {
	updateActionSelection(gs, 1)
}

func updateActionSelection(gs *gamestate.GameState, change int) {

	selected := gs.ActionScreenState.Actions.Selected + change
	max := len(gs.ActionScreenState.Items)
	if selected < 0 {
		selected = 0
	} else if max > 0 && selected >= max {
		selected = max - 1
	}

	gs.ActionScreenState.Actions.Selected = selected
}

func ActionDialogActivate(engine *ecs.Engine, gs *gamestate.GameState) {
	activatedEntity := gs.ActionScreenState.Items[gs.ActionScreenState.Actions.Selected]

	if activatedEntity.HasComponent(constants.Consumable) || activatedEntity.HasComponent(constants.ConsumeEffect) {
		ConsumeConsumableComponent(engine, gs, activatedEntity)
		gs.ScreenState = gamestate.GameScreen
	} else if activatedEntity.HasComponent(constants.Equipable) {
		EquipEntity(gs, activatedEntity)
	} else if activatedEntity.HasComponent(constants.Stairs) {
		UseStairs(gs, activatedEntity)
	} else if activatedEntity.HasComponent(constants.Money) {
		PickupEntity(gs, activatedEntity)
	}

	gs.ScreenState = gamestate.GameScreen
}

func ActionDialogExit(gs *gamestate.GameState) {
	gs.ScreenState = gamestate.GameScreen
}
