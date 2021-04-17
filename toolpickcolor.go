package main

import (
	. "gopaint/reza"
)

type ToolPickColor struct {
	ToolBasic
}

func (tool *ToolPickColor) initialize() {

}

func (tool *ToolPickColor) Dispose() {

}

func (tool *ToolPickColor) prepare() {

}

func (tool *ToolPickColor) changeSize(size int) {

}

func (tool *ToolPickColor) draw(e *ToolDrawEvent) {

}

func (tool *ToolPickColor) mouseDownEvent(e *ToolMouseEvent) {

}

func (tool *ToolPickColor) mouseMoveEvent(e *ToolMouseEvent) {

}

func (tool *ToolPickColor) mouseUpEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	image := e.image
	x := e.pt.X
	y := e.pt.Y
	pickedColor := image.GetColorAt(x, y)
	if mbutton == MouseButtonLeft {
		mainWindow.color1.SetColor(pickedColor)

	} else if mbutton == MouseButtonRight {
		mainWindow.color2.SetColor(pickedColor)
	}
}
