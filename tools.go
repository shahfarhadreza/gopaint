package main

import (
	. "gopaint/reza"

	"github.com/shahfarhadreza/go-gdiplus"

	"github.com/fogleman/gg"
	"github.com/lxn/win"
)

type ToolMouseEvent struct {
	pt      Point
	lastPt  Point
	mbutton int
	context *gdiplus.Graphics
	image   *DrawingImage
	canvas  *DrawingCanvas
}

type ToolKeyEvent struct {
	keycode int
	context *gdiplus.Graphics
	image   *DrawingImage
	canvas  *DrawingCanvas
}

type ToolDrawEvent struct {
	gdi32    *Graphics
	graphics *gdiplus.Graphics
	mouse    Point
}

type Tool interface {
	initialize()
	Dispose()
	prepare()            // gets called everytime user chooses this tool
	leave()              // gets called everytime user switches from this to another tool
	changeSize(size int) // gets called when user changes size in the size dropdown button from the ribbon
	draw(e *ToolDrawEvent)
	getCursor(ptMouse *Point) win.HCURSOR
	mouseMoveEvent(e *ToolMouseEvent)
	mouseDownEvent(e *ToolMouseEvent)
	mouseUpEvent(e *ToolMouseEvent)
	keyPressEvent(e *ToolKeyEvent)
}

type ToolBasic struct {
}

func (tool *ToolBasic) keyPressEvent(e *ToolKeyEvent) {

}

func (tool *ToolBasic) changeSize(size int) {

}

func (tool *ToolBasic) leave() {

}

func (tool *ToolBasic) getCursor(ptMouse *Point) win.HCURSOR {
	return mainWindow.hCursorArrow
}

func AsGdiplusColor(color *Color) *gdiplus.Color {
	return gdiplus.NewColor(color.R, color.G, color.B, color.A)
}

func FromGdiplusColor(color *gdiplus.Color) Color {
	return Rgba(color.GetR(), color.GetG(), color.GetB(), color.GetA())
}

func GetColorForeground() (color gdiplus.Color) {
	c := mainWindow.color1.GetColor()
	color.Argb = gdiplus.MakeARGB(c.A, c.R, c.G, c.B)
	return
}

func GetColorBackground() (color gdiplus.Color) {
	c := mainWindow.color2.GetColor()
	color.Argb = gdiplus.MakeARGB(c.A, c.R, c.G, c.B)
	return
}

func GetColorForeBack(mouseButton int) (color gdiplus.Color) {
	if mouseButton == MouseButtonRight {
		return GetColorBackground()
	}
	return GetColorForeground()
}

func GetOutlineAndFillColors(mouseButton int) (outline, fill gdiplus.Color) {
	if mouseButton == MouseButtonRight {
		return GetColorBackground(), GetColorForeground()
	}
	return GetColorForeground(), GetColorBackground()
}

// brush and pen must be disposed by the caller
func GetPenAndBrush(mbutton int, penwidth float32) (*gdiplus.Pen, *gdiplus.Brush) {
	outline, fill := GetOutlineAndFillColors(mbutton)
	var pen *gdiplus.Pen = nil
	var brush *gdiplus.Brush = nil
	if mainWindow.menuSolidOutline.IsToggled() {
		pen = gdiplus.NewPen(&outline, penwidth)
	}
	if mainWindow.menuSolidFill.IsToggled() {
		brush = &gdiplus.NewSolidBrush(&fill).Brush
	}
	return pen, brush
}

func ggDrawPolygon(gc *gg.Context, points []Point) {
	gc.NewSubPath()
	for i, pt := range points {
		if i == 0 {
			gc.MoveTo(float64(pt.X), float64(pt.Y))
		} else {
			gc.LineTo(float64(pt.X), float64(pt.Y))
		}
	}
	gc.ClosePath()
}

func GetStartAndEnd(start, end Point) (s gdiplus.Point, e gdiplus.Point) {
	startPoint := start
	endPoint := end
	if startPoint.X > endPoint.X {
		startPoint.X, endPoint.X = endPoint.X, startPoint.X
	}
	if startPoint.Y > endPoint.Y {
		startPoint.Y, endPoint.Y = endPoint.Y, startPoint.Y
	}
	s.X, s.Y = int32(startPoint.X), int32(startPoint.Y)
	e.X, e.Y = int32(endPoint.X), int32(endPoint.Y)
	return
}
