package ecs_test

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/ecs"
)

func TestWhite(t *testing.T) {
	thisColor, err := ecs.ParseHexColorFast("#FFFFFF")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, uint8(255), thisColor.R)
	assert.Equal(t, uint8(255), thisColor.G)
	assert.Equal(t, uint8(255), thisColor.B)
}

func TestBlack(t *testing.T) {
	thisColor, err := ecs.ParseHexColorFast("#000000")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, uint8(0), thisColor.R)
	assert.Equal(t, uint8(0), thisColor.G)
	assert.Equal(t, uint8(0), thisColor.B)
}

func TestRandom(t *testing.T) {
	r := rand.Int63n(255)
	g := rand.Int63n(255)
	b := rand.Int63n(255)
	colorStr := "#" + toHex(r) + toHex(g) + toHex(b)
	thisColor, err := ecs.ParseHexColorFast(colorStr)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, uint8(r), thisColor.R)
	assert.Equal(t, uint8(g), thisColor.G)
	assert.Equal(t, uint8(b), thisColor.B)
}

func toHex(n int64) string {
	src := []byte{byte(n)}
	return hex.EncodeToString(src)
}
