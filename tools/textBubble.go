package tools

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// Expanded version of gg.Context.DrawRoundedRectangle(...)
func DrawRoundedRectangleThoughtBubble(dc *gg.Context, x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	dc.NewSubPath()
	dc.MoveTo(x1, y0)
	dc.LineTo(x2, y0)
	dc.DrawArc(x2, y1, r, gg.Radians(270), gg.Radians(360))
	dc.LineTo(x3, y2)
	dc.DrawArc(x2, y2, r, gg.Radians(0), gg.Radians(90))

	xA0, xA1, xA2 := w*.85, w*.8, w*.75
	yA0, yA1, yA2 := y+h, y+h+25, y+h
	dc.LineTo(xA0, yA0)
	dc.LineTo(xA1, yA1)
	dc.LineTo(xA2, yA2)

	dc.LineTo(x1, y3)
	dc.DrawArc(x1, y2, r, gg.Radians(90), gg.Radians(180))
	dc.LineTo(x0, y1)
	dc.DrawArc(x1, y1, r, gg.Radians(180), gg.Radians(270))
	dc.ClosePath()
}

func CreateTextBubble(width, height float64, inputText string, size float64) image.Image {
	// startTime := time.Now()
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic("")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: size,
	})

	dc := gg.NewContext(int(width), int(height))
	dc.SetFontFace(face)

	DrawRoundedRectangleThoughtBubble(dc, 3, 3, width-6, height-30, 10)
	dc.SetColor(White)
	dc.FillPreserve()
	dc.SetColor(SkyBlue)
	dc.SetLineWidth(4)
	dc.Stroke()

	dc.SetColor(Black)
	dc.DrawStringWrapped(inputText, 5, 5, 0.0, 0.0, width-10, 1, gg.AlignCenter)
	img := dc.Image()

	// fmt.Println(time.Since(startTime).Seconds()) // 1ms
	return img
}

func CreateTextImg(inputText string, width, height, size float64, c color.Color) image.Image {
	// startTime := time.Now()
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic("")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: size,
	})

	dc := gg.NewContext(int(width), int(height))
	dc.SetFontFace(face)
	dc.SetColor(c)
	dc.DrawStringWrapped(inputText, 0, 0, 0.0, 0.0, width, 1, gg.AlignCenter)
	img := dc.Image()

	// fmt.Println(time.Since(startTime).Seconds()) // 1ms
	return img
}
