package game

import (
	"image"
	"image/color"
	"log"

	"github.com/gary23b/sprites/imagedraw"
	"github.com/gary23b/sprites/tools"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func MakeTextBubbleShape(w, h float64) image.Image {
	i := imagedraw.NewImageDraw(int(w), int(h))
	i.Color(tools.White)
	//i.CircleAroundPoint(0, 0, 10, 0, 360, 20)
	i.FakeEllipseAroundPoint(0, 0, w*0.4, h*0.4, 0, 100)
	i.BucketFillPoint(0, 0, tools.Water)

	return i.GetImage()
}

func CreateTextBubble(inputText string, size int) image.Image {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}
	mplusNormalFont = text.FaceWithLineHeight(mplusNormalFont, float64(size))

	baseImage := ebiten.NewImageFromImage(MakeTextBubbleShape(200, 100))
	text.Draw(baseImage, inputText, mplusNormalFont, 30, 50, color.Black)

	return baseImage
}
