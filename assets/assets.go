package assets

import (
	"embed"
	"fmt"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/*
var Data embed.FS

var mplusFontCached map[float64]font.Face

var fontCached = map[float64]font.Face{}

func MplusFont(size float64) font.Face {
	if mplusFontCached == nil {
		mplusFontCached = map[float64]font.Face{}
	}
	fnt, ok := mplusFontCached[size]
	if !ok {
		// tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
		tt, err := opentype.Parse(fonts.PressStart2P_ttf)
		if err != nil {
			panic(err)
		}

		fnt, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			panic(err)
		}
		mplusFontCached[size] = fnt
	}
	return fnt
}

func LoadFontCached(size float64) font.Face {
	fnt, ok := fontCached[size]
	if !ok {
		fnt = loadFont(size)
		fontCached[size] = fnt
	}
	return fnt
}

func loadFont(size float64) font.Face {
	//fontFile := "fonts/NotCourierSans.ttf"
	// fontFile := "fonts/DejaVuSansMono-Bold.ttf"
	//fontFile := "fonts/SLC_.ttf"
	//fontBytes, err := ioutil.ReadFile(fontFile)

	fontFile := "fonts/FreeMono.ttf"
	fontBytes, err := Data.ReadFile(fontFile)

	if err != nil {
		panic(err)
	}
	tt, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(fmt.Sprintf("Error parsing font bytes: %v", err))
	}

	return truetype.NewFace(tt, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
