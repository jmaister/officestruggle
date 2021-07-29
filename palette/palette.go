package palette

import (
	"image/color"
	"math"
)

type Hue int

const (
	Red       Hue = 0
	Flame     Hue = 15
	Orange    Hue = 20
	Amber     Hue = 45
	Yellow    Hue = 60
	Lime      Hue = 75
	Green     Hue = 120
	Emerald   Hue = 150
	Turquoise Hue = 165
	Cyan      Hue = 180
	Sky       Hue = 195
	Azure     Hue = 210
	Blue      Hue = 240
	Violet    Hue = 270
	Purple    Hue = 285
	Magenta   Hue = 300
	Pink      Hue = 315
	Crimson   Hue = 345
	Gray
	Sepia
)

// https://en.wikipedia.org/wiki/Hue
// HSL
func PColor(h Hue, brightness float64) color.Color {
	// H:   0
	// S: 240,
	// L:  30, 48, 66, 84, 102, 120, 144, 168, 192, 216
	if brightness < 0 || brightness > 1 {
		brightness = 0.5
	}

	H := float64(h) / 360.0
	S := 1.0
	L := (brightness*186 + 30) / 240.0

	return toRGB(H, S, L)
}

// https://stackoverflow.com/questions/2353211/hsl-to-rgb-color-conversion
func toRGB(h float64, s float64, l float64) color.Color {

	r := 0.0
	g := 0.0
	b := 0.0

	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		q := l * (1 + s)
		if q > 0.5 {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hueToRgb(p, q, h+1.0/3.0)
		g = hueToRgb(p, q, h)
		b = hueToRgb(p, q, h-1.0/3.0)
	}
	return color.RGBA{
		R: to255(r),
		G: to255(g),
		B: to255(b),
		A: uint8(255),
	}
}

func hueToRgb(p float64, q float64, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6.0
	}
	return p
}

func to255(v float64) uint8 {
	return uint8(math.Min(255, 256*v))
}
