package main

import (
	"unsafe"

	. "gopaint/reza"

	win "github.com/lxn/win"
)

const canvasMargin = 7

const ResizeTypeNone = 10
const ResizeTypeWidth = 11
const ResizeTypeHeight = 12
const ResizeTypeBoth = 13

// Workspace is a window panel where the drawing canvas gets placed
type Workspace struct {
	// Inherit data from the window type
	Window
	// Own data
	doubleBuffer   *DoubleBuffer
	canvas         *DrawingCanvas
	resizePreview  Window
	resizeType     int
	canvasPos      Point
	ptMouseDown    Point
	rcBoxBottom    Rect
	rcBoxCorner    Rect
	rcBoxRight     Rect
	xMinScroll     int
	xCurrentScroll int
	xMaxScroll     int
	yMinScroll     int
	yCurrentScroll int
	yMaxScroll     int
}

func NewWorkspace(parent Window) *Workspace {
	workspace := &Workspace{Window: NewWindow()}
	workspace.Init(parent)
	return workspace
}

func (work *Workspace) IsMouseOnResize() (resizeType int) {
	var winpt win.POINT
	win.GetCursorPos(&winpt)
	win.ScreenToClient(work.GetHandle(), &winpt)
	pt := Point{X: int(winpt.X), Y: int(winpt.Y)}
	if work.rcBoxBottom.IsPointInside(&pt) {
		return ResizeTypeHeight
	} else if work.rcBoxRight.IsPointInside(&pt) {
		return ResizeTypeWidth
	} else if work.rcBoxCorner.IsPointInside(&pt) {
		return ResizeTypeBoth
	}
	return ResizeTypeNone
}

func (work *Workspace) Init(parent Window) {
	logInfo("init workspace...")
	work.Create("", win.WS_CHILD|win.WS_CLIPCHILDREN|win.WS_VISIBLE|win.WS_VSCROLL|win.WS_HSCROLL, 10, 10, 10, 10, parent)
	work.SetMouseMoveEventHandler(work.MouseMove)
	work.SetMouseDownEventHandler(work.MouseDown)
	work.SetMouseUpEventHandler(work.MouseUp)
	work.SetPaintEventHandler(work.Paint)
	work.SetResizeEventHandler(work.OnResize)
	work.SetSetCursorEventHandler(work.updateCursor)
	work.SetHScrollEventHandler(work.HScroll)
	work.SetVScrollEventHandler(work.VScroll)

	color := Rgb(255, 255, 255)
	brushWhite := win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH))

	work.resizePreview = NewWindow()
	work.resizePreview.CreateEx("", win.WS_POPUP, win.WS_EX_LAYERED, 10, 10, 700, 100, work)
	work.resizePreview.SetPaintEventHandler(func(g *Graphics, rect *Rect) {
		wrect := rect.AsRECT()
		FillRect(g.GetHDC(), &wrect, brushWhite)
		g.DrawDashedRectangle(rect, NewRgb(0, 0, 0))
	})
	SetLayeredWindowAttributes(work.resizePreview.GetHandle(), color.AsCOLORREF(), 0, LWA_COLORKEY)

	work.canvas = NewDrawingCanvas(work)
	work.canvas.NewImage(app.DefaultCanvasSize.Width, app.DefaultCanvasSize.Height)
	//work.canvas.OpenImage(".\\images\\cloudy.jpg")
	// test
	//work.canvas.Resize(1100, 620)
	logInfo("Done initializing workspace")

}
func (work *Workspace) Dispose() {
	if work.canvas != nil {
		work.canvas.Dispose()
	}
	if work.doubleBuffer != nil {
		work.doubleBuffer.Dispose()
	}
}

func (work *Workspace) ScrollUp() {
	yNewPos := work.yCurrentScroll - 50
	work.UpdateVScroll(yNewPos)
}

func (work *Workspace) ScrollDown() {
	yNewPos := work.yCurrentScroll + 50
	work.UpdateVScroll(yNewPos)
}

func (work *Workspace) updateCursor() bool {
	resizeType := work.IsMouseOnResize()
	if resizeType == ResizeTypeHeight {
		win.SetCursor(mainWindow.hCursorSizeNS)
	} else if resizeType == ResizeTypeWidth {
		win.SetCursor(mainWindow.hCursorSizeWE)
	} else if resizeType == ResizeTypeBoth {
		win.SetCursor(mainWindow.hCursorSizeNWSE)
	} else {
		win.SetCursor(mainWindow.hCursorArrow)
	}
	return true
}

func (work *Workspace) UpdateVScroll(newValue int) {
	yNewPos := newValue
	// New position must be between 0 and the screen height.
	yNewPos = Max(0, yNewPos)
	yNewPos = Min(work.yMaxScroll, yNewPos)
	// If the current position does not change, do not scroll.
	if yNewPos == work.yCurrentScroll {
		return
	}
	// Reset the current scroll position.
	work.yCurrentScroll = yNewPos

	// Reset the scroll bar.
	var si win.SCROLLINFO
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = win.SIF_POS
	si.NPos = int32(work.yCurrentScroll)
	win.SetScrollInfo(work.GetHandle(), win.SB_VERT, &si, true)

	canvas := work.canvas

	if canvas != nil {
		work.canvasPos = Point{X: -(work.xCurrentScroll - canvasMargin), Y: -(work.yCurrentScroll - canvasMargin)}
		win.SetWindowPos(canvas.GetHandle(), 0,
			int32(work.canvasPos.X), int32(work.canvasPos.Y),
			0, 0, win.SWP_NOZORDER|win.SWP_NOSIZE|win.SWP_NOREDRAW)
		canvas.RepaintVisible()
		canvas.Update()
	}
	//work.Repaint()
	work.InvalidateRect(nil, false)
	work.Update()
}

func (work *Workspace) UpdateHScroll(newValue int) {
	// New position must be between 0 and the screen width.
	xNewPos := newValue
	xNewPos = Max(0, xNewPos)
	xNewPos = Min(work.xMaxScroll, xNewPos)
	// If the current position does not change, do not scroll.
	if xNewPos == work.xCurrentScroll {
		return
	}
	// Reset the current scroll position.
	work.xCurrentScroll = xNewPos

	// Reset the scroll bar.
	var si win.SCROLLINFO
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = win.SIF_POS
	si.NPos = int32(work.xCurrentScroll)
	win.SetScrollInfo(work.GetHandle(), win.SB_HORZ, &si, true)

	canvas := work.canvas

	if canvas != nil {
		work.canvasPos = Point{X: -(work.xCurrentScroll - canvasMargin), Y: -(work.yCurrentScroll - canvasMargin)}
		win.SetWindowPos(canvas.GetHandle(), 0,
			int32(work.canvasPos.X), int32(work.canvasPos.Y),
			0, 0, win.SWP_NOZORDER|win.SWP_NOSIZE|win.SWP_NOREDRAW)
		canvas.RepaintVisible()
		canvas.Update()
	}
	//work.Repaint()
	work.InvalidateRect(nil, false)
	work.Update()
}

func (work *Workspace) VScroll(stype, position int) {
	yNewPos := 0 // new position
	switch stype {
	case win.SB_PAGEUP:
		yNewPos = work.yCurrentScroll - 100
	case win.SB_PAGEDOWN:
		yNewPos = work.yCurrentScroll + 100
	case win.SB_LINEUP:
		yNewPos = work.yCurrentScroll - 50
	case win.SB_LINEDOWN:
		yNewPos = work.yCurrentScroll + 50
	case win.SB_THUMBPOSITION, win.SB_THUMBTRACK:
		yNewPos = position
	default:
		yNewPos = work.yCurrentScroll
	}
	work.UpdateVScroll(yNewPos)
}

func (work *Workspace) HScroll(stype, position int) {
	xNewPos := 0 // new position
	switch stype {
	case win.SB_PAGEUP:
		xNewPos = work.xCurrentScroll - 100
	case win.SB_PAGEDOWN:
		xNewPos = work.xCurrentScroll + 100
	case win.SB_LINEUP:
		xNewPos = work.xCurrentScroll - 50
	case win.SB_LINEDOWN:
		xNewPos = work.xCurrentScroll + 50
	case win.SB_THUMBPOSITION, win.SB_THUMBTRACK:
		xNewPos = position
	default:
		xNewPos = work.xCurrentScroll
	}
	work.UpdateHScroll(xNewPos)
}

func (work *Workspace) OnResize(clientNotused *Rect) {
	var si win.SCROLLINFO

	client := work.GetClientRect()

	if work.doubleBuffer != nil {
		work.doubleBuffer.Dispose()
	}
	work.doubleBuffer = NewDoubleBuffer(work, &client, &mainWindow.workspaceColor)

	canvasSize := work.canvas.GetSize()

	contentWidth := canvasSize.Width + (canvasMargin * 2)   // mul by 2 for both side margins
	contentHeight := canvasSize.Height + (canvasMargin * 2) // mul by 2 for both sides

	workspaceWidth := client.Width()
	work.xMaxScroll = Max(contentWidth-workspaceWidth, 0)
	work.xCurrentScroll = Min(work.xCurrentScroll, work.xMaxScroll)
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = win.SIF_RANGE | win.SIF_PAGE | win.SIF_POS
	si.NMin = int32(work.xMinScroll)
	si.NMax = int32(contentWidth)
	si.NPage = uint32(workspaceWidth)
	si.NPos = int32(work.xCurrentScroll)
	win.SetScrollInfo(work.GetHandle(), win.SB_HORZ, &si, true)

	// if both scrollbars are gonna show up then we need to recalculate our client area
	// otherwise scrollbar range won't be accurate since the scrollbars
	// themselves gonna eat some spaces
	client = work.GetClientRect()

	workspaceHeight := client.Height()
	work.yMaxScroll = Max(contentHeight-workspaceHeight, 0)
	work.yCurrentScroll = Min(work.yCurrentScroll, work.yMaxScroll)
	si.CbSize = uint32(unsafe.Sizeof(si))
	si.FMask = win.SIF_RANGE | win.SIF_PAGE | win.SIF_POS
	si.NMin = int32(work.yMinScroll)
	si.NMax = int32(contentHeight)
	si.NPage = uint32(workspaceHeight)
	si.NPos = int32(work.yCurrentScroll)
	win.SetScrollInfo(work.GetHandle(), win.SB_VERT, &si, true)

	if work.canvas != nil {
		work.canvasPos = Point{X: -(work.xCurrentScroll - canvasMargin), Y: -(work.yCurrentScroll - canvasMargin)}
		work.canvas.SetPosition(work.canvasPos.X, work.canvasPos.Y)
	}
}

func (work *Workspace) MouseDown(pt *Point, mbutton int) {
	win.SetCapture(work.GetHandle())
	work.resizeType = work.IsMouseOnResize()
	work.ptMouseDown = app.GetCursorPos()

	if work.resizeType != ResizeTypeNone {
		logInfo("Resize!!!")
		canvas := work.canvas
		preview := work.resizePreview
		wrect := work.GetWindowRect()
		crect := canvas.GetWindowRect()
		previewLeft := crect.Left
		diffLeft := 0
		previewWidth := crect.Width()
		previewTop := crect.Top
		diffTop := 0
		previewHeight := crect.Height()
		if previewLeft < wrect.Left {
			diffLeft = wrect.Left - previewLeft
			previewLeft = wrect.Left
		}
		if (crect.Left + previewWidth) > wrect.Right {
			previewWidth -= (crect.Left + previewWidth) - wrect.Right
		}
		if previewTop < wrect.Top {
			diffTop = wrect.Top - previewTop
			previewTop = wrect.Top
		}
		if (crect.Top + previewHeight) > wrect.Bottom {
			previewHeight -= (crect.Top + previewHeight) - wrect.Bottom
		}
		preview.MoveWindow(previewLeft, previewTop, previewWidth-diffLeft, previewHeight-diffTop, true)
		preview.SetVisible(true)
		preview.Update()
	}
}

func (work *Workspace) MouseMove(pt *Point, mbutton int) {
	if mbutton == MouseButtonLeft {
		canvas := work.canvas
		preview := work.resizePreview
		wrect := work.GetWindowRect()
		crect := canvas.GetWindowRect()
		ptMouse := app.GetCursorPos()
		ptMouseDiff := ptMouse.Distance(&work.ptMouseDown)
		newHeight := crect.Height() + ptMouseDiff.Y
		newWidth := crect.Width() + ptMouseDiff.X
		if newHeight < 1 {
			newHeight = 1
		}
		if newWidth < 1 {
			newWidth = 1
		}
		previewLeft := crect.Left
		diffLeft := 0
		previewWidth := crect.Width()
		previewTop := crect.Top
		diffTop := 0
		previewHeight := crect.Height()
		if previewLeft < wrect.Left {
			diffLeft = wrect.Left - previewLeft
			previewLeft = wrect.Left
		}
		if (crect.Left + previewWidth) > wrect.Right {
			previewWidth -= (crect.Left + previewWidth) - wrect.Right
		}
		if previewTop < wrect.Top {
			diffTop = wrect.Top - previewTop
			previewTop = wrect.Top
		}
		if (crect.Top + previewHeight) > wrect.Bottom {
			previewHeight -= (crect.Top + previewHeight) - wrect.Bottom
		}
		if work.resizeType == ResizeTypeHeight {
			preview.MoveWindow(previewLeft, previewTop, previewWidth-diffLeft, newHeight-diffTop, true)
			preview.Update()

		} else if work.resizeType == ResizeTypeWidth {
			preview.MoveWindow(previewLeft, previewTop, newWidth-diffLeft, previewHeight-diffTop, true)
			preview.Update()

		} else if work.resizeType == ResizeTypeBoth {
			preview.MoveWindow(previewLeft, previewTop, newWidth-diffLeft, newHeight-diffTop, true)
			preview.Update()
		}
	}
}

func (work *Workspace) MouseUp(pt *Point, mbutton int) {
	if mbutton == MouseButtonLeft {
		if work.resizeType != ResizeTypeNone {
			canvas := work.canvas
			crect := canvas.GetWindowRect()
			ptMouse := app.GetCursorPos()
			ptMouseDiff := ptMouse.Distance(&work.ptMouseDown)
			newHeight := crect.Height() + ptMouseDiff.Y
			newWidth := crect.Width() + ptMouseDiff.X
			if newHeight < 1 {
				newHeight = 1
			}
			if newWidth < 1 {
				newWidth = 1
			}
			if work.resizeType == ResizeTypeHeight {
				canvas.Resize(crect.Width(), newHeight)
			} else if work.resizeType == ResizeTypeWidth {
				canvas.Resize(newWidth, crect.Height())
			} else {
				canvas.Resize(newWidth, newHeight)
			}
			work.resizePreview.SetVisible(false)
			work.resizeType = ResizeTypeNone
			work.RequestLayout()
		}
	}
	win.ReleaseCapture()
}

func (work *Workspace) CalcResizeHandles(rect *Rect) {
	const size = 6
	work.rcBoxBottom.Left = rect.CenterX()
	work.rcBoxBottom.Top = rect.Bottom
	work.rcBoxBottom.Right = work.rcBoxBottom.Left + size
	work.rcBoxBottom.Bottom = work.rcBoxBottom.Top + size

	work.rcBoxRight.Left = rect.Right
	work.rcBoxRight.Top = rect.CenterY()
	work.rcBoxRight.Right = work.rcBoxRight.Left + size
	work.rcBoxRight.Bottom = work.rcBoxRight.Top + size

	work.rcBoxCorner.Left = rect.Right
	work.rcBoxCorner.Top = rect.Bottom
	work.rcBoxCorner.Right = work.rcBoxCorner.Left + size
	work.rcBoxCorner.Bottom = work.rcBoxCorner.Top + size
}

func (work *Workspace) Paint(gOrg *Graphics, rect *Rect) {
	if work.doubleBuffer == nil {
		return
	}
	db := work.doubleBuffer
	defer db.BitBlt(gOrg.GetHDC())

	g := db.GetGraphics()
	canvasSize := work.canvas.GetSize()
	rectCanvas := Rect{
		Left:   work.canvasPos.X,
		Top:    work.canvasPos.Y,
		Right:  work.canvasPos.X + canvasSize.Width,
		Bottom: work.canvasPos.Y + canvasSize.Height,
	}

	rectShadow := Rect{Left: (canvasMargin + 8) - work.xCurrentScroll, Top: (canvasMargin + 8) - work.yCurrentScroll}
	rectShadow.Right = rectShadow.Left + canvasSize.Width
	rectShadow.Bottom = rectShadow.Top + canvasSize.Height

	g.FillRect(&rectShadow, NewRgb(183, 194, 211))

	work.CalcResizeHandles(&rectCanvas)

	g.FillRectangle(&work.rcBoxBottom, NewRgb(70, 70, 70), NewRgb(255, 255, 255))
	g.FillRectangle(&work.rcBoxRight, NewRgb(70, 70, 70), NewRgb(255, 255, 255))
	g.FillRectangle(&work.rcBoxCorner, NewRgb(70, 70, 70), NewRgb(255, 255, 255))

}
