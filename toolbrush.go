package main

import (
	. "gopaint/reza"
)

type ToolBrush struct {
	ToolBasic
	size int
}

func (tool *ToolBrush) initialize() {
	tool.size = 8
}

func (tool *ToolBrush) Dispose() {

}

func (tool *ToolBrush) prepare() {
	mainWindow.bsize.SetEnabled(true)
	menu := mainWindow.bsizeMenu
	items := menu.GetItems()
	for _, item := range items {
		(item.(PopupSizeMenuItem)).SetToggled(false)
	}
	(items[0].(PopupSizeMenuItem)).SetSize(1)
	(items[1].(PopupSizeMenuItem)).SetSize(3)
	(items[2].(PopupSizeMenuItem)).SetSize(5)
	(items[3].(PopupSizeMenuItem)).SetSize(8)
	switch tool.size {
	case 1:
		(items[0].(PopupSizeMenuItem)).SetToggled(true)
	case 3:
		(items[1].(PopupSizeMenuItem)).SetToggled(true)
	case 5:
		(items[2].(PopupSizeMenuItem)).SetToggled(true)
	case 8:
		(items[3].(PopupSizeMenuItem)).SetToggled(true)
	default:
		(items[3].(PopupSizeMenuItem)).SetToggled(true)
	}
}

func (tool *ToolBrush) changeSize(size int) {
	tool.size = size
}

func (tool *ToolBrush) draw(e *ToolDrawEvent) {
	pt := e.mouse
	g := e.gdi32
	gdipluscolor := getColorForeground()
	color := fromGdiplusColor(&gdipluscolor)
	halfSize := (tool.size + 1) / 2
	rect := &Rect{
		Left:   pt.X - halfSize,
		Top:    pt.Y - halfSize,
		Right:  pt.X + halfSize,
		Bottom: pt.Y + halfSize,
	}
	pen := NewSolidPen(1, &color)
	defer pen.Dispose()
	brush := NewSolidBrush(&color)
	defer brush.Dispose()
	g.DrawEllipseEx(rect, pen, brush)
}

func (tool *ToolBrush) mouseDownEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	//g := e.context
	gc := e.image.context2
	x := float64(e.pt.X)
	y := float64(e.pt.Y)
	//lastX := float64(e.lastPt.X)
	//lastY := float64(e.lastPt.Y)
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		color := getColorForeBack(mbutton)
		/*
			brush := gdiplus.NewSolidBrush(&color)
			ellipseSize := float32(tool.size) //(float32(tool.size) / 2.0) - 0.1
			halfSize := ellipseSize / 2
			g.FillEllipse(&brush.Brush, x-halfSize, y-halfSize, ellipseSize, ellipseSize)
			brush.Dispose()
		*/
		gc.SetRGBA255(int(color.GetB()), int(color.GetG()), int(color.GetR()), int(color.GetA()))
		gc.DrawPoint(x, y, (float64(tool.size)/2.0)-0.1)
		gc.Fill()
	}
}

func (tool *ToolBrush) mouseMoveEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	gc := e.image.context2
	x := float64(e.pt.X)
	y := float64(e.pt.Y)
	lastX := float64(e.lastPt.X)
	lastY := float64(e.lastPt.Y)
	//dx := x - lastX
	//dy := y - lastY
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		color := getColorForeBack(mbutton)
		/*
			brush := gdiplus.NewSolidBrush(&color)
			ellipseSize := float32(tool.size) //(float32(tool.size) / 2.0) - 0.1
			halfSize := ellipseSize / 2
			g.FillEllipse(&brush.Brush, x-halfSize, y-halfSize, ellipseSize, ellipseSize)
			brush.Dispose()
		*/
		gc.SetRGBA255(int(color.GetB()), int(color.GetG()), int(color.GetR()), int(color.GetA()))
		gc.SetLineWidth(float64(tool.size))
		gc.MoveTo(lastX+0.5, lastY+0.5)
		gc.LineTo(x+0.5, y+0.5)
		gc.Stroke()
	}
}

func (tool *ToolBrush) mouseUpEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	//gc := e.context
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {

	}
}
