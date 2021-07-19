package systems

import (
	"fmt"
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

	red := float64(ri / 255.0)
	green := float64(gi / 255.0)
	blue := float64(bi / 255.0)

	hue := 0.0
	saturation := 0.0

	max := math.Max(math.Max(red, green), blue)
	min := math.Min(math.Min(red, green), blue)

	fmt.Println(max, min)

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

	fmt.Println(hue, saturation, value)

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

	fmt.Println(red, green, blue)

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(math.Round(blue)),
		A: 0,
	}
}
