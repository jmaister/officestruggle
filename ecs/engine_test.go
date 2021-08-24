package ecs_test

import (
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/constants"
	"jordiburgos.com/officestruggle/ecs"
	"jordiburgos.com/officestruggle/game"
)

var engine = ecs.NewEngine()
var gs = game.NewGameState(engine)

func TestGetEntities(t *testing.T) {
	ents1 := GetEntitiesFn()
	ents2 := GetEntitiesWithFilterFn()

	assert.Equal(t, len(ents1), len(ents2))
}

func GetEntitiesFn() ecs.EntityList {
	return engine.Entities.GetEntities([]string{constants.Position, constants.Apparence})
}

func BenchmarkGetEntities(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetEntitiesFn()
	}
}

func GetEntitiesWithFilterFn() ecs.EntityList {
	return engine.Entities.GetEntitiesWithFilter(filter)
}

func filter(entity *ecs.Entity) bool {
	return entity.HasComponent(constants.Position) && entity.HasComponent(constants.Apparence)

}

func BenchmarkGetEntitiesWithFilter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetEntitiesWithFilterFn()
	}
}
