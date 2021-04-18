package reza

import (
	win "github.com/lxn/win"
)

// PopupMenu is a popup window to show menu
type PopupMenu interface {
	// embed the popup window type
	PopupWindow
	SetLargeItem(large bool)
}

type popupMenuData struct {
	// embed the popup window type
	PopupWindow
	isBigIcon  bool
	itemHeight int
	leftMargin int
}

// PopupMenuItem is an item of a popup/dropdown menu
type PopupMenuItem interface {
	PopupItem
	HasText() bool
	GetText() string
	SetText(text string)
	IsEnabled() bool
	SetEnabled(enabled bool)
	SetToggled(toogle bool)
	IsToggled() bool
	SetRectText(rc Rect)
	GetRectText() Rect
	GetIcon() *BitmapImage
}

type popupMenuItemData struct {
	// Embed
	popupItemData
	text    string
	toggled bool
	enabled bool
	icon    *BitmapImage
	rcText  Rect
}

// MenuItemInfo is for simple declaration
type MenuItemInfo struct {
	Text     string
	IconPath string
	Sperator bool
	OnClick  func(e *PopupItemEvent)
	AssignTo *PopupMenuItem
}

var menuBorderColor = Rgb(210, 210, 210)
var capSepHeight = 22

func NewPopupMenu(parent Window, items []MenuItemInfo) PopupMenu {
	menu := &popupMenuData{PopupWindow: NewPopupWindow()}
	menu.initWithItems(parent, items)
	return menu
}

// HasText returns true if has text
func (menuitem *popupMenuItemData) HasText() bool {
	return len(menuitem.text) > 0
}

// GetText returns the text
func (menuitem *popupMenuItemData) GetText() string {
	return menuitem.text
}

// SetText sets the specified string as the item caption
func (menuitem *popupMenuItemData) SetText(text string) {
	menuitem.text = text
}

func (menuitem *popupMenuItemData) IsEnabled() bool {
	return menuitem.enabled
}

func (menuitem *popupMenuItemData) SetEnabled(enabled bool) {
	menuitem.enabled = enabled
}

// SetRectText does something
func (menuitem *popupMenuItemData) SetRectText(rc Rect) {
	menuitem.rcText = rc
}

// GetRectText also does something
func (menuitem *popupMenuItemData) GetRectText() Rect {
	return menuitem.rcText
}

// GetIcon returns the item icon, if has any
func (menuitem *popupMenuItemData) GetIcon() *BitmapImage {
	return menuitem.icon
}

func (item *popupMenuItemData) SetToggled(toogle bool) {
	item.toggled = toogle
}

func (item *popupMenuItemData) IsToggled() bool {
	return item.toggled
}

// init initializes the menu
func (menu *popupMenuData) initWithItems(parent Window, items []MenuItemInfo) PopupMenu {
	menu.PopupWindow.Init(parent)
	menu.SetPaintEventHandler(menu.paint)
	menu.SetMeasureContentSize(menu.measureContentSize)
	/*
		if menu.isBigIcon {
			menu.itemHeight = 48
			menu.leftMargin = 52
		} else {*/
	menu.itemHeight = 24
	menu.leftMargin = 26
	//}
	if items != nil {
		menu.AddItems(items)
	}
	return menu
}

func (menu *popupMenuData) SetLargeItem(large bool) {
	menu.isBigIcon = large
	if large {
		menu.itemHeight = 48
		menu.leftMargin = 52
	} else {
		menu.itemHeight = 24
		menu.leftMargin = 26
	}
}

// AddItems adds item array to the menu
func (menu *popupMenuData) AddItems(newItems []MenuItemInfo) PopupMenu {
	for _, newItem := range newItems {
		menuitem := &popupMenuItemData{}
		menuitem.text = newItem.Text
		menuitem.isSperator = newItem.Sperator
		menuitem.clickEvent = newItem.OnClick
		menuitem.enabled = true
		if newItem.AssignTo != nil {
			*newItem.AssignTo = menuitem
		}
		// Check if we were given a icon path
		if len(newItem.IconPath) > 0 {
			menuitem.icon, _ = CreateBitmapImage(newItem.IconPath, true)
		}
		menu.AddItem(menuitem)
	}
	menu.Repaint()
	return menu
}

func (menu *popupMenuData) measureContentSize(g *Graphics) (width int, height int) {
	contentWidth := 1
	contentHeight := 1
	// measure width
	font := menu.GetFont()
	items := menu.GetItems()
	for i := range items {
		item := items[i].(PopupMenuItem)
		if item.HasText() {
			rect := g.MeasureText(item.GetText(), win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, font)
			if rect.Width() > contentWidth {
				contentWidth = rect.Width()
			}
		}
	}
	contentWidth += menu.leftMargin + 40 // some extra width for a nice look..idk :)
	// measure height
	for i := range items {
		item := items[i].(PopupMenuItem)
		if item.IsSperator() {
			if item.HasText() {
				contentHeight += capSepHeight
			}
		} else {
			contentHeight += menu.itemHeight
		}
	}
	return contentWidth, contentHeight
}

func (menu *popupMenuData) measureItems(rect *Rect) {
	itemTop := rect.Top

	// Update items data
	items := menu.GetItems()
	for i := range items {
		item := items[i].(PopupMenuItem)
		prect := item.GetRect() // returns pointer
		if item.IsSperator() {
			if item.HasText() {
				prect.Left = rect.Left
				prect.Top = itemTop
				prect.Right = rect.Right
				prect.Bottom = prect.Top + capSepHeight
				itemTop += capSepHeight
			} else {
				prect.Left = rect.Left + menu.leftMargin + 1
				prect.Top = itemTop
				prect.Right = rect.Right
				prect.Bottom = prect.Top
			}
		} else {
			prect.Left = rect.Left
			prect.Top = itemTop
			prect.Right = rect.Right
			prect.Bottom = prect.Top + menu.itemHeight

			rcText := *prect
			rcText.Left = prect.Left + menu.leftMargin + 8

			item.SetRectText(rcText)

			itemTop += menu.itemHeight
		}
	}
}

func (menu *popupMenuData) drawItems(g *Graphics, rect *Rect) {
	font := menu.GetFont()
	items := menu.GetItems()
	for i := range items {
		item := items[i].(*popupMenuItemData)
		prect := item.GetRect() // returns pointer
		if item.IsSperator() {
			// Draw a captioned seperator if text isn't empty
			if item.HasText() {
				rcSeperator := *item.GetRect()
				g.FillRectangle(&rcSeperator, NewRgb(255, 255, 255), NewRgb(255, 255, 255))
				g.DrawLine(rect.Left, rect.Top, rect.Left, rect.Bottom, &menuBorderColor)
				g.DrawLine(rect.Right-1, rect.Top, rect.Right-1, rect.Bottom, &menuBorderColor)
				g.DrawLine(rect.Left+2, rcSeperator.Top, rect.Right-3, rcSeperator.Top, &menuBorderColor)
				g.DrawLine(rect.Left+2, rcSeperator.Bottom, rect.Right-3, rcSeperator.Bottom, &menuBorderColor)

				rcSeperator.Left += 8

				g.DrawText(item.GetText(), &rcSeperator, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, NewRgb(80, 80, 80), font)
			} else {
				g.DrawLine(prect.Left+2, prect.Top, prect.Right-3, prect.Top, &menuBorderColor)
			}

		} else {
			if item.IsHighlighted() {
				rcHover := *prect // copy
				rcHover.Left += 2
				rcHover.Right -= 2
				rcHover.Top += 2
				rcHover.Bottom--
				g.FillRectangle(&rcHover, NewRgb(168, 210, 253), NewRgb(237, 244, 252))
			}
			rectIcon := *prect
			rectIcon.Right = prect.Left + menu.leftMargin
			if item.toggled {
				rectIcon.Left += 2
				rectIcon.Right -= 2
				rectIcon.Top += 2
				rectIcon.Bottom -= 2
				g.FillRectangle(&rectIcon, NewRgb(100, 165, 230), NewRgb(205, 230, 252))
			}
			rectIcon = *prect
			rectIcon.Right = prect.Left + menu.leftMargin
			if item.GetIcon() != nil {
				icon := item.GetIcon()
				g.DrawBitmapImageCenter(icon, &rectIcon, !item.IsEnabled())
			}
			// Draw a menu item
			rcText := item.GetRectText()
			if item.IsEnabled() {
				g.DrawText(item.GetText(), &rcText, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, NewRgb(60, 60, 60), font)
			} else {
				g.DrawText(item.GetText(), &rcText, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, NewRgb(140, 140, 140), font)
			}
		}
	}
}

func (menu *popupMenuData) paint(gOrg *Graphics, rect *Rect) {
	var db DoubleBuffer
	g := db.BeginDoubleBuffer(gOrg, rect, &menuBorderColor, NewRgb(255, 255, 255))
	// Draw left margin
	g.DrawLine(rect.Left+menu.leftMargin, rect.Top+4, rect.Left+menu.leftMargin, rect.Bottom-4, &menuBorderColor)
	// Measure items rects
	menu.measureItems(rect)
	// Draw items
	menu.drawItems(g, rect)
	db.EndDoubleBuffer()
}
