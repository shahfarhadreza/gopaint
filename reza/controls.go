package reza

import (
	"syscall"

	win "github.com/lxn/win"
)

const CheckboxAndTextGap = 6 // hardcoded constant. TODO: find a way to retriev this value from windows theme data..maybe?

type Label interface {
	Window
	SetAutoSize(enable bool)
}

type labelData struct {
	Window
	autoSize bool
}

type TextBox interface {
	Window
}

type Button interface {
	Window
	// Public
	SetClickEventHandler(f func(sender Button))
}

type buttonData struct {
	Window
	clickEventHandler func(sender Button)
}

type ImageViewer interface {
	Window
}

type imageViewerData struct {
	Window
}

func (btn *buttonData) SetClickEventHandler(f func(sender Button)) {
	btn.clickEventHandler = f
}

// Override and handle WM_COMMAND msg for the click event
func (button *buttonData) ReflectedMsg(reflectedFrom Window, msg uint32, wParam, lParam uintptr) {
	switch msg {
	case win.WM_COMMAND:
		if button.clickEventHandler != nil {
			button.clickEventHandler(button)
		}
	}
}

func (lbl *labelData) SetAutoSize(enable bool) {
	lbl.autoSize = enable
}

func CreateSubClassedWindow(iface Window, class, text string, x, y, width, height int, style uint, parent Window) Window {
	window := iface.asWindowData()
	window.parent = parent
	var handleParent win.HWND = 0
	if parent != nil {
		handleParent = parent.GetHandle()
		parentData := parent.asWindowData()
		parentData.addChild(window)
	}
	//window.Create(text, win.WS_TABSTOP|win.WS_VISIBLE|win.WS_CHILD|style, x, y, width, height, parent)
	classNameUTF16, _ := syscall.UTF16PtrFromString(class)
	textUTF16, _ := syscall.UTF16PtrFromString(text)

	styleEx := 0
	if class == "Edit" {
		styleEx = win.WS_EX_CLIENTEDGE
	}

	window.handle = win.CreateWindowEx(
		uint32(styleEx), classNameUTF16, textUTF16,
		uint32(win.WS_TABSTOP|win.WS_VISIBLE|win.WS_CHILD|style),
		int32(x), int32(y), int32(width), int32(height), handleParent, 0, app.hInstance, nil)

	window.orgWndProc = win.SetWindowLongPtr(window.handle, win.GWLP_WNDPROC, syscall.NewCallback(basicWindowProc))

	windowList[window.handle] = iface
	window.dockType = DockNone
	window.SetDefaultGuiFont()
	return iface
}

func CreateButton(text string, x, y, width, height int, style uint, parent Window) Button {
	btn := &buttonData{Window: NewWindow()}
	autoSize := false
	if width == 0 || height == 0 {
		autoSize = true
		width, height = 1, 1
	}
	CreateSubClassedWindow(btn, "Button", text, x, y, width, height, style, parent)
	if autoSize {
		// Calculate size according to the text
		hdc := win.GetDC(btn.GetHandle())
		g := NewGraphics(hdc)
		rect := g.MeasureText(text, win.DT_LEFT, btn.GetFont())

		if (style & win.BS_AUTORADIOBUTTON) != 0 {
			checkWidth := win.GetSystemMetrics(win.SM_CXMENUCHECK)
			rect.Right += int(checkWidth) + CheckboxAndTextGap
		} else if (style & win.BS_AUTOCHECKBOX) != 0 {
			checkWidth := win.GetSystemMetrics(win.SM_CXMENUCHECK)
			rect.Right += int(checkWidth) + CheckboxAndTextGap
		}
		win.ReleaseDC(btn.GetHandle(), hdc)
		btn.SetSize(rect.Width(), rect.Height())
	}
	return btn
}

func CreateGroup(text string, x, y, width, height int, parent Window) Button {
	btn := &buttonData{Window: NewWindow()}
	CreateSubClassedWindow(btn, "Button", text, x, y, width, height, win.BS_GROUPBOX, parent)
	return btn
}

func CreateLabel(text string, x, y, width, height int, parent Window) Label {
	lbl := &labelData{Window: NewWindow()}
	lbl.autoSize = false
	if width == 0 || height == 0 {
		lbl.autoSize = true
		width, height = 1, 1
	}
	CreateSubClassedWindow(lbl, "Static", text, x, y, width, height, 0, parent)
	if lbl.autoSize {
		// Calculate size according to the text
		hdc := win.GetDC(lbl.GetHandle())
		g := NewGraphics(hdc)
		rect := g.MeasureText(text, win.DT_LEFT, lbl.GetFont())
		win.ReleaseDC(lbl.GetHandle(), hdc)
		lbl.SetSize(rect.Width(), rect.Height())
	}
	return lbl
}

func (lbl *labelData) SetText(text string) {
	lbl.Window.SetText(text)
	if lbl.autoSize {
		// Calculate size according to the text
		hdc := win.GetDC(lbl.GetHandle())
		g := NewGraphics(hdc)
		rect := g.MeasureText(text, win.DT_LEFT, lbl.GetFont())
		win.ReleaseDC(lbl.GetHandle(), hdc)
		lbl.SetSize(rect.Width(), rect.Height())
	}
}

func CreateTextBox(text string, x, y, width, height int, style uint, parent Window) TextBox {
	wnd := NewWindow()
	CreateSubClassedWindow(wnd, "Edit", text, x, y, width, height, 0, parent)
	return wnd
}

func CreateImageViewer(imagepath string, x, y int, parent Window) ImageViewer {
	viewer := &imageViewerData{Window: NewWindow()}
	viewer.Create("", win.WS_CHILD|win.WS_CLIPCHILDREN|win.WS_VISIBLE, x, y, 1, 1, parent)
	img, _ := CreateBitmapImage(imagepath, false)
	if img != nil {
		viewer.SetSize(img.Width, img.Height)
	}
	viewer.SetPaintEventHandler(func(g *Graphics, rect *Rect) {
		if img != nil {
			g.DrawBitmapImage(img, 0, 0, false)
		}
	})
	return viewer
}
