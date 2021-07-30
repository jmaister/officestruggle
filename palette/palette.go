package palette

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

type Hue int

const (
	Red       Hue = 0
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
	Gray      Hue = -1
	Sepia     Hue = -2
)

var Hues = []Hue{
	Red,
	Orange,
	Amber,
	Yellow,
	Lime,
	Green,
	Emerald,
	Turquoise,
	Cyan,
	Sky,
	Azure,
	Blue,
	Violet,
	Purple,
	Magenta,
	Pink,
	Gray,
	Sepia,
}

var cache = map[float64]color.Color{}

// https://en.wikipedia.org/wiki/Hue
// HCL
func calculateColor(h Hue, brightness float64) color.Color {
	s := 0.85

	// brightness from 0.0 to 1.0 into 0.1 to 0.9 to avoid extremes
	if brightness < 0 || brightness > 1 {
		brightness = 0.5
	}
	l := brightness*(0.80-0.15) + 0.15

	if h == Gray {
		return colorful.Hsl(0.0, 0.0, l)
	} else if h == Sepia {
		return colorful.Hsl(30.0, 0.6, l)
	}

	return colorful.Hsl(float64(h), s, l)
}

func PColor(h Hue, brightness float64) color.Color {
	key := float64(h)*10.0 + brightness
	cl, ok := cache[key]
	if ok {
		return cl
	}

	cl = calculateColor(h, brightness)
	cache[key] = cl
	return cl
}
