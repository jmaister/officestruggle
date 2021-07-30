package systems_test

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"gotest.tools/assert"
	"jordiburgos.com/officestruggle/systems"
)

func TestWhite(t *testing.T) {
	thisColor := systems.ParseHexColorFast("#FFFFFF")
	tr, tg, tb, ta := thisColor.RGBA()

	assert.Equal(t, uint32(255), tr>>8)
	assert.Equal(t, uint32(255), tg>>8)
	assert.Equal(t, uint32(255), tb>>8)
	assert.Equal(t, uint32(255), ta>>8)
}

func TestBlack(t *testing.T) {
	thisColor := systems.ParseHexColorFast("#000000")
	tr, tg, tb, ta := thisColor.RGBA()

	assert.Equal(t, uint32(0), tr>>8)
	assert.Equal(t, uint32(0), tg>>8)
	assert.Equal(t, uint32(0), tb>>8)
	assert.Equal(t, uint32(255), ta>>8)
}

func TestRandom(t *testing.T) {
	r := rand.Int63n(255)
	g := rand.Int63n(255)
	b := rand.Int63n(255)
	colorStr := "#" + toHex(r) + toHex(g) + toHex(b)
	thisColor := systems.ParseHexColorFast(colorStr)

	tr, tg, tb, ta := thisColor.RGBA()

	assert.Equal(t, uint32(r), tr>>8)
	assert.Equal(t, uint32(g), tg>>8)
	assert.Equal(t, uint32(b), tb>>8)
	assert.Equal(t, uint32(255), ta>>8)
}

func toHex(n int64) string {
	src := []byte{byte(n)}
	return hex.EncodeToString(src)
}
