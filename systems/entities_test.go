package systems_test

import (
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/state"
	"jordiburgos.com/officestruggle/systems"
)

func TestNewLightningScroll(t *testing.T) {
	engine := ecs.NewEngine()

	scroll := systems.NewLightningScroll(engine.NewEntity())

	desc := state.GetLongDescription(scroll)
	assert.Equal(t, "Lightning Scroll", desc)

	apparence := scroll.GetComponent(constants.Apparence).(state.ApparenceComponent)
	s := string(apparence.Char)
	assert.Equal(t, 1, len(s))
}
