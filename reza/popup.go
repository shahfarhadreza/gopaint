package reza

import (
	"log"

	win "github.com/lxn/win"
)

// PopupWindow is window that pops up out of nowhere ;)
type PopupWindow interface {
	Window
	Init(parent Window)
	Popup(x, y int)
	AddItem(item PopupItem)
	GetItems() []PopupItem
	SetMeasureContentSize(f func(g *Graphics) (width int, height int))
}

type popupWindowData struct {
	// Inherit data from the window type
	Window
	// private fields
	done        bool      // when done (user clicks elsewhere) window gets hidden
	itemClicked PopupItem // item to sent click event after done
	items       []PopupItem
	// public fields
	MeasureContentSize func(g *Graphics) (width int, height int)
}

// PopupItem is a basic interface to any popup item
type PopupItem interface {
	GetRect() *Rect
	SetHighlight(value bool)
	IsHighlighted() bool
	HasClickEvent() bool
	SetClickEvent(func(e *PopupItemEvent))
	IsSperator() bool
	SetSperator(value bool)
	asPopupItemData() *popupItemData
}

// popupItemData is item belongs to a popup window
type popupItemData struct {
	rect       Rect
	highlight  bool
	isSperator bool
	clickEvent func(e *PopupItemEvent)
}

// PopupItemEvent is an event
type PopupItemEvent struct {
	Item        PopupItem
	PopupWindow PopupWindow
}

func NewPopupWindow() PopupWindow {
	return &popupWindowData{Window: NewWindow()}
}

func (item *popupItemData) asPopupItemData() *popupItemData {
	return item
}

// GetRect returns rect
func (item *popupItemData) GetRect() *Rect {
	return &item.rect
}

// IsHighlighted returns true if highlighted
func (item *popupItemData) IsHighlighted() bool {
	return item.highlight
}

// SetHighlight sets
func (item *popupItemData) SetHighlight(value bool) {
	item.highlight = value
}

// SetClickEvent sets
func (item *popupItemData) SetClickEvent(e func(e *PopupItemEvent)) {
	item.clickEvent = e
}

// HasClickEvent returns true if has event
func (item *popupItemData) HasClickEvent() bool {
	return item.clickEvent != nil
}

// IsSperator returns true if item is a sperator
func (item *popupItemData) IsSperator() bool {
	return item.isSperator
}

// SetSperator sets
func (item *popupItemData) SetSperator(value bool) {
	item.isSperator = value
}

// GetItems returns the item list
func (popup *popupWindowData) GetItems() []PopupItem {
	return popup.items
}

func (popup *popupWindowData) AddItem(item PopupItem) {
	popup.items = append(popup.items, item)
}

func (popup *popupWindowData) SetMeasureContentSize(f func(g *Graphics) (width int, height int)) {
	popup.MeasureContentSize = f
}

// Init initializes
func (popup *popupWindowData) Init(parent Window) {
	popup.CreateEx("DropDownMenu", win.WS_POPUP, win.WS_EX_NOACTIVATE|win.WS_EX_TOPMOST, 10, 10, 10, 10, parent)
	popup.SetMouseDownEventHandler(popup.mouseDown)
	popup.SetMouseMoveEventHandler(popup.mouseMove)
	popup.SetMouseUpEventHandler(popup.mouseUp)
	popup.SetMouseActivateEventHandler(func() uintptr {
		return 3 // MA_NOACTIVATE
	})
	popup.items = make([]PopupItem, 0)
}

func (popup *popupWindowData) mouseDown(ptClient *Point, mbutton int) {
	rect := popup.GetWindowRect()
	pt := app.GetCursorPos()
	if !rect.IsPointInside(&pt) {
		popup.done = true
	}
}

func (popup *popupWindowData) mouseMove(pt *Point, mbutton int) {
	for i := range popup.items {
		popup.items[i].SetHighlight(false)
	}
	for i := range popup.items {
		if pt.IsInsideRect(popup.items[i].GetRect()) {
			popup.items[i].SetHighlight(true)
			break
		}
	}
	popup.Repaint()
}

func (popup *popupWindowData) mouseUp(pt *Point, mbutton int) {
	if mbutton == MouseButtonLeft {
		for _, item := range popup.items {
			if !item.IsSperator() && pt.IsInsideRect(item.GetRect()) {
				log.Printf("Clicked\n")
				popup.done = true
				popup.itemClicked = item
				break
			}
		}
	}
}

func (popup *popupWindowData) getContentSize() (width, height int) {
	hdc := win.GetDC(popup.GetHandle())
	g := &Graphics{hdc}
	defFont := popup.GetFont()
	win.SelectObject(hdc, win.HGDIOBJ(defFont))
	cwidth, cheight := popup.MeasureContentSize(g)
	win.ReleaseDC(popup.GetHandle(), hdc)
	return cwidth, cheight
}

// Popup does NOTHING..like...NOTHING
func (popup *popupWindowData) Popup(x, y int) {
	if popup.MeasureContentSize == nil {
		logError("'MeasureContentSize' is 'nil' can't popup!!!")
		return
	}
	// First, calculate how much width and height do we need
	cwidth, cheight := popup.getContentSize()

	totalHeight := cheight + 2
	totalWidth := cwidth + 2

	hwndPopup := popup.GetHandle()
	hwndOwner := popup.GetParent().GetHandle()

	// Now set parameters accordingly and show the menu
	win.SetWindowPos(hwndPopup, win.HWND_TOP, int32(x), int32(y), int32(totalWidth), int32(totalHeight), win.SWP_NOZORDER|win.SWP_NOACTIVATE)
	win.ShowWindow(hwndPopup, win.SW_SHOWNOACTIVATE)

	win.SetCapture(hwndOwner)

	popup.itemClicked = nil
	popup.done = false

	var msg win.MSG
	for {
		if win.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		if popup.done {
			break
		}
		/*
		 *  If our owner stopped being the active window
		 *  (e.g., the user Alt+Tab'd to another window
		 *  in the meantime), then stop.
		 */
		hwndActive := win.GetActiveWindow()
		if hwndActive != hwndOwner && !win.IsChild(hwndActive, hwndOwner) {
			break
		}
		if GetCapture() != hwndOwner {
			break
		}
		/*
		 *  At this point, we get to snoop at all input messages
		 *  before they get dispatched.  This allows us to
		 *  route all input to our popup window even if really
		 *  belongs to somebody else.
		 *
		 *  All mouse messages are remunged and directed at our
		 *  popup menu.  If the mouse message arrives as client
		 *  coordinates, then we have to convert it from the
		 *  client coordinates of the original target to the
		 *  client coordinates of the new target.
		 */
		switch msg.Message {
		/*
		*  These mouse messages arrive in client coordinates,
		*  so in addition to stealing the message, we also
		*  need to convert the coordinates.
		 */
		case win.WM_MOUSEMOVE,
			win.WM_LBUTTONDOWN,
			win.WM_LBUTTONUP,
			win.WM_LBUTTONDBLCLK,
			win.WM_RBUTTONDOWN,
			win.WM_RBUTTONUP,
			win.WM_RBUTTONDBLCLK,
			win.WM_MBUTTONDOWN,
			win.WM_MBUTTONUP,
			win.WM_MBUTTONDBLCLK:
			var pt win.POINT
			pt.X = int32(win.LOWORD(uint32(msg.LParam)))
			pt.Y = int32(win.HIWORD(uint32(msg.LParam)))
			MapWindowPoints(msg.HWnd, hwndPopup, &pt, 1)
			msg.LParam = MAKELPARAM(uint16(pt.X), uint16(pt.Y))
			msg.HWnd = hwndPopup
		/*
		 *  These mouse messages arrive in screen coordinates,
		 *  so we just need to steal the message.
		 */
		case win.WM_NCMOUSEMOVE,
			win.WM_NCLBUTTONDOWN,
			win.WM_NCLBUTTONUP,
			win.WM_NCLBUTTONDBLCLK,
			win.WM_NCRBUTTONDOWN,
			win.WM_NCRBUTTONUP,
			win.WM_NCRBUTTONDBLCLK,
			win.WM_NCMBUTTONDOWN,
			win.WM_NCMBUTTONUP,
			win.WM_NCMBUTTONDBLCLK:
			msg.HWnd = hwndPopup
		/*
		 *  Steal all keyboard messages, too.
		 */
		case win.WM_KEYDOWN:
		case win.WM_KEYUP:
		case win.WM_CHAR:
		case win.WM_DEADCHAR:
		case win.WM_SYSKEYDOWN:
		case win.WM_SYSKEYUP:
		case win.WM_SYSCHAR:
		case win.WM_SYSDEADCHAR:
			msg.HWnd = hwndPopup
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
		if popup.done {
			break
		}
		/*
		 *  If our owner stopped being the active window
		 *  (e.g., the user Alt+Tab'd to another window
		 *  in the meantime), then stop.
		 */
		hwndActive = win.GetActiveWindow()
		if hwndActive != hwndOwner && !win.IsChild(hwndActive, hwndOwner) {
			break
		}
		if GetCapture() != hwndOwner {
			break
		}
	}
	win.ShowWindow(popup.GetHandle(), win.SW_HIDE)
	win.ReleaseCapture()
	if popup.itemClicked != nil {
		data := popup.itemClicked.asPopupItemData()
		if data.clickEvent != nil {
			e := PopupItemEvent{Item: popup.itemClicked, PopupWindow: popup}
			data.clickEvent(&e)
		}
	}
	/*
	 *  If we got a WM_QUIT message, then re-post it so the caller's
	 *  message loop will see it.
	 */
	if msg.Message == win.WM_QUIT {
		win.PostQuitMessage(int32(msg.WParam))
	}
}
