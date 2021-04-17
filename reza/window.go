package reza

import (
	"syscall"
	"unsafe"

	win "github.com/lxn/win"
)

var windowBeingCreated Window = nil
var windowList map[win.HWND]Window = make(map[win.HWND]Window)

type Margins struct {
	Left, Right, Top, Bottom int
}

const (
	// MouseButtonNone means no button pressed
	MouseButtonNone = 0
	// MouseButtonLeft means left button pressed
	MouseButtonLeft = 1
	// MouseButtonMiddle means middle button pressed
	MouseButtonMiddle = 2
	// MouseButtonRight means right button pressed
	MouseButtonRight = 3
)

// Dock Type
const (
	DockNone   = 0
	DockTop    = 1
	DockBottom = 2
	DockLeft   = 3
	DockRight  = 4
	DockFill   = 5
)

type MouseWheelEvent struct {
	WheelDelta    int
	MousePosition Point
	VirtualKey    int
}

// Window is a common access interface to any Window/Win32 object/control
type Window interface {
	GetHandle() win.HWND
	Create(text string, style uint, x, y, width, height int, parent Window)
	CreateEx(text string, style, styleEx uint, x, y, width, height int, parent Window)
	Dispose()
	MoveWindow(x, y, width, height int, repaint bool)
	GetSize() Size
	GetClientRect() Rect
	GetWindowRect() Rect
	SetPosition(x, y int) bool
	SetSize(width, height int) bool
	Repaint()
	SetVisible(visible bool)
	SetText(text string)
	SetFont(font win.HFONT)
	GetFont() win.HFONT
	GetText() string
	SetMargin(left, right, top, bottom int)
	SetDockType(dtype int)
	GetDockType() int
	GetMargin() (left, right, top, bottom int)
	GetChildrens() []Window
	HasParent() bool
	GetParent() Window
	GetDC() win.HDC
	ReleaseDC(hdc win.HDC)
	RequestLayout()
	InvalidateRect(rect *Rect, eraseBackground bool)
	Update()
	// Overridables
	ReflectedMsg(reflectedFrom Window, msg uint32, wParam, lParam uintptr)
	// Event Handlers
	SetPaintEventHandler(func(g *Graphics, rc *Rect))
	SetMouseMoveEventHandler(func(pt *Point, mbutton int))
	SetMouseDownEventHandler(func(pt *Point, mbutton int))
	SetMouseUpEventHandler(func(pt *Point, mbutton int))
	SetMouseLeaveEventHandler(func())
	SetMouseWheelEventHandler(func(e *MouseWheelEvent))
	SetKeyDownEventHandler(func(keycode int))
	SetKeyUpEventHandler(func(keycode int))
	SetKeyPressEventHandler(func(keycode int))
	SetResizeEventHandler(func(client *Rect))
	SetKillFocusEventHandler(func())
	SetDestroyEventHandler(func())
	SetCloseEventHandler(func() bool)
	SetPositionChangedEventHandler(func(wp *win.WINDOWPOS))
	SetMouseActivateEventHandler(func() uintptr)
	SetSetCursorEventHandler(func() bool)
	SetHScrollEventHandler(func(stype, position int))
	SetVScrollEventHandler(func(stype, position int))
	// only for internal use
	asWindowData() *windowData
}

// windowData is a win32 window control
type windowData struct {
	// Private data
	handle             win.HWND
	orgWndProc         uintptr
	parent             Window
	childrens          []Window // childrens
	dockType           int
	marginLeft         int
	marginRight        int
	marginTop          int
	marginBottom       int
	font               win.HFONT
	isMouseInside      bool
	quickAccessToolbar bool
	// All the events
	paintHandler         func(g *Graphics, rc *Rect)
	MouseMoveHandler     func(pt *Point, mbutton int)
	MouseDownEvent       func(pt *Point, mbutton int)
	MouseUpEvent         func(pt *Point, mbutton int)
	MouseLeaveEvent      func()
	MouseWheelEvent      func(e *MouseWheelEvent)
	KeyDownEvent         func(keycode int)
	KeyUpEvent           func(keycode int)
	keyPressHandler      func(keycode int)
	ResizeEvent          func(client *Rect)
	KillFocusEvent       func()
	DestroyEvent         func()
	CloseEvent           func() bool
	PositionChangedEvent func(wp *win.WINDOWPOS)
	MouseActivateEvent   func() uintptr
	SetCursorEvent       func() bool
	HScrollEvent         func(stype, position int)
	VScrollEvent         func(stype, position int)
}

func NewWindow() Window {
	return &windowData{}
}

func (window *windowData) SetDefaultGuiFont() {
	appGuiFont := app.GetGuiFont()
	if appGuiFont != nil {
		window.font = appGuiFont.GetHandle()
	} else {
		window.font = win.HFONT(win.GetStockObject(win.DEFAULT_GUI_FONT))
	}
	win.SendMessage(window.handle, win.WM_SETFONT, uintptr(window.font), 1)
}

// windowData simply returns the receiver.
func (window *windowData) asWindowData() *windowData {
	return window
}

func (window *windowData) createWindowImpl(text string, style, styleEx uint, x, y, width, height int, parent Window) {
	logInfo("'createWindowImpl'")

	classNameUTF16, _ := syscall.UTF16PtrFromString("reza-ui-window-class")
	textUTF16, _ := syscall.UTF16PtrFromString(text)

	windowBeingCreated = window

	registerWindow(classNameUTF16, basicWindowProc)

	window.parent = parent
	var handleParent win.HWND = 0
	if parent != nil {
		handleParent = parent.GetHandle()
		parentData := parent.asWindowData()
		parentData.addChild(window)
	}

	logInfo("'CreateWindowEx'.....")
	window.handle = win.CreateWindowEx(
		uint32(styleEx), classNameUTF16, textUTF16,
		uint32(style), int32(x), int32(y), int32(width), int32(height), handleParent, 0, app.hInstance, nil)
	logInfo("Done with 'CreateWindowEx'!!!!")

	window.dockType = DockNone
	// Set a nice font (Default gui font)
	window.SetDefaultGuiFont()

	windowBeingCreated = nil
}

func (window *windowData) Dispose() {
	if window.handle != 0 {
		win.DestroyWindow(window.handle)
	}
	delete(windowList, window.handle)
}

func basicWindowProc(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	iwindow := windowList[hWnd]
	if iwindow == nil {
		// window is being created through 'CreateWindowEx', not added to the list yet. add it
		if windowBeingCreated == nil {
			panic("Serious bug in the application!!! - in func basicWindowProc, 'windowBeingCreated' is nil")
		}
		windowList[hWnd] = windowBeingCreated
		iwindow = windowBeingCreated
	}
	if iwindow == nil {
		panic("Serious bug in the application!!! - in func basicWindowProc, 'window' is nil")
	}
	window := iwindow.asWindowData()

	//fCallDWP := true
	var ret uintptr
	//if window.quickAccessToolbar {
	//ret = customWindowProc(window, hWnd, msg, wParam, lParam, &fCallDWP)
	//}
	// Winproc worker for the rest of the application.
	//if fCallDWP {
	ret = appWindowProc(window, hWnd, msg, wParam, lParam)
	//}
	return ret
}

/*
func customWindowProc(window *windowData, hWnd win.HWND, msg uint32, wParam, lParam uintptr, pfCallDWP *bool) uintptr {
	fCallDWP := true // Pass on to DefWindowProc?

	var lRet uintptr
	fCallDWP = DwmDefWindowProc(hWnd, msg, wParam, lParam, &lRet) == 0

	switch msg {
	case win.WM_CREATE:
		var rcClient win.RECT
		win.GetWindowRect(hWnd, &rcClient)
		// Inform application of the frame change.
		win.SetWindowPos(
			hWnd,
			0,
			rcClient.Left, rcClient.Top,
			rcClient.Right-rcClient.Left, rcClient.Bottom-rcClient.Top,
			win.SWP_FRAMECHANGED)
		fCallDWP = true
		lRet = 0
	case win.WM_PAINT:
		if window.quickAccessToolbar {
			var ps win.PAINTSTRUCT

			htheme := win.OpenThemeData(hWnd, syscall.StringToUTF16Ptr("BUTTON"))

			hdcPaint := win.BeginPaint(hWnd, &ps)
			//g := &Graphics{hdc: hdcPaint}

			rectBtn := win.RECT{Left: 0, Top: 0, Right: 50, Bottom: 24}

			win.DrawThemeBackground(htheme, hdcPaint, WP_MINBUTTON, MINBS_NORMAL, &rectBtn, nil)

			//rect := Rect{Left: 0, Top: 0, Right: 40, Bottom: 20}
			//g.FillRect(&rect, Rgb(255, 0, 0))

			win.EndPaint(hWnd, &ps)

			win.CloseThemeData(htheme)
			fCallDWP = true
			lRet = 0
		}
		case win.WM_ACTIVATE:
			var margins MARGINS
			margins.CxLeftWidth = 0
			margins.CxRightWidth = 0
			margins.CyBottomHeight = 0

			titleBarHeight := GetThemeSysSize(win.SM_CYSIZE) + GetThemeSysSize(SM_CXPADDEDBORDER)*2

			margins.CyTopHeight = 0
			DwmExtendFrameIntoClientArea(hWnd, &margins)
			fCallDWP = true
			lRet = 0
	case win.WM_NCCALCSIZE:
		if wParam == 1 {
			fCallDWP = false
			lRet = 0
		}
	}
	*pfCallDWP = fCallDWP
	return lRet
}
*/

// ReflectedMsg gets called from parent Window.
// We keep an empty body so that every struct that
// embeds Window won't have to also implement this method.
// Only who wants to (A button for example would want to
// override this and check 'WM_COMMAND' msg for click event
// handling.)
func (window *windowData) ReflectedMsg(reflectedFrom Window, msg uint32, wParam, lParam uintptr) {

}

func GET_WHEEL_DELTA_WPARAM(dw uintptr) int16 {
	return int16(dw >> 16 & 0xffff)
}

func appWindowProc(window *windowData, hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_COMMAND:
		hwndControl := win.HWND(lParam)
		windowControl := windowList[hwndControl]
		if windowControl != nil {
			windowControl.ReflectedMsg(window, msg, wParam, lParam)
		}
		return 0
	case win.WM_CLOSE:
		if window.CloseEvent != nil {
			if window.CloseEvent() {
				window.Dispose()
			}
		} else {
			window.Dispose()
		}
		return 0
	case win.WM_DESTROY:
		if window.DestroyEvent != nil {
			window.DestroyEvent()
		}
	case win.WM_WINDOWPOSCHANGED:
		if window.PositionChangedEvent != nil {
			wp := (*win.WINDOWPOS)(unsafe.Pointer(lParam))
			window.PositionChangedEvent(wp)
		}
	case win.WM_SIZE:
		window.RequestLayout()
	case win.WM_HSCROLL:
		if window.HScrollEvent != nil {
			window.HScrollEvent(int(win.LOWORD(uint32(wParam))), int(win.HIWORD(uint32(wParam))))
			return 0
		}
	case win.WM_VSCROLL:
		if window.VScrollEvent != nil {
			window.VScrollEvent(int(win.LOWORD(uint32(wParam))), int(win.HIWORD(uint32(wParam))))
			return 0
		}
	case win.WM_SETCURSOR:
		if window.SetCursorEvent != nil {
			if window.SetCursorEvent() {
				return 1
			}
			return 0
		}
	case win.WM_KILLFOCUS:
		if window.KillFocusEvent != nil {
			window.KillFocusEvent()
		}
	case win.WM_CHAR:
		if window.keyPressHandler != nil {
			window.keyPressHandler(int(wParam))
			return 0
		}
	case win.WM_MOUSEACTIVATE:
		if window.MouseActivateEvent != nil {
			return window.MouseActivateEvent()
		}
	case win.WM_MOUSEWHEEL:
		if window.MouseWheelEvent != nil {
			e := MouseWheelEvent{
				WheelDelta:    int(GET_WHEEL_DELTA_WPARAM(wParam)),
				MousePosition: Point{int(win.LOWORD(uint32(lParam))), int(win.HIWORD(uint32(lParam)))},
				VirtualKey:    int(win.LOWORD(uint32(wParam))),
			}
			window.MouseWheelEvent(&e)
			return 0
		}
	case win.WM_MOUSELEAVE:
		window.isMouseInside = false
		if window.MouseLeaveEvent != nil {
			window.MouseLeaveEvent()
		}
	case win.WM_MOUSEMOVE:
		if !window.isMouseInside {
			var mouseEvent win.TRACKMOUSEEVENT
			mouseEvent.CbSize = uint32(unsafe.Sizeof(mouseEvent))
			mouseEvent.DwFlags = win.TME_LEAVE
			mouseEvent.HwndTrack = hWnd
			mouseEvent.DwHoverTime = 0
			win.TrackMouseEvent(&mouseEvent)
			window.isMouseInside = true
		}
		if window.MouseMoveHandler != nil {
			var winpt win.POINT
			win.GetCursorPos(&winpt)
			win.ScreenToClient(hWnd, &winpt)
			pt := Point{X: int(winpt.X), Y: int(winpt.Y)}
			currentMouseButton := MouseButtonNone
			if (wParam & win.MK_LBUTTON) != 0 {
				currentMouseButton = MouseButtonLeft
			} else if (wParam & win.MK_RBUTTON) != 0 {
				currentMouseButton = MouseButtonRight
			}
			window.MouseMoveHandler(&pt, currentMouseButton)
		}
	case win.WM_LBUTTONDOWN:
		if window.MouseDownEvent != nil {
			var winpt win.POINT
			win.GetCursorPos(&winpt)
			win.ScreenToClient(hWnd, &winpt)
			pt := Point{X: int(winpt.X), Y: int(winpt.Y)}
			window.MouseDownEvent(&pt, MouseButtonLeft)
		}
	case win.WM_LBUTTONUP:
		if window.MouseUpEvent != nil {
			var winpt win.POINT
			win.GetCursorPos(&winpt)
			win.ScreenToClient(hWnd, &winpt)
			pt := Point{X: int(winpt.X), Y: int(winpt.Y)}
			window.MouseUpEvent(&pt, MouseButtonLeft)
		}
	case win.WM_RBUTTONDOWN:
		if window.MouseDownEvent != nil {
			var winpt win.POINT
			win.GetCursorPos(&winpt)
			win.ScreenToClient(hWnd, &winpt)
			pt := Point{X: int(winpt.X), Y: int(winpt.Y)}
			window.MouseDownEvent(&pt, MouseButtonRight)
		}
	case win.WM_RBUTTONUP:
		if window.MouseUpEvent != nil {
			var winpt win.POINT
			win.GetCursorPos(&winpt)
			win.ScreenToClient(hWnd, &winpt)
			pt := Point{X: int(winpt.X), Y: int(winpt.Y)}
			window.MouseUpEvent(&pt, MouseButtonRight)
		}
	case win.WM_PAINT:
		if window.orgWndProc == 0 {
			var ps win.PAINTSTRUCT
			hdcPaint := win.BeginPaint(hWnd, &ps)
			g := &Graphics{hdc: hdcPaint}
			if window.paintHandler != nil {
				win.SelectObject(g.hdc, win.HGDIOBJ(window.font))
				// Call paint event
				rect := window.GetClientRect()
				window.paintHandler(g, &rect)
			} else {
				color := FromCOLORREF(win.COLORREF(win.GetSysColor(win.COLOR_WINDOW)))
				rect := FromRECT(&ps.RcPaint)
				g.FillRect(&rect, &color)
			}
			win.EndPaint(hWnd, &ps)
		}
	case win.WM_ERASEBKGND:
		if window.orgWndProc == 0 {
			/*
				color := FromCOLORREF(win.COLORREF(win.GetSysColor(win.COLOR_WINDOW)))
				rect := window.GetClientRect()
				g := &Graphics{hdc: win.HDC(wParam)}
				g.FillRect(&rect, color)
			*/
			return 1
		}
	}
	if window.orgWndProc != 0 {
		return win.CallWindowProc(window.orgWndProc, hWnd, msg, wParam, lParam)
	}
	return win.DefWindowProc(hWnd, msg, wParam, lParam)
}

func registerWindow(classname *uint16, wndproc func(win.HWND, uint32, uintptr, uintptr) uintptr) bool {
	logInfo("'registerWindow'")
	var wcex win.WNDCLASSEX
	wcex.CbSize = uint32(unsafe.Sizeof(wcex))
	wcex.Style = 0 //win.CS_HREDRAW | win.CS_VREDRAW
	wcex.LpfnWndProc = syscall.NewCallback(wndproc)
	wcex.CbClsExtra = 0
	wcex.CbWndExtra = 0
	wcex.HInstance = app.hInstance
	wcex.HCursor = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	wcex.HbrBackground = win.HBRUSH(win.COLOR_WINDOW)
	wcex.LpszClassName = classname
	wcex.HIcon = win.LoadIcon(app.hInstance, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wcex.HIconSm = win.LoadIcon(app.hInstance, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wcex.LpszMenuName = nil
	return win.RegisterClassEx(&wcex) != 0
}

func (window *windowData) GetText() string {
	hwnd := window.handle
	textLength := win.SendMessage(hwnd, win.WM_GETTEXTLENGTH, 0, 0)
	buf := make([]uint16, textLength+1)
	win.SendMessage(hwnd, win.WM_GETTEXT, uintptr(textLength+1), uintptr(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

func (window *windowData) RequestLayout() {
	logInfo("RequestLayout...")
	rectClient := window.GetClientRect()
	if window.quickAccessToolbar {
		rectClient.Top += 30
	}
	controls := window.GetChildrens()
	hdwp := win.BeginDeferWindowPos(int32(len(controls)))
	for _, control := range controls {
		// Dock if needed
		size := control.GetSize()
		marginLeft, marginRight, marginTop, marginBottom := control.GetMargin()
		switch control.GetDockType() {
		case DockFill:
			// dock fill to available space
			//logInfo("Dock fill to available space")
			left := rectClient.Left + marginLeft
			right := rectClient.Right - marginRight
			top := rectClient.Top + marginTop
			bottom := rectClient.Bottom - marginBottom
			width := right - left
			height := bottom - top
			//control.MoveWindow(left, top, width, height, true)
			win.DeferWindowPos(hdwp, control.GetHandle(), 0,
				int32(left), int32(top), int32(width), int32(height),
				win.SWP_SHOWWINDOW|win.SWP_FRAMECHANGED)
		case DockLeft:
			// dock to bottom
			//logInfo("Dock to bottom")
			rectClient.Left += marginLeft

			rectClient.Top += marginTop
			rectClient.Bottom -= marginBottom
			availableHeight := rectClient.Height()
			//control.MoveWindow(rectClient.Left, rectClient.Top, size.Width, availableHeight, true)
			win.DeferWindowPos(hdwp, control.GetHandle(), 0,
				int32(rectClient.Left), int32(rectClient.Top), int32(size.Width), int32(availableHeight),
				win.SWP_SHOWWINDOW|win.SWP_FRAMECHANGED)

			rectClient.Top -= marginTop
			rectClient.Bottom += marginBottom

			rectClient.Left += size.Width + marginRight
		case DockTop:
			// dock to top
			//logInfo("Dock to top")
			rectClient.Left += marginLeft
			rectClient.Top += marginTop
			rectClient.Right -= marginRight
			availableWidth := rectClient.Width()
			//control.MoveWindow(rectClient.Left, rectClient.Top, availableWidth, size.Height, true)
			win.DeferWindowPos(hdwp, control.GetHandle(), 0,
				int32(rectClient.Left), int32(rectClient.Top), int32(availableWidth), int32(size.Height),
				win.SWP_SHOWWINDOW|win.SWP_FRAMECHANGED)

			rectClient.Left -= marginLeft
			rectClient.Right += marginRight

			rectClient.Top += size.Height + marginBottom
		case DockBottom:
			// dock to bottom
			//logInfo("Dock to bottom")
			rectClient.Left += marginLeft
			rectClient.Right -= marginRight

			rectClient.Bottom -= marginBottom
			availableWidth := rectClient.Width()
			//control.MoveWindow(rectClient.Left, rectClient.Bottom-size.Height, availableWidth, size.Height, true)
			win.DeferWindowPos(hdwp, control.GetHandle(), 0,
				int32(rectClient.Left), int32(rectClient.Bottom-size.Height), int32(availableWidth), int32(size.Height),
				win.SWP_SHOWWINDOW|win.SWP_FRAMECHANGED)

			rectClient.Left -= marginLeft
			rectClient.Right += marginRight

			rectClient.Bottom -= size.Height + marginTop
		default:
			// do nothing
		}
	}
	win.EndDeferWindowPos(hdwp)
	if window.ResizeEvent != nil {
		// send the updated 'rectClient'
		window.ResizeEvent(&rectClient)
	}
	window.Repaint()
}

// Create creates the window
func (window *windowData) Create(text string, style uint, x, y, width, height int, parent Window) {
	window.createWindowImpl(text, style, 0, x, y, width, height, parent)
}

// CreateEx creates the window with the given extended styles
func (window *windowData) CreateEx(text string, style, styleEx uint, x, y, width, height int, parent Window) {
	window.createWindowImpl(text, style, styleEx, x, y, width, height, parent)
}

// Event Handlers
func (window *windowData) SetPaintEventHandler(f func(g *Graphics, rc *Rect)) {
	window.paintHandler = f
}
func (window *windowData) SetMouseMoveEventHandler(f func(pt *Point, mbutton int)) {
	window.MouseMoveHandler = f
}
func (window *windowData) SetMouseDownEventHandler(f func(pt *Point, mbutton int)) {
	window.MouseDownEvent = f
}
func (window *windowData) SetMouseUpEventHandler(f func(pt *Point, mbutton int)) {
	window.MouseUpEvent = f
}
func (window *windowData) SetMouseLeaveEventHandler(f func()) {
	window.MouseLeaveEvent = f
}
func (window *windowData) SetMouseWheelEventHandler(f func(e *MouseWheelEvent)) {
	window.MouseWheelEvent = f
}
func (window *windowData) SetKeyDownEventHandler(f func(keycode int)) {
	window.KeyDownEvent = f
}
func (window *windowData) SetKeyUpEventHandler(f func(keycode int)) {
	window.KeyUpEvent = f
}
func (window *windowData) SetKeyPressEventHandler(f func(keycode int)) {
	window.keyPressHandler = f
}
func (window *windowData) SetResizeEventHandler(f func(client *Rect)) {
	window.ResizeEvent = f
}
func (window *windowData) SetKillFocusEventHandler(f func()) {
	window.KillFocusEvent = f
}
func (window *windowData) SetDestroyEventHandler(f func()) {
	window.DestroyEvent = f
}
func (window *windowData) SetCloseEventHandler(f func() bool) {
	window.CloseEvent = f
}
func (window *windowData) SetPositionChangedEventHandler(f func(wp *win.WINDOWPOS)) {
	window.PositionChangedEvent = f
}
func (window *windowData) SetMouseActivateEventHandler(f func() uintptr) {
	window.MouseActivateEvent = f
}
func (window *windowData) SetSetCursorEventHandler(f func() bool) {
	window.SetCursorEvent = f
}
func (window *windowData) SetHScrollEventHandler(f func(stype, position int)) {
	window.HScrollEvent = f
}
func (window *windowData) SetVScrollEventHandler(f func(stype, position int)) {
	window.VScrollEvent = f
}

func (window *windowData) SetFont(font win.HFONT) {
	window.font = font
	window.Repaint()
}

func (window *windowData) GetFont() win.HFONT {
	return window.font
}

func (window *windowData) HasParent() bool {
	return window.parent != nil
}

func (window *windowData) GetParent() Window {
	return window.parent
}

func (window *windowData) addChild(child Window) {
	window.childrens = append(window.childrens, child)
}

func (window *windowData) GetChildrens() []Window {
	return window.childrens
}

func (window *windowData) SetMargin(left, right, top, bottom int) {
	window.marginLeft = left
	window.marginRight = right
	window.marginTop = top
	window.marginBottom = bottom
	if window.HasParent() {
		window.GetParent().RequestLayout()
	}
}

func (window *windowData) SetDockType(dtype int) {
	window.dockType = dtype
	if window.HasParent() {
		window.GetParent().RequestLayout()
	}
}

func (window *windowData) GetDockType() int {
	return window.dockType
}

func (window *windowData) GetMargin() (left, right, top, bottom int) {
	return window.marginLeft, window.marginRight, window.marginTop, window.marginBottom
}

// GetHandle returns handle/id of the window
func (window *windowData) GetHandle() win.HWND {
	return window.handle
}

func (window *windowData) GetDC() win.HDC {
	return win.GetDC(window.handle)
}

func (window *windowData) ReleaseDC(hdc win.HDC) {
	win.ReleaseDC(window.handle, hdc)
}

// MoveWindow sets the x, y and width, height of current window
func (window *windowData) MoveWindow(x, y, width, height int, repaint bool) {
	win.MoveWindow(window.handle, int32(x), int32(y), int32(width), int32(height), repaint)
}

// GetSize returns...size?..ugh
func (window *windowData) GetSize() Size {
	var rc win.RECT
	win.GetClientRect(window.handle, &rc)
	return Size{Width: int(rc.Right - rc.Left), Height: int(rc.Bottom - rc.Top)}
}

// GetClientRect returns...the client rect?..ugh
func (window *windowData) GetClientRect() Rect {
	var rc win.RECT
	win.GetClientRect(window.handle, &rc)
	return Rect{Left: int(rc.Left), Top: int(rc.Top), Right: int(rc.Right), Bottom: int(rc.Bottom)}
}

// GetWindowRect returns screen space rect
func (window *windowData) GetWindowRect() Rect {
	var rc win.RECT
	win.GetWindowRect(window.handle, &rc)
	return Rect{Left: int(rc.Left), Top: int(rc.Top), Right: int(rc.Right), Bottom: int(rc.Bottom)}
}

// SetPosition sets the position
func (window *windowData) SetPosition(x, y int) bool {
	return win.SetWindowPos(window.handle, 0, int32(x), int32(y), 0, 0, win.SWP_NOZORDER|win.SWP_NOSIZE)
}

// SetSize sets the size
func (window *windowData) SetSize(width, height int) bool {
	ret := win.SetWindowPos(window.handle, 0, 0, 0, int32(width), int32(height), win.SWP_NOZORDER|win.SWP_NOMOVE)
	return ret
}

// Repaint re-paints the window
func (window *windowData) Repaint() {
	win.InvalidateRect(window.handle, nil, true)
}

func (window *windowData) InvalidateRect(rect *Rect, eraseBackground bool) {
	if rect == nil {
		win.InvalidateRect(window.handle, nil, eraseBackground)
	} else {
		wrect := rect.AsRECT()
		win.InvalidateRect(window.handle, &wrect, eraseBackground)
	}
}

func (window *windowData) Update() {
	win.UpdateWindow(window.handle)
}

// SetVisible sets
func (window *windowData) SetVisible(visible bool) {
	if visible {
		win.ShowWindow(window.handle, win.SW_SHOW)
	} else {
		win.ShowWindow(window.handle, win.SW_HIDE)
	}
}

// SetText sets the caption of the window
func (window *windowData) SetText(text string) {
	textUTF16, _ := syscall.UTF16PtrFromString(text)
	win.SendMessage(window.handle, win.WM_SETTEXT, 0, uintptr(unsafe.Pointer(textUTF16)))
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func MAKEWPARAM(low, high uint16) uintptr {
	return uintptr(low) | uintptr(high)<<16
}

func MAKELPARAM(low, high uint16) uintptr {
	return MAKEWPARAM(low, high)
}
