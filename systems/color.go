package systems

import (
	"image/color"
	"math"
)

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

//https://github.com/google/closure-library/blob/master/closure/goog/color/color.js
func RGBToHSV(c color.Color) (int, float64, float64) {
	ri, gi, bi, _ := c.RGBA()

	red := float64(ri >> 8)
	green := float64(gi >> 8)
	blue := float64(bi >> 8)

	hue := 0.0
	saturation := 0.0

	max := math.Max(math.Max(red, green), blue)
	min := math.Min(math.Min(red, green), blue)

	value := max
	if min == max {
		hue = 0
		saturation = 0

	} else {
		delta := max - min
		saturation = delta / max

		if red == max {
			hue = (green - blue) / delta
		} else if green == max {
			hue = 2 + ((blue - red) / delta)
		} else {
			hue = 4 + ((red - green) / delta)
		}

		hue *= 60
		if hue < 0 {
			hue += 360

		}
		if hue > 360 {
			hue -= 360
		}
	}

	return int(hue), saturation, value
}

// https://github.com/google/closure-library/blob/master/closure/goog/color/color.js
func HSVToRGB(hi int, s float64, v float64) color.RGBA {
	red := 0.0
	green := 0.0
	blue := 0.0

	h := float64(hi)

	if s == 0 {
		red = v
		green = v
		blue = v
	} else {
		sextant := math.Floor(float64(h) / 60.0)
		remainder := (h / 60.0) - sextant
		val1 := v * (1.0 - s)
		val2 := v * (1.0 - (s * remainder))
		val3 := v * (1.0 - (s * (1 - remainder)))

		switch {
		case sextant == 1:
			red = val2
			green = v
			blue = val1
		case sextant == 2:
			red = val1
			green = v
			blue = val3
		case sextant == 3:
			red = val1
			green = val2
			blue = v
		case sextant == 4:
			red = val3
			green = val1
			blue = v
		case sextant == 5:
			red = v
			green = val1
			blue = val2
		case sextant == 6 || sextant == 0:
			red = v
			green = val3
			blue = val1
		}

	}

	return color.RGBA{
		R: uint8(red * 255),
		G: uint8(green * 255),
		B: uint8(blue * 255),
		A: uint8(255),
	}
}

func ColorBlend(ca color.RGBA, cb color.RGBA, mix float64) color.Color {

	ra, ga, ba, _ := ca.RGBA()
	rb, gb, bb, _ := cb.RGBA()

	r1 := float64(ra)
	g1 := float64(ga)
	b1 := float64(ba)

	r2 := float64(rb)
	g2 := float64(gb)
	b2 := float64(bb)

	return color.RGBA{
		R: uint8((r1*mix + r2*(1-mix)) / 256),
		G: uint8((g1*mix + g2*(1-mix)) / 255),
		B: uint8((b1*mix + b2*(1-mix)) / 255),
		A: uint8(255),
	}
}
