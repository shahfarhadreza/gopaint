package main

import (
	. "gopaint/reza"
	"strconv"

	"github.com/shahfarhadreza/go-gdiplus"

	win "github.com/lxn/win"
)

const (
	SelectActionNone = iota
	SelectActionSelecting
	SelectActionMoving
	SelectActionResizing
)

type ToolSelect struct {
	ToolBasic
	penBorder     *Pen
	selection     *SelectionRect
	startPoint    Point
	currentAction int
	selected      bool
	bitmap        *BitmapGraphics
}

func (tool *ToolSelect) initialize() {
	tool.selection = NewSelectionRect()
	tool.penBorder = NewUserStylePen(1, NewRgb(0, 0, 0), []uint32{3, 4})
	tool.currentAction = SelectActionNone
}

func (tool *ToolSelect) Dispose() {
	if tool.penBorder != nil {
		tool.penBorder.Dispose()
	}
	if tool.selection != nil {
		tool.selection.Dispose()
	}
	if tool.bitmap != nil {
		tool.bitmap.Dispose()
	}
}

func (tool *ToolSelect) prepare() {
	tool.currentAction = SelectActionNone
	tool.selected = false
	tool.selection.Clear()
	if tool.bitmap != nil {
		tool.bitmap.Dispose()
		tool.bitmap = nil
	}
	//tool.SelectAll()
}

func (tool *ToolSelect) leave() {
	tool.finalizeSelection()
	tool.selected = false
	tool.currentAction = SelectActionNone
	tool.selection.Clear()
	tool.updateStatus()
	if tool.bitmap != nil {
		tool.bitmap.Dispose()
		tool.bitmap = nil
	}
}

func (tool *ToolSelect) getCursor(ptMouse *Point) win.HCURSOR {
	if tool.selected {
		rect := tool.selection.GetRect()
		if onpoint, point := tool.selection.GetClosestRectPoint(ptMouse, 6); onpoint {
			switch point {
			case RectPointTop, RectPointBottom:
				return mainWindow.hCursorSizeNS
			case RectPointLeft, RectPointRight:
				return mainWindow.hCursorSizeWE
			case RectPointTopLeft, RectPointBottomRight:
				return mainWindow.hCursorSizeNWSE
			case RectPointTopRight, RectPointBottomLeft:
				return mainWindow.hCursorSizeNESW
			}
		} else {
			if rect.IsPointInside(ptMouse) {
				return mainWindow.hCursorMove
			}
		}
	}
	return mainWindow.hCursorArrow
}

func (tool *ToolSelect) draw(e *ToolDrawEvent) {
	g := e.gdi32
	if tool.selection != nil {
		rect := tool.selection.GetRect()
		if !tool.selection.IsEmpty() {
			if tool.bitmap != nil {
				canvas := mainWindow.workspace.canvas
				visibleRect := canvas.GetVisibleRect()
				newRect := rect
				if newRect.Bottom > visibleRect.Bottom {
					newRect.Bottom = visibleRect.Bottom
				}
				if newRect.Right > visibleRect.Right {
					newRect.Right = visibleRect.Right
				}
				g.BitBlt(newRect.Left, newRect.Top,
					newRect.Width(), newRect.Height(), tool.bitmap.Hdc,
					0, 0, win.SRCCOPY)
			}
			if tool.currentAction == SelectActionSelecting || tool.currentAction == SelectActionMoving {
				g.DrawRectangleEx(&rect, tool.penBorder, nil)
			} else {
				tool.selection.Draw(g)
			}
		}
	}
}

func (tool *ToolSelect) SelectAll() {
	image := mainWindow.workspace.canvas.image
	tool.selected = true
	tool.currentAction = SelectActionNone
	tool.selection.SetRect(&Rect{
		Left:   0,
		Top:    0,
		Right:  image.Width(),
		Bottom: image.Height(),
	})
}

func (tool *ToolSelect) Deselect() {
	tool.finalizeSelection()
	tool.selected = false
	tool.currentAction = SelectActionNone
	tool.selection.Clear()
}

func (tool *ToolSelect) finalizeSelection() {
	image := mainWindow.workspace.canvas.image
	if tool.bitmap != nil {
		rect := tool.selection.GetRect()
		image.context3.BitBlt(rect.Left, rect.Top, rect.Width(), rect.Height(), tool.bitmap.Hdc, 0, 0, win.SRCCOPY)
		tool.bitmap.Dispose()
		tool.bitmap = nil
	}
}

func (tool *ToolSelect) DeleteSelection() {
	if tool.selection.IsEmpty() {
		return
	}
	image := mainWindow.workspace.canvas.image
	if tool.bitmap == nil {
		// replace the area with background color
		context := image.context
		color := GetColorBackground()
		brush := gdiplus.NewSolidBrush(&color)
		rect := tool.selection.GetRect()
		w, h := rect.Width(), rect.Height()
		context.FillRectangleI(brush.AsBrush(), int32(rect.Left), int32(rect.Top), int32(w), int32(h))
		brush.Dispose()
	} else {
		tool.bitmap.Dispose()
		tool.bitmap = nil
	}
	tool.selected = false
	tool.currentAction = SelectActionNone
	tool.selection.Clear()
}

func (tool *ToolSelect) mouseDownEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		tool.startPoint = e.pt
		rect := tool.selection.GetRect()
		if tool.selected {
			if onpoint, point := tool.selection.GetClosestRectPoint(&e.pt, 6); onpoint {
				switch point {
				case RectPointTop, RectPointBottom:
					tool.currentAction = SelectActionResizing
				case RectPointLeft, RectPointRight:
					tool.currentAction = SelectActionResizing
				case RectPointTopLeft, RectPointBottomRight:
					tool.currentAction = SelectActionResizing
				case RectPointTopRight, RectPointBottomLeft:
					tool.currentAction = SelectActionResizing
				}
			} else if rect.IsPointInside(&e.pt) {
				tool.currentAction = SelectActionMoving
				if !tool.selection.IsEmpty() {
					w, h := rect.Width(), rect.Height()
					if tool.bitmap == nil {
						//log.Printf("w %d, h %d, x %d y %d\n", tool.rect.Width(), tool.rect.Height(), tool.rect.Left, tool.rect.Top)
						tool.bitmap = NewBitmapGraphics(w, h)
						bitmapContext := tool.bitmap.Graphics
						bitmapContext.BitBlt(0, 0, w, h, e.image.context3.GetHDC(), rect.Left, rect.Top, win.SRCCOPY)
						// replace the area with background color
						context := e.image.context3
						gcolor := GetColorBackground()
						color := FromGdiplusColor(&gcolor)
						brush := NewSolidBrush(&color)
						pen := NewSolidPen(1, &color)
						context.FillRectangleEx(&rect, pen, brush)
						pen.Dispose()
						brush.Dispose()
					}
				}
			} else {
				tool.finalizeSelection()
				tool.selected = false
				tool.currentAction = SelectActionSelecting
				tool.selection.Clear()
			}
		} else {
			tool.selected = false
			tool.currentAction = SelectActionSelecting
			tool.selection.Clear()
		}
		tool.updateStatus()
	}
}

func (tool *ToolSelect) mouseMoveEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		if tool.currentAction == SelectActionSelecting {
			startPoint, endPoint := GetStartAndEnd(tool.startPoint, e.pt)
			rect := Rect{
				Left:   int(startPoint.X),
				Top:    int(startPoint.Y),
				Right:  int(endPoint.X),
				Bottom: int(endPoint.Y),
			}
			tool.selection.SetRect(&rect)
			tool.updateStatus()
		} else if tool.currentAction == SelectActionMoving {
			rect := tool.selection.GetRect()
			rectOrigin := Point{
				X: rect.Left,
				Y: rect.Top,
			}
			ptDist := e.pt.Distance(&e.lastPt)
			newPos := Point{
				X: rectOrigin.X + ptDist.X,
				Y: rectOrigin.Y + ptDist.Y,
			}
			rect = Rect{
				Left:   newPos.X,
				Top:    newPos.Y,
				Right:  newPos.X + rect.Width(),
				Bottom: newPos.Y + rect.Height(),
			}
			tool.selection.SetRect(&rect)
		}
	}
}

func (tool *ToolSelect) mouseUpEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	//image := e.image
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		if tool.currentAction == SelectActionSelecting {
			// if we have at least 1x1 selected
			if !tool.selection.IsEmpty() {
				tool.selected = true
			} else {
				tool.selected = false
			}
		}
		tool.currentAction = SelectActionNone
	}
}

func (tool *ToolSelect) updateStatus() {
	status := mainWindow.statusSelSize
	if status != nil {
		if !tool.selection.IsEmpty() {
			rect := tool.selection.GetRect()
			status.Update(strconv.Itoa(rect.Width()) + ", " + strconv.Itoa(rect.Height()) + "px")
		} else {
			status.Update("")
		}
	}
}
