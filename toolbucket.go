package main

import (
	. "gopaint/reza"
	"image/color"
)

type ToolBucket struct {
	ToolBasic
}

func (tool *ToolBucket) initialize() {

}

func (tool *ToolBucket) Dispose() {

}

func (tool *ToolBucket) prepare() {

}

func (tool *ToolBucket) draw(e *ToolDrawEvent) {

}

type PointStack struct {
	top  *PointElement
	size int
}

type PointElement struct {
	value Point
	next  *PointElement
}

func (s *PointStack) IsEmpty() bool {
	return s.size == 0
}

func (s *PointStack) Push(value Point) {
	s.top = &PointElement{value, s.top}
	s.size++
}

func (s *PointStack) Pop() (value Point) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return Point{}
}

// BGRA
func fastSet(p *DrawingImage, i int, c *color.RGBA) {
	p.Pix[i+0] = c.B
	p.Pix[i+1] = c.G
	p.Pix[i+2] = c.R
	p.Pix[i+3] = c.A
}

// BGRA
func hasSimilarColor(p *DrawingImage, i int, c *Color, tolerance int) bool {
	/*
		c1 := int(math.Abs(float64(c.R-p.Pix[i+2]))) <= tolerance
		c2 := int(math.Abs(float64(c.G-p.Pix[i+1]))) <= tolerance
		c3 := int(math.Abs(float64(c.B-p.Pix[i+0]))) <= tolerance
		c4 := int(math.Abs(float64(c.A-p.Pix[i+3]))) <= tolerance
		return c1 && c2 && c3 && c4*/
	return (c.R == p.Pix[i+2] && c.G == p.Pix[i+1] && c.B == p.Pix[i+0] && c.A == p.Pix[i+3])
}

func fastPixOffset(p *DrawingImage, x, y int) int {
	return y*p.Stride + x*4
}

func floodFillScanline(image *DrawingImage, q Point, oldColor, newColor Color) {
	var x1 int
	var spanAbove, spanBelow bool
	w, h := image.Width(), image.Height()
	stack := new(PointStack)
	stack.Push(q)
	newColorRgba := newColor.AsRGBA()
	tolerance := 50
	for !stack.IsEmpty() {
		p := stack.Pop()
		y := p.Y
		x1 = p.X
		for {
			pixelIndex := fastPixOffset(image, x1, y)
			if !(x1 >= 0 && hasSimilarColor(image, pixelIndex, &oldColor, tolerance)) {
				break
			}
			x1--
		}
		x1++
		spanAbove, spanBelow = false, false
		for {
			pixelIndex := fastPixOffset(image, x1, y)
			if !(x1 < w && hasSimilarColor(image, pixelIndex, &oldColor, tolerance)) {
				break
			}
			fastSet(image, pixelIndex, &newColorRgba)
			pixelIndex = fastPixOffset(image, x1, (y - 1))
			if !spanAbove && y > 0 && hasSimilarColor(image, pixelIndex, &oldColor, tolerance) {
				stack.Push(Point{X: x1, Y: y - 1})
				spanAbove = true
			} else if spanAbove && y > 0 && !hasSimilarColor(image, pixelIndex, &oldColor, tolerance) {
				spanAbove = false
			}
			pixelIndex = fastPixOffset(image, x1, (y + 1))
			if !spanBelow && y < h-1 && hasSimilarColor(image, pixelIndex, &oldColor, tolerance) {
				stack.Push(Point{X: x1, Y: y + 1})
				spanBelow = true
			} else if spanBelow && y < h-1 && !hasSimilarColor(image, pixelIndex, &oldColor, tolerance) {
				spanBelow = false
			}
			x1++
		}
	}
	stack = nil
}

func (tool *ToolBucket) mouseDownEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	image := e.image
	x := e.pt.X
	y := e.pt.Y
	clickedPixel := image.GetColorAt(x, y)
	c := getColorForeBack(mbutton)
	floodWith := Rgba(c.GetR(), c.GetG(), c.GetB(), c.GetA())
	floodFillScanline(image, Point{X: x, Y: y}, clickedPixel, floodWith)
}

func (tool *ToolBucket) mouseMoveEvent(e *ToolMouseEvent) {

}

func (tool *ToolBucket) mouseUpEvent(e *ToolMouseEvent) {

}
