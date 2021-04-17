package main

import (
	. "gopaint/reza"
)

type ToolEraser struct {
	ToolBasic
	size int
}

func (tool *ToolEraser) initialize() {
	tool.size = 5
}

func (tool *ToolEraser) Dispose() {

}

func (tool *ToolEraser) prepare() {
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

func (tool *ToolEraser) changeSize(size int) {
	tool.size = size
}

func (tool *ToolEraser) draw(e *ToolDrawEvent) {

}

func (tool *ToolEraser) mouseDownEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	gc := e.image.context2
	x := float64(e.pt.X)
	y := float64(e.pt.Y)
	if mbutton == MouseButtonLeft {
		color := getColorBackground()
		gc.SetRGBA255(int(color.GetB()), int(color.GetG()), int(color.GetR()), int(color.GetA()))
		gc.DrawPoint(x, y, (float64(tool.size)/2.0)-0.1)
		gc.Fill()
	}
}

func (tool *ToolEraser) mouseMoveEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	gc := e.image.context2
	x := float64(e.pt.X)
	y := float64(e.pt.Y)
	lastX := float64(e.lastPt.X)
	lastY := float64(e.lastPt.Y)
	if mbutton == MouseButtonLeft {
		color := getColorBackground()
		gc.SetRGBA255(int(color.GetB()), int(color.GetG()), int(color.GetR()), int(color.GetA()))
		gc.SetLineWidth(float64(tool.size))
		gc.MoveTo(lastX, lastY)
		gc.LineTo(x, y)
		gc.Stroke()
	}
}

func (tool *ToolEraser) mouseUpEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	//gc := e.context
	if mbutton == MouseButtonLeft {

	}
}
