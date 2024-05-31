package tools

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func MakeTextBubbleShape(w, h float64) image.Image {
	baseImage := ebiten.NewImage(int(w), int(h))
	baseImage.Fill(color.White)

	return baseImage
}

func CreateTextBubble(inputText string) image.Image {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 300
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	baseImage := ebiten.NewImage(100, 100)
	baseImage.Fill(color.White)
	text.Draw(baseImage, inputText, mplusNormalFont, 20, 20, color.Black)

	return baseImage
}
