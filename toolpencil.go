package main

import (
	. "gopaint/reza"
)

type ToolPencil struct {
	ToolBasic
	size int
}

func (tool *ToolPencil) initialize() {
	tool.size = 1
}

func (tool *ToolPencil) Dispose() {

}

func (tool *ToolPencil) prepare() {
	mainWindow.bsize.SetEnabled(true)
	menu := mainWindow.bsizeMenu
	items := menu.GetItems()
	for _, item := range items {
		(item.(PopupSizeMenuItem)).SetToggled(false)
	}
	(items[0].(PopupSizeMenuItem)).SetSize(1)
	(items[1].(PopupSizeMenuItem)).SetSize(2)
	(items[2].(PopupSizeMenuItem)).SetSize(3)
	(items[3].(PopupSizeMenuItem)).SetSize(4)
	switch tool.size {
	case 1:
		(items[0].(PopupSizeMenuItem)).SetToggled(true)
	case 2:
		(items[1].(PopupSizeMenuItem)).SetToggled(true)
	case 3:
		(items[2].(PopupSizeMenuItem)).SetToggled(true)
	case 4:
		(items[3].(PopupSizeMenuItem)).SetToggled(true)
	default:
		(items[3].(PopupSizeMenuItem)).SetToggled(true)
	}
}

func (tool *ToolPencil) changeSize(size int) {
	tool.size = size
}

func (tool *ToolPencil) draw(e *ToolDrawEvent) {

}

func (tool *ToolPencil) mouseMoveEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	gc := e.image.context2
	x := float64(e.pt.X)
	y := float64(e.pt.Y)
	lastX := float64(e.lastPt.X)
	lastY := float64(e.lastPt.Y)
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		color := getColorForeBack(mbutton)
		gc.SetRGBA255(int(color.GetB()), int(color.GetG()), int(color.GetR()), int(color.GetA()))
		gc.SetLineWidth(float64(tool.size))
		gc.MoveTo(lastX+0.5, lastY+0.5)
		gc.LineTo(x+0.5, y+0.5)
		gc.Stroke()
	}
}

func (tool *ToolPencil) mouseDownEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	gc := e.image.context2
	x := float64(e.pt.X)
	y := float64(e.pt.Y)
	lastX := float64(e.lastPt.X)
	lastY := float64(e.lastPt.Y)
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		color := getColorForeBack(mbutton)
		gc.SetRGBA255(int(color.GetB()), int(color.GetG()), int(color.GetR()), int(color.GetA()))
		gc.SetLineWidth(float64(tool.size))
		gc.MoveTo(lastX+0.5, lastY+0.5)
		gc.LineTo(x+0.5, y+0.5)
		gc.Stroke()
	}
}

func (tool *ToolPencil) mouseUpEvent(e *ToolMouseEvent) {

}
