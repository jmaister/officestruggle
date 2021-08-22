package systems

import (
	"image/color"
)

func ParseHexColorFast(s string) color.Color {
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

func ColorBlend(ca color.Color, cb color.Color, mix float64) color.Color {

	ra, ga, ba, _ := ca.RGBA()
	rb, gb, bb, _ := cb.RGBA()

	r1 := float64(ra)
	g1 := float64(ga)
	b1 := float64(ba)

	r2 := float64(rb)
	g2 := float64(gb)
	b2 := float64(bb)

	return color.RGBA{
		R: uint8((r1*mix + r2*(1-mix)) / 255),
		G: uint8((g1*mix + g2*(1-mix)) / 255),
		B: uint8((b1*mix + b2*(1-mix)) / 255),
		A: uint8(255),
	}
}
