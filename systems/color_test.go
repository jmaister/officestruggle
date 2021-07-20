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

	assert.Equal(t, uint8(255), thisColor.R)
	assert.Equal(t, uint8(255), thisColor.G)
	assert.Equal(t, uint8(255), thisColor.B)
}

func TestBlack(t *testing.T) {
	thisColor := systems.ParseHexColorFast("#000000")

	assert.Equal(t, uint8(0), thisColor.R)
	assert.Equal(t, uint8(0), thisColor.G)
	assert.Equal(t, uint8(0), thisColor.B)
}

func TestRandom(t *testing.T) {
	r := rand.Int63n(255)
	g := rand.Int63n(255)
	b := rand.Int63n(255)
	colorStr := "#" + toHex(r) + toHex(g) + toHex(b)
	thisColor := systems.ParseHexColorFast(colorStr)

	assert.Equal(t, uint8(r), thisColor.R)
	assert.Equal(t, uint8(g), thisColor.G)
	assert.Equal(t, uint8(b), thisColor.B)
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
	assert.Equal(t, uint32(0), a)
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
	assert.Equal(t, uint32(0), a)
}
