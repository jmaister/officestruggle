package systems

import (
	"fmt"
	"time"

	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/gamestate"
	"jordiburgos.com/officestruggle/interfaces"
	"jordiburgos.com/officestruggle/state"
)

func GiveXP(gs *gamestate.GameState, target *ecs.Entity, fromEntity *ecs.Entity) {

	levelingComponent, ok := target.GetComponent(constants.Leveling).(state.LevelingComponent)

	if ok {
		xpGiver, ok := fromEntity.GetComponent(constants.XPGiver).(state.XPGiverComponent)
		if ok {
			xpIncrease := xpGiver.Calculate()
			newLevelingCmp, hasChanged, hasIncreasedLevel := calculateNewLevelingComponent(levelingComponent, xpIncrease)
			if hasChanged {
				// TODO: trigger xp animation
				gs.Log(constants.Good, fmt.Sprintf("You gained %d experience points.", xpIncrease))
			}
			if hasIncreasedLevel {
				// TODO: Update players stats current and max
				gs.Log(constants.Good, fmt.Sprintf("You advance to level %d.", newLevelingCmp.CurrentLevel))
				// Trigger level increase animation
				pos := state.GetPosition(target)
				target.AddComponent(state.AnimatedComponent{
					Animation: LevelUpAnimation{
						AnimationInfo: interfaces.AnimationInfo{
							StartTime: time.Now(),
							Duration:  750 * time.Millisecond,
							Source: interfaces.Point{
								X: pos.X,
								Y: pos.Y,
							},
						},
					},
				})
			}
			target.ReplaceComponent(newLevelingCmp)
		} else {
			gs.Log(constants.Bad, fmt.Sprintf("%s does not give XP.", state.GetDescription(fromEntity)))
		}
	} else {
		gs.Log(constants.Bad, fmt.Sprintf("%s can't receive XP.", state.GetDescription(target)))
	}
}

// returns the new LevelingComponent, true/false if XP has changed, true/false if it has incresed level
func calculateNewLevelingComponent(levelingComponent state.LevelingComponent, xp int) (state.LevelingComponent, bool, bool) {
	if xp == 0 {
		return levelingComponent, false, false
	}

	levelingComponent.CurrentXP += xp

	xpToNextLevel := levelingComponent.GetNextLevelXP()
	if levelingComponent.CurrentXP >= xpToNextLevel {
		levelingComponent.CurrentLevel += 1
		levelingComponent.CurrentXP -= xpToNextLevel
		return levelingComponent, true, true
	}
	return levelingComponent, true, false

}
