package reza

// PopupSizeMenu is a popup menu to show list of sizes to choose from
type PopupSizeMenu interface {
	// Inherit data from the window type
	PopupWindow
}

type popupSizeMenuData struct {
	// Inherit data from the window type
	PopupWindow
}

// PopupSizeMenuItem is a size menu item
type PopupSizeMenuItem interface {
	PopupItem
	SetSize(size int)
	GetSize() int
	IsToggled() bool
	SetToggled(value bool)
}

type popupSizeMenuItemData struct {
	// Embed
	popupItemData
	// public fields
	size    int
	toggled bool
}

type SizeMenuItemInfo struct {
	Size     int // size in pixel
	Toggled  bool
	OnClick  func(e *PopupItemEvent)
	AssignTo *PopupSizeMenuItem
}

const sizeMenuItemHeight = 40

func NewPopupSizeMenu(parent Window, items []SizeMenuItemInfo) PopupSizeMenu {
	menu := &popupSizeMenuData{PopupWindow: NewPopupWindow()}
	menu.initWithItems(parent, items)
	return menu
}

func (item *popupSizeMenuItemData) SetSize(size int) {
	item.size = size
}

// GetSize returns size
func (item *popupSizeMenuItemData) GetSize() int {
	return item.size
}

// IsToggled returns true if item is toggled
func (item *popupSizeMenuItemData) IsToggled() bool {
	return item.toggled
}

// SetToggled sets whether the item will be toggled or not
func (item *popupSizeMenuItemData) SetToggled(value bool) {
	item.toggled = value
}

// Init initializes the menu
func (menu *popupSizeMenuData) initWithItems(parent Window, items []SizeMenuItemInfo) PopupSizeMenu {
	menu.PopupWindow.Init(parent)
	menu.SetPaintEventHandler(menu.paint)
	menu.SetMeasureContentSize(menu.measureContentSize)
	if items != nil {
		menu.AddItems(items)
	}
	return menu
}

// AddItems adds item array to the menu
func (menu *popupSizeMenuData) AddItems(items []SizeMenuItemInfo) PopupSizeMenu {
	for _, item := range items {
		newSizeItem := &popupSizeMenuItemData{}
		newSizeItem.SetClickEvent(item.OnClick)
		newSizeItem.SetToggled(item.Toggled)
		newSizeItem.SetSize(item.Size)
		if item.AssignTo != nil {
			*item.AssignTo = newSizeItem
		}
		menu.AddItem(newSizeItem)
	}
	menu.Repaint()
	return menu
}

func (menu *popupSizeMenuData) measureContentSize(g *Graphics) (width int, height int) {
	contentWidth := 130
	contentHeight := 0
	// measure height
	items := menu.GetItems()
	for i := range items {
		item := items[i].(PopupSizeMenuItem)
		item.GetRect()
		contentHeight += sizeMenuItemHeight
	}
	return contentWidth, contentHeight
}

func (menu *popupSizeMenuData) measureItems(rect *Rect) {
	itemTop := rect.Top
	// Update items data
	// NOTE: 'PopupWindow' is always initialized with '&popupWindowData{}'...and always must be!
	items := menu.PopupWindow.(*popupWindowData).items
	for i := range items {
		item := items[i].(PopupSizeMenuItem)
		prect := item.GetRect() // returns pointer
		prect.Left = rect.Left
		prect.Top = itemTop
		prect.Right = rect.Right
		prect.Bottom = prect.Top + sizeMenuItemHeight

		itemTop += sizeMenuItemHeight
	}
}

func (menu *popupSizeMenuData) drawItems(g *Graphics, rect *Rect) {
	items := menu.GetItems()
	for i := range items {
		item := items[i].(PopupSizeMenuItem)
		prect := item.GetRect() // returns pointer
		size := item.GetSize()
		rcHover := *prect // copy
		rcHover.Left += 2
		rcHover.Right -= 2
		rcHover.Top += 2
		rcHover.Bottom--
		if item.IsHighlighted() {
			if item.IsToggled() {
				g.FillRectangle(&rcHover, NewRgb(125, 179, 234), NewRgb(219, 235, 252))
			} else {
				g.FillRectangle(&rcHover, NewRgb(168, 210, 253), NewRgb(237, 244, 252))
			}
		} else if item.IsToggled() {
			g.FillRectangle(&rcHover, NewRgb(100, 165, 230), NewRgb(206, 229, 252))
		}
		centerY := prect.CenterY()
		rcLine := *prect // copy
		rcLine.Left += 6
		rcLine.Right -= 6
		rcLine.Top = centerY
		rcLine.Bottom = rcLine.Top + size
		if size > 1 {
			// lets keep it perfectly in the center
			halfSize := size / 2
			rcLine.Top -= halfSize
			rcLine.Bottom -= halfSize
		}
		g.FillRectangle(&rcLine, NewRgb(0, 0, 0), NewRgb(0, 0, 0))
	}
}

func (menu *popupSizeMenuData) paint(gOrg *Graphics, rect *Rect) {
	var db DoubleBuffer
	g := db.BeginDoubleBuffer(gOrg, rect, &menuBorderColor, NewRgb(255, 255, 255))
	defer db.EndDoubleBuffer()
	// Measure items rects
	menu.measureItems(rect)
	// Draw items
	menu.drawItems(g, rect)
}
