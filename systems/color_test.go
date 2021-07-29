package systems_test

import (
	"encoding/hex"
	"image/color"
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

func TestRGBToHSV(t *testing.T) {
	c := color.RGBA{
		R: 255,
		G: 255,
		B: 17,
		A: 0,
	}

	h, s, v := systems.RGBToHSV(c)
	assert.Equal(t, 60, h)
	assert.Assert(t, 0.93-s < 0.001)
	assert.Assert(t, 1-v < 0.001)
}

func TestHSVToRGB(t *testing.T) {
	h := 60
	s := 0.93
	v := 1.0

	c := systems.HSVToRGB(h, s, v)

	r, g, b, a := c.RGBA()

	assert.Equal(t, uint32(255), r>>8)
	assert.Equal(t, uint32(255), g>>8)
	assert.Equal(t, uint32(17), b>>8)
	assert.Equal(t, uint32(255), a>>8)
}

func TestHSVToRGB_2(t *testing.T) {
	h := 124
	s := 0.6
	v := 0.65

	c := systems.HSVToRGB(h, s, v)

	r, g, b, a := c.RGBA()

	assert.Equal(t, uint32(66), r>>8)
	assert.Equal(t, uint32(165), g>>8)
	assert.Equal(t, uint32(72), b>>8)
	assert.Equal(t, uint32(255), a>>8)
}
