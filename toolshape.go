package main

import (
	. "gopaint/reza"

	"github.com/shahfarhadreza/go-gdiplus"
)

type ShapeDrawer interface {
	draw(args *ToolDrawShapeArgs)
}

type ToolShape struct {
	ToolBasic
	ShapeDrawer
	strokeWidth int
	startPoint  Point
	endPoint    Point
	isDrawing   bool
	mbutton     int
}

type ToolDrawShapeArgs struct {
	gdi32         *Graphics
	context       *gdiplus.Graphics
	pen           *gdiplus.Pen
	brush         *gdiplus.Brush
	rect          gdiplus.Rect
	startPointOrg gdiplus.Point
	endPointOrg   gdiplus.Point
	startPoint    gdiplus.Point
	endPoint      gdiplus.Point
}

func (tool *ToolShape) initialize() {
	tool.strokeWidth = 1
	tool.isDrawing = false
}

func (tool *ToolShape) Dispose() {

}

func (tool *ToolShape) prepare() {
	tool.isDrawing = false
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
	switch tool.strokeWidth {
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

func (tool *ToolShape) changeSize(size int) {
	tool.strokeWidth = size
}

func (tool *ToolShape) draw(e *ToolDrawEvent) {
	g := e.graphics
	mbutton := tool.mbutton
	if tool.isDrawing {
		startPointOrg, endPointOrg := gdiplus.Point{
			X: int32(tool.startPoint.X), Y: int32(tool.startPoint.Y),
		}, gdiplus.Point{
			X: int32(tool.endPoint.X), Y: int32(tool.endPoint.Y),
		}
		startPoint, endPoint := GetStartAndEnd(tool.startPoint, tool.endPoint)
		rect := gdiplus.Rect{
			X:      startPoint.X,
			Y:      startPoint.Y,
			Width:  endPoint.X - startPoint.X,
			Height: endPoint.Y - startPoint.Y,
		}
		pen, brush := GetPenAndBrush(mbutton, float32(tool.strokeWidth))
		if pen == nil && brush == nil {
			return
		}
		args := &ToolDrawShapeArgs{
			gdi32:         e.gdi32,
			context:       g,
			pen:           pen,
			brush:         brush,
			rect:          rect,
			startPoint:    startPoint,
			endPoint:      endPoint,
			startPointOrg: startPointOrg,
			endPointOrg:   endPointOrg,
		}
		tool.ShapeDrawer.draw(args)
		if pen != nil {
			pen.Dispose()
		}
		if brush != nil {
			brush.Dispose()
		}
	}
}

func (tool *ToolShape) mouseDownEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		tool.startPoint = e.pt
		tool.endPoint = e.pt
		tool.isDrawing = true
		tool.mbutton = mbutton
	}
}

func (tool *ToolShape) mouseMoveEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight && tool.isDrawing {
		tool.endPoint = e.pt
	}
}

func (tool *ToolShape) mouseUpEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	g := e.image.context
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		if tool.isDrawing {
			tool.isDrawing = false
			tool.endPoint = e.pt
			startPointOrg, endPointOrg := gdiplus.Point{
				X: int32(tool.startPoint.X), Y: int32(tool.startPoint.Y),
			}, gdiplus.Point{
				X: int32(tool.endPoint.X), Y: int32(tool.endPoint.Y),
			}
			startPoint, endPoint := GetStartAndEnd(tool.startPoint, tool.endPoint)
			rect := gdiplus.Rect{
				X:      startPoint.X,
				Y:      startPoint.Y,
				Width:  endPoint.X - startPoint.X,
				Height: endPoint.Y - startPoint.Y,
			}
			pen, brush := GetPenAndBrush(mbutton, float32(tool.strokeWidth))
			if pen == nil && brush == nil {
				return
			}
			args := &ToolDrawShapeArgs{
				gdi32:         e.image.context3,
				context:       g,
				pen:           pen,
				brush:         brush,
				rect:          rect,
				startPoint:    startPoint,
				endPoint:      endPoint,
				startPointOrg: startPointOrg,
				endPointOrg:   endPointOrg,
			}
			tool.ShapeDrawer.draw(args)
			if pen != nil {
				pen.Dispose()
			}
			if brush != nil {
				brush.Dispose()
			}
		}
	}
}
