package systems

import (
	"errors"
	"image/color"
	"math"
)

var errInvalidFormat = errors.New("invalid format")

func ParseHexColorFast(s string) color.RGBA {
	c := color.RGBA{}
	c.A = 0xff

	if s[0] != '#' {
		panic("Invalid color [" + s + "]")
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		panic("Invalid color [" + s + "]")
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		panic("Invalid color [" + s + "]")

	}
	return c
}

// Calculate lighter or darker color.
// luminosity goes from -1 (darker) to +1 (lighter)
func Lighten(c color.Color, luminosity float64) color.Color {

	r, g, b, a := c.RGBA()

	newColor := color.RGBA{
		R: calcLighteness(r, luminosity),
		G: calcLighteness(g, luminosity),
		B: calcLighteness(b, luminosity),
		A: uint8(a),
	}
	return newColor
}

func calcLighteness(c uint32, luminosity float64) uint8 {
	fc := float64(c) + float64(c)*luminosity
	cl := math.Round(math.Max(math.Min(255, fc), 0))
	return uint8(cl)

}
