package main

import (
	"fmt"
	"gopaint/gdiplus"
	. "gopaint/reza"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	win "github.com/lxn/win"
)

// DrawingCanvas is the main drawing canvas
type DrawingCanvas struct {
	// Embed the Window interface
	Window
	// Double buffer data
	mbitmap win.HBITMAP
	mhdc    win.HDC
	context *gdiplus.Graphics
	// Extras
	gridPen *Pen
	// Own data
	image *DrawingImage
	// test
	firstMove bool
	lastPt    Point
}

func NewDrawingCanvas(parent Window) *DrawingCanvas {
	canvas := &DrawingCanvas{Window: NewWindow()}
	canvas.Init(parent)
	return canvas
}

func (canvas *DrawingCanvas) Init(parent Window) {
	logInfo("initializing canvas...")
	canvas.firstMove = true
	canvas.gridPen = NewDashPen(1, NewRgb(120, 120, 120))
	canvas.Create("", win.WS_CHILD|win.WS_VISIBLE|win.WS_CLIPCHILDREN, 10, 10, 10, 10, parent)
	canvas.SetPaintEventHandler(canvas.Paint)
	canvas.SetMouseMoveEventHandler(canvas.MouseMove)
	canvas.SetMouseDownEventHandler(canvas.MouseDown)
	canvas.SetMouseUpEventHandler(canvas.MouseUp)
	canvas.SetMouseWheelEventHandler(func(e *MouseWheelEvent) {
		work := mainWindow.workspace
		if e.WheelDelta > 0 {
			work.ScrollUp()
		} else {
			work.ScrollDown()
		}
	})
	canvas.SetMouseLeaveEventHandler(func() {
		status := mainWindow.statusMousePos
		status.Update("")
	})
	canvas.SetKeyPressEventHandler(func(keycode int) {
		tool := mainWindow.tools.GetCurrentTool()
		if tool != nil {
			e := ToolKeyEvent{
				keycode: keycode,
				context: canvas.image.context,
				image:   canvas.image,
				canvas:  canvas,
			}
			tool.keyPressEvent(&e)
		}
	})
	canvas.SetSetCursorEventHandler(canvas.UpdateCursor)
	canvas.SetResizeEventHandler(canvas.OnResize)
	logInfo("Done initializing canvas")
}

func (canvas *DrawingCanvas) Dispose() {
	logInfo("Disposing canvas...")
	if canvas.image != nil {
		canvas.image.Dispose()
	}
	if canvas.context != nil {
		canvas.context.Dispose()
	}
	if canvas.mhdc != 0 {
		win.DeleteDC(canvas.mhdc)
	}
	if canvas.mbitmap != 0 {
		win.DeleteObject(win.HGDIOBJ(canvas.mbitmap))
	}
	canvas.Window.Dispose()
	if canvas.gridPen != nil {
		canvas.gridPen.Dispose()
	}
}

func (canvas *DrawingCanvas) UpdateCursor() bool {
	tool := mainWindow.tools.GetCurrentTool()
	if tool != nil {
		var winpt win.POINT
		win.GetCursorPos(&winpt)
		win.ScreenToClient(canvas.GetHandle(), &winpt)
		ptMouse := Point{X: int(winpt.X), Y: int(winpt.Y)}
		toolCursor := tool.getCursor(&ptMouse)
		win.SetCursor(toolCursor)
	} else {
		win.SetCursor(mainWindow.hCursorArrow)
	}
	return true
}

func (canvas *DrawingCanvas) NewImage(width, height int) {
	logInfo("'NewImage' - Clear the canvas and create a blank white image...")

	newImage := NewDrawingImage(width, height)
	logInfo("Clear...")
	newImage.filepath = ""
	newImage.Clear(gdiplus.NewColor(255, 255, 255, 255))
	/*
		pen := gdiplus.NewPen(gdiplus.NewColor(0, 0, 0, 255), 7)
		newImage.context.DrawLine(pen, 10.0, 210.0, 160.0, 400.0)
	*/
	//SavePNG("reza.png", newImage.hbitmap)
	/*
		color := Rgb(100, 0, 200)
		gc.SetFillColor(color.AsRGBA())

		gc.BeginPath() // Initialize a new path
		gc.ArcTo(80, 80, 50, 50, 0, math.Pi*2)
		gc.Close()
		gc.FillStroke()
	*/
	if canvas.image != nil {
		canvas.image.Dispose()
	}
	canvas.image = newImage
	canvas.SetSize(width, height)
	canvas.UpdateStatus()
}

func (canvas *DrawingCanvas) Resize(width, height int) {
	prevWidth := canvas.GetSize().Width
	prevHeight := canvas.GetSize().Height
	log.Printf("Resize from (%d x %d) to (%d x %d)\n", prevWidth, prevHeight, width, height)
	if prevWidth == width && prevHeight == height {
		return
	}
	if canvas.image == nil {
		log.Panicln("WHY is canvas image INVALID!!!")
	}
	// We allocate new data with given new size
	// then we copy/draw the old data/image into it
	newImage := NewDrawingImage(width, height)
	color := GetColorBackground()
	newImage.Clear(&color)

	rect2 := canvas.image.Bounds()
	draw.Draw(newImage, rect2, canvas.image, image.Point{}, draw.Src)

	newImage.filepath = canvas.image.filepath
	newImage.sizeOnDisk = canvas.image.sizeOnDisk
	newImage.lastSaved = canvas.image.lastSaved

	canvas.image.Dispose()
	canvas.image = newImage

	canvas.SetSize(width, height)
	canvas.UpdateStatus()
}

func (canvas *DrawingCanvas) OpenImage(filename string) bool {
	log.Printf("Open image '%s'...\n", filename)
	catFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return false
	}
	defer catFile.Close()

	// Obtain file size
	fi, err := catFile.Stat()
	if err != nil {
		// Could not obtain stat, handle error
		log.Println(err)
		return false
	}
	filesize := fi.Size()
	modDate := fi.ModTime().Format("02-Jan-06 3:04 PM")

	log.Println("Decoding...")
	var imageData image.Image
	ext := filepath.Ext(filename)
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		imageData, err = jpeg.Decode(catFile)

	} else {
		imageData, _, err = image.Decode(catFile)
	}
	if err != nil {
		fmt.Println(err)
		return false
	}
	log.Println("Done decoding!")

	bounds := imageData.Bounds()
	width, height := bounds.Size().X, bounds.Size().Y

	newImage := NewDrawingImage(width, height)
	logInfo("Copy image data into canvas image....")
	draw.Draw(newImage, newImage.Bounds(), imageData, image.Point{}, draw.Src)

	newImage.filepath = filename
	newImage.sizeOnDisk = filesize
	newImage.lastSaved = modDate
	if canvas.image != nil {
		canvas.image.Dispose()
	}
	canvas.image = newImage
	canvas.SetSize(width, height)
	canvas.UpdateStatus()
	canvas.Repaint()
	log.Println("Done opening image")
	return true
}

func (canvas *DrawingCanvas) SaveImage(filePath string) bool {
	log.Printf("Saving image '%s'...\n", filePath)
	fd, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return false
	}
	defer fd.Close()

	ext := filepath.Ext(filePath)
	format := FindFormatFromExt(ext)
	if format != nil {
		format.Function(fd, canvas.image, true)
	} else {
		log.Printf("Unknown file format/extension (%s)\n", ext)
		return false
	}
	log.Printf("Done Saving image\n")
	return true
}

func (canvas *DrawingCanvas) UpdateStatus() {
	size := canvas.GetSize()
	scz := mainWindow.statusCanvasSize
	if scz != nil {
		scz.Update(strconv.Itoa(size.Width) + " x " + strconv.Itoa(size.Height) + "px")
	}
	sfz := mainWindow.statusFileSize
	imageSize, available := canvas.image.SizeOnDisk()
	if available {
		sfz.SetVisible(true)
		sfz.Update(imageSize)
	} else {
		sfz.SetVisible(false)
	}
}

func (canvas *DrawingCanvas) UpdateMousePosStatus() {
	status := mainWindow.statusMousePos
	wndRect := canvas.GetWindowRect()
	ptScreen := app.GetCursorPos()
	if status != nil {
		if wndRect.IsPointInside(&ptScreen) {
			status.Update(strconv.Itoa(ptScreen.X) + ", " + strconv.Itoa(ptScreen.Y) + "px")
		} else {
			status.Update("")
		}
	}
}

func (canvas *DrawingCanvas) MouseDown(pt *Point, mbutton int) {
	win.SetCapture(canvas.GetHandle())
	tool := mainWindow.tools.GetCurrentTool()
	if canvas.firstMove {
		canvas.lastPt = *pt
		canvas.firstMove = false
	}

	e := ToolMouseEvent{
		pt:      *pt,
		lastPt:  canvas.lastPt,
		mbutton: mbutton,
		context: canvas.image.context,
		image:   canvas.image,
		canvas:  canvas,
	}
	tool.mouseDownEvent(&e)
	canvas.Repaint()
	canvas.lastPt = *pt
}

func (canvas *DrawingCanvas) MouseUp(pt *Point, mbutton int) {
	tool := mainWindow.tools.GetCurrentTool()
	if canvas.firstMove {
		canvas.lastPt = *pt
		canvas.firstMove = false
	}
	e := ToolMouseEvent{
		pt:      *pt,
		lastPt:  canvas.lastPt,
		mbutton: mbutton,
		context: canvas.image.context,
		image:   canvas.image,
		canvas:  canvas,
	}
	tool.mouseUpEvent(&e)
	canvas.RepaintVisible()
	canvas.lastPt = *pt
	win.ReleaseCapture()
}

func (canvas *DrawingCanvas) MouseMove(mousepoint *Point, mbutton int) {
	pt := *mousepoint
	tool := mainWindow.tools.GetCurrentTool()
	if canvas.firstMove {
		canvas.lastPt = pt
		canvas.firstMove = false
	}
	e := ToolMouseEvent{
		pt:      pt,
		lastPt:  canvas.lastPt,
		mbutton: mbutton,
		context: canvas.image.context,
		image:   canvas.image,
		canvas:  canvas,
	}
	tool.mouseMoveEvent(&e)
	canvas.UpdateMousePosStatus()
	canvas.RepaintVisible()
	canvas.lastPt = pt
}

func (canvas *DrawingCanvas) OnResize(rect *Rect) {
	logInfo("canvas resize...")
	if canvas.mhdc != 0 {
		win.DeleteDC(canvas.mhdc)
	}
	if canvas.mbitmap != 0 {
		win.DeleteObject(win.HGDIOBJ(canvas.mbitmap))
	}
	rcVisible := *rect ///canvas.GetVisibleRect()

	hdc := win.GetDC(canvas.GetHandle())
	canvas.mhdc = win.CreateCompatibleDC(hdc)
	canvas.mbitmap = win.CreateCompatibleBitmap(hdc, int32(rcVisible.Width()), int32(rcVisible.Height()))
	win.SelectObject(canvas.mhdc, win.HGDIOBJ(canvas.mbitmap))
	win.ReleaseDC(canvas.GetHandle(), hdc)

	brushBack := win.GetStockObject(win.BLACK_BRUSH)
	wrect := rcVisible.AsRECT()
	FillRect(canvas.mhdc, &wrect, win.HBRUSH(brushBack))
	if canvas.context != nil {
		canvas.context.Dispose()
	}
	canvas.context = gdiplus.NewGraphicsFromHDC(canvas.mhdc)
	//canvas.context.SetSmoothingMode(gdiplus.SmoothingModeHighSpeed)
}

// This returns only the visible area of the canvas, not the whole client rect
func (canvas *DrawingCanvas) GetVisibleRect() Rect {
	work := mainWindow.workspace
	if work == nil {
		return canvas.GetClientRect()
	}
	// Get Screen Space Rects
	rcWork := work.GetWindowRect()
	rcCanvas := canvas.GetWindowRect()

	// See if some portion of the left or top side of the canvas is hidden/we scrolled
	// the width and height of workspace area should be enough
	// ---(that just also includes the workspace scrollbars sizes too which is not a big deal, i hope :)

	canvasTotalWidth := rcCanvas.Width() + (canvasMargin * 2)   // mul by 2 for both sides
	canvasTotalHeight := rcCanvas.Height() + (canvasMargin * 2) // mul by 2 for both sides

	var rcVisible Rect
	// If the whole width fits the workspace we just assign the original width with left being 0
	if canvasTotalWidth <= rcWork.Width() {
		rcVisible.Right = rcCanvas.Width()
	} else {
		if rcCanvas.Left < rcWork.Left {
			rcVisible.Left = rcWork.Left - rcCanvas.Left
		}
		rcVisible.Right = rcVisible.Left + rcWork.Width()
	}
	// If the whole height fits the workspace we just assign the original height with top being 0
	if canvasTotalHeight <= rcWork.Height() {
		rcVisible.Bottom = rcCanvas.Height()
	} else {
		if rcCanvas.Top < rcWork.Left {
			rcVisible.Top = rcWork.Top - rcCanvas.Top
		}
		rcVisible.Bottom = rcVisible.Top + rcWork.Height()
	}
	return rcVisible
}

func (canvas *DrawingCanvas) RepaintVisible() {
	rect := canvas.GetVisibleRect()
	canvas.InvalidateRect(&rect, false)
}

func (canvas *DrawingCanvas) DrawGridLines(g *Graphics, rect *Rect, rcVisible *Rect) {
	g.SelectObject(canvas.gridPen)
	for x := 10; x <= rect.Right; x += 10 {
		if x > rcVisible.Left && x < rcVisible.Right {
			g.DrawLineOnly(x, 0, x, rect.Bottom)
		}
	}
	for y := 10; y <= rect.Bottom; y += 10 {
		if y > rcVisible.Top && y < rcVisible.Bottom {
			g.DrawLineOnly(0, y, rect.Right, y)
		}
	}
}

func (canvas *DrawingCanvas) Paint(g *Graphics, rect *Rect) {
	if canvas.context == nil {
		return
	}
	image := canvas.image
	if image == nil {
		return
	}
	rcVisible := canvas.GetVisibleRect()
	gmem := NewGraphics(canvas.mhdc)

	gmem.BitBlt(rcVisible.Left, rcVisible.Top,
		rcVisible.Width(), rcVisible.Height(), image.memdc,
		rcVisible.Left, rcVisible.Top, win.SRCCOPY)
	//gmem.AlphaBlend(rcVisible.Left, rcVisible.Top, rcVisible.Width(), rcVisible.Height(),
	//image.memdc, rcVisible.Left, rcVisible.Top, rcVisible.Width(), rcVisible.Height())

	if mainWindow.bShowGridlines.IsToggled() {
		canvas.DrawGridLines(gmem, rect, &rcVisible)
	}

	var winpt win.POINT
	win.GetCursorPos(&winpt)
	win.ScreenToClient(canvas.GetHandle(), &winpt)
	ptMouse := Point{X: int(winpt.X), Y: int(winpt.Y)}
	tool := mainWindow.tools.GetCurrentTool()
	if tool != nil {
		e := ToolDrawEvent{
			gdi32:    gmem,
			graphics: canvas.context,
			mouse:    ptMouse,
		}
		tool.draw(&e)
	}

	g.BitBlt(rcVisible.Left, rcVisible.Top,
		rcVisible.Width(), rcVisible.Height(),
		canvas.mhdc, rcVisible.Left, rcVisible.Top, win.SRCCOPY)
}
