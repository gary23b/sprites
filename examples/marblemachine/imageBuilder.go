package main

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/gary23b/sprites"
)

func createCircleImage(radius float64, c color.Color) image.Image {
	dc := gg.NewContext(int(radius*3), int(radius*3))
	dc.DrawCircle(radius*1.5, radius*1.5, radius)
	dc.SetColor(c)
	dc.Fill()
	dc.SetColor(sprites.White)
	dc.SetLineWidth(3)
	dc.DrawLine(radius*1.5, radius*1.5, radius*2.0, radius*1.5)
	dc.Stroke()

	return dc.Image()
}

func createRectangleImage(width, height float64, c color.Color) image.Image {
	dc := gg.NewContext(int(width), int(height))
	dc.SetColor(c)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()
	return dc.Image()
}

func createSegmentImage(width, height float64, c color.Color) image.Image {
	hw := width / 2

	dc := gg.NewContext(int(width), int(height))
	dc.SetColor(c)

	dc.MoveTo(0, hw)
	dc.DrawArc(hw, hw, hw, 1.0*math.Pi, 2*math.Pi)
	dc.LineTo(width, height-hw)
	dc.DrawArc(hw, height-hw, hw, 0*math.Pi, 1*math.Pi)
	dc.LineTo(0, hw)
	dc.FillPreserve()
	dc.Stroke()

	return dc.Image()
}
