package game

import (
	"log"
	"math/rand"
	"os"

	"github.com/norendren/go-fov/fov"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/dungeon"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/grid"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

func NewGameState(engine *ecs.Engine) *gamestate.GameState {

	// Dungeon
	g := grid.Grid{
		Width:  80,
		Height: 34,
		Levels: 2,
		Map: grid.Rect{
			X:      16,
			Y:      8,
			Width:  74,
			Height: 30,
		},
		MessageLog: grid.Rect{
			X:      0,
			Y:      0,
			Width:  79,
			Height: 5,
		},
		PlayerHud: grid.Rect{
			X:      0,
			Y:      6,
			Width:  79,
			Height: 1,
		},
		InfoBar: grid.Rect{
			X:      16,
			Y:      39,
			Width:  79,
			Height: 10,
		},
		GameInventory: grid.Rect{
			X:      0,
			Y:      8,
			Width:  15,
			Height: 40,
		},
		Inventory: grid.Rect{
			X:      2,
			Y:      5,
			Width:  79,
			Height: 29,
		},
		Equipment: grid.Rect{
			X:      50,
			Y:      5,
			Width:  20,
			Height: 29,
		},
	}
	dungeonTiles, startingTile := dungeon.CreateDungeon(g.Map, dungeon.DungeonOptions{
		MinRoomSize:  6,
		MaxRoomSize:  12,
		MaxRoomCount: 40,
	}, g.Levels)

	for _, tile := range dungeonTiles {
		tileEntity := engine.NewEntity()
		if tile.Sprite == grid.Wall {
			state.NewWall(tileEntity, tile.X, tile.Y, tile.Z)
		} else if tile.Sprite == grid.Floor {
			state.NewFloor(tileEntity, tile.X, tile.Y, tile.Z)
		}
	}

	// Player
	player := state.NewPlayer(engine.NewEntity())
	state.ApplyPosition(player, startingTile.X, startingTile.Y, startingTile.Z)

	for level := 0; level < g.Levels; level++ {
		visitables := engine.Entities.GetEntities([]string{constants.IsFloor})
		visitables = systems.FilterZ(visitables, level)
		// Enemies
		for i := 0; i < 10; i++ {
			v := visitables[rand.Intn(len(visitables))]
			pos := state.GetPosition(v)
			goblin := state.NewGlobin(engine.NewEntity())
			state.ApplyPosition(goblin, pos.X, pos.Y, pos.Z)
		}
		// Health potions
		for i := 0; i < 10; i++ {
			v := visitables[rand.Intn(len(visitables))]
			pos := state.GetPosition(v)
			potion := state.NewHealthPotion(engine.NewEntity())
			state.ApplyPosition(potion, pos.X, pos.Y, pos.Z)
		}
		// Swords
		for i := 0; i < 10; i++ {
			v := visitables[rand.Intn(len(visitables))]
			pos := state.GetPosition(v)
			potion := state.NewSword(engine.NewEntity())
			state.ApplyPosition(potion, pos.X, pos.Y, pos.Z)
		}
		// Lightning Scroll
		for i := 0; i < 10; i++ {
			v := visitables[rand.Intn(len(visitables))]
			pos := state.GetPosition(v)
			scroll := systems.NewLightningScroll(engine.NewEntity())
			state.ApplyPosition(scroll, pos.X, pos.Y, pos.Z)
		}
		// Paralize Scroll
		for i := 0; i < 10; i++ {
			v := visitables[rand.Intn(len(visitables))]
			pos := state.GetPosition(v)
			scroll := systems.NewParalizeScroll(engine.NewEntity())
			state.ApplyPosition(scroll, pos.X, pos.Y, pos.Z)
		}
	}

	return &gamestate.GameState{
		Engine:      engine,
		Fov:         fov.New(),
		Grid:        &g,
		Player:      player,
		ScreenState: gamestate.WelcomeScreen,
		InventoryScreenState: gamestate.InventoryScreenState{
			InventoryState: gamestate.ListState{
				Selected:  0,
				IsFocused: true,
			},
			EquipmentState: gamestate.ListState{
				Selected:  0,
				IsFocused: false,
			},
		},
		IsPlayerTurn: true,
		L:            log.New(os.Stderr, "", 0),
		ScreenWidth:  1920,
		ScreenHeight: 1080,
		TileWidth:    20,
		TileHeight:   20,
	}
}
