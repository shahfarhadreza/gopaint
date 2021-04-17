package reza

import (
	"log"
	"syscall"

	win "github.com/lxn/win"
)

// Ribbon is the top ribbon which contains all the tools
type Ribbon interface {
	// Promote Window interface methods
	Window
	AddTab(text string) *RibbonTab
	SetCurrentTab(tab *RibbonTab)
	SetApplicationMenu(caption string, menu PopupMenu)
	IsRepaintSuspended() bool
	SuspendRepaint()
	ResumeRepaint()
}

type ribbonData struct {
	// Embedd Window interface
	Window
	tabs           []*RibbonTab
	tabSelected    *RibbonTab
	initialized    bool
	rectAppMenu    Rect
	appMenu        PopupMenu
	appMenuState   int // 0 - normal, 1 - highlight, 2 - pressed
	appMenuCaption string
	suspendRepaint bool
}

const tabHeaderHeight = 22
const tabHeaderWidth = 56

var imageCheck *BitmapImage = nil

func NewRibbon(parent Window) Ribbon {
	ribbon := &ribbonData{Window: NewWindow()}
	ribbon.init(parent)
	return ribbon
}

// Init initializes the ribbon
func (ribbon *ribbonData) init(parent Window) {
	logInfo("init ribbon")

	ribbon.initialized = false

	ribbon.tabs = make([]*RibbonTab, 0)
	ribbon.Create("", win.WS_CHILD|win.WS_CLIPCHILDREN|win.WS_VISIBLE, 10, 10, 10, 10, parent)
	ribbon.SetPaintEventHandler(ribbon.paint)
	ribbon.SetMouseMoveEventHandler(ribbon.mouseMove)
	ribbon.SetMouseDownEventHandler(ribbon.mouseDown)
	ribbon.SetMouseUpEventHandler(ribbon.mouseUp)

	// We just need only one image instence. We can draw it over and over whenever and whereever we need
	if imageCheck == nil {
		imageCheck, _ = CreateBitmapImage(".\\icons\\check.png", false)
	}
	ribbon.initialized = true
}

// Dispose cleans up stuffs from the memory
func (ribbon *ribbonData) Dispose() {
	logInfo("Disposing ribbon...")
	for _, tab := range ribbon.tabs {
		for _, sec := range tab.sections {
			for _, ibutton := range sec.buttons {
				ibutton.Dispose()
			}
		}
	}
	ribbon.Window.Dispose()
}

// Prevents it to repaint itself while setting various elements properties
func (ribbon *ribbonData) SuspendRepaint() {
	ribbon.suspendRepaint = true
}

func (ribbon *ribbonData) ResumeRepaint() {
	ribbon.suspendRepaint = false
	ribbon.Repaint()
}

func (ribbon *ribbonData) IsRepaintSuspended() bool {
	return ribbon.suspendRepaint
}

func (ribbon *ribbonData) AddTab(text string) *RibbonTab {
	tab := &RibbonTab{}
	tab.init(ribbon, text)
	tab.x = 0
	tab.width = tabHeaderWidth
	ribbon.tabs = append(ribbon.tabs, tab)
	return tab
}

// SetCurrentTab sets the given tab as selected
func (ribbon *ribbonData) SetCurrentTab(tab *RibbonTab) {
	ribbon.tabSelected = tab
	if !ribbon.IsRepaintSuspended() {
		ribbon.Repaint()
	}
}

// SetApplicationMenu sets the main menu
func (ribbon *ribbonData) SetApplicationMenu(caption string, menu PopupMenu) {
	ribbon.appMenu = menu
	ribbon.appMenuCaption = caption
	if !ribbon.IsRepaintSuspended() {
		ribbon.Repaint()
	}
}

func (ribbon *ribbonData) removeStateAllButtons(state uint8) {
	for _, sec := range ribbon.tabSelected.sections {
		for _, ibutton := range sec.buttons {
			button := ibutton.(*RibbonButtonData)
			button.removeState(state)
		}
	}
}

func (ribbon *ribbonData) mouseDown(pt *Point, mbutton int) {
	win.SetCapture(ribbon.GetHandle())
	ribbon.removeStateAllButtons(RibbonButtonStatePressed)
	if mbutton == MouseButtonLeft {
		if pt.IsInsideRect(&ribbon.rectAppMenu) {
			rect := ribbon.rectAppMenu
			ribbonRect := ribbon.GetWindowRect()
			ribbon.appMenu.Popup(ribbonRect.Left+rect.Left+1, ribbonRect.Top+rect.Bottom)
		} else {
			for _, tab := range ribbon.tabs {
				if pt.IsInsideRect(&tab.rectHeader) {
					ribbon.tabSelected = tab
					break
				}
			}
			for _, sec := range ribbon.tabSelected.sections {
				for _, ibutton := range sec.buttons {
					button := ibutton.(*RibbonButtonData)
					rect := button.rect
					if pt.IsInsideRect(&rect) {
						if !button.enabled {
							if button.splitDropDown && pt.IsInsideRect(&button.splitRects[1]) {
								button.addState(RibbonButtonStatePressed)
								button.splitStateIndex = 1
							}
						} else {
							button.addState(RibbonButtonStatePressed)
							button.splitStateIndex = 0
							if pt.IsInsideRect(&button.splitRects[1]) {
								button.splitStateIndex = 1
							}
						}
						fShowDropDown := func() {
							ribbon.Repaint()
							ribbonRect := ribbon.GetWindowRect()
							button.dropDownMenu.Popup(ribbonRect.Left+rect.Left, ribbonRect.Top+rect.Bottom)
							button.removeState(RibbonButtonStateHot)
							button.removeState(RibbonButtonStatePressed)
						}
						if button.splitDropDown {
							if button.splitStateIndex == 1 {
								fShowDropDown()
							}

						} else if button.dropDownMenu != nil {
							if button.enabled {
								fShowDropDown()
							}
						}
						break
					}
				}
			}
		}
	}
	ribbon.Repaint()
}

func (ribbon *ribbonData) mouseUp(pt *Point, mbutton int) {
	ribbon.removeStateAllButtons(RibbonButtonStatePressed)
	if mbutton == MouseButtonLeft {
		for _, sec := range ribbon.tabSelected.sections {
			for _, ibutton := range sec.buttons {
				button := ibutton.(*RibbonButtonData)
				rect := button.rect
				if button.enabled {
					if pt.IsInsideRect(&rect) {
						if button.dropDownMenu == nil {
							log.Printf("Clicked on '%s'\n", button.caption)
							if button.checkbox {
								button.ToggleState(RibbonButtonStateChecked)
							}
							if button.clickEvent != nil {
								e := RibbonButtonEvent{Button: button}
								button.clickEvent(&e)
							}
						} else if button.splitDropDown && pt.IsInsideRect(&button.splitRects[0]) {
							log.Printf("Clicked on '%s'\n", button.caption)
							if button.clickEvent != nil {
								e := RibbonButtonEvent{Button: button}
								button.clickEvent(&e)
							}
						}
						break
					}
				}
			}
		}
	}
	ribbon.Repaint()
	win.ReleaseCapture()
}

func (ribbon *ribbonData) mouseMove(pt *Point, mbutton int) {
	ribbon.appMenuState = 0
	for _, tab := range ribbon.tabs {
		tab.highlight = false
	}
	if pt.IsInsideRect(&ribbon.rectAppMenu) {
		ribbon.appMenuState = 1
	} else {
		for _, tab := range ribbon.tabs {
			if pt.IsInsideRect(&tab.rectHeader) {
				tab.highlight = true
				break
			}
		}
		ribbon.removeStateAllButtons(RibbonButtonStateHot)
		for _, sec := range ribbon.tabSelected.sections {
			for _, ibutton := range sec.buttons {
				button := ibutton.(*RibbonButtonData)
				rect := button.rect
				if pt.IsInsideRect(&rect) {
					button.splitStateIndex = 0
					if pt.IsInsideRect(&button.splitRects[1]) {
						button.splitStateIndex = 1
					}
					if mbutton == MouseButtonNone {
						button.addState(RibbonButtonStateHot)
					}
				}
			}
		}
	}
	ribbon.Repaint()
}

func (ribbon *ribbonData) measureRibbonSections(g *Graphics, rcTab *Rect) {
	tab := ribbon.tabSelected
	if tab == nil {
		log.Println("bug in the application - 'ribbon.tabSelected' is nil")
		return
	}

	newLeft := rcTab.Left
	sectionWidth := 0

	for _, section := range tab.sections {
		// first update sections top, bottom (which is independent of content size)
		section.rect.Top = rcTab.Top
		section.rect.Bottom = rcTab.Bottom
		section.rect.Left = newLeft

		// now update contents and update sections width accordingly
		contentWidth := section.measureButtons(g)
		if len(section.buttons) != 0 {
			sectionWidth = contentWidth
		} else {
			sectionWidth = 150 // temporary size for test purpose!!!
		}

		section.rect.Right = section.rect.Left + sectionWidth

		newLeft += sectionWidth
	}
}

func (ribbon *ribbonData) measureTabHeaders() {
	newLeft := tabHeaderWidth
	for _, tab := range ribbon.tabs {
		tab.x = newLeft
		tab.rectHeader = Rect{Left: tab.x, Top: 0,
			Right: tab.x + tab.width, Bottom: tabHeaderHeight}
		newLeft += tab.width
	}
}

func (sec *RibbonSection) measureButtons(g *Graphics) (width int) {
	// Update buttons
	sectionMargin := 4
	buttonMargin := 2
	buttonLeft := buttonMargin + sectionMargin
	buttonTop := buttonMargin + sectionMargin
	buttonSmallSize := 20
	const bigButtonHeight = 66
	buttonRowCount := 0
	contentWidth := buttonMargin + sectionMargin
	lastButtonWidth := 0
	font := sec.tab.ribbon.GetFont()
	for _, ibutton := range sec.buttons {
		button := ibutton.(*RibbonButtonData)
		//style := button.style
		size := button.size
		// Small buttons are 20x20 in size
		if size == RibbonButtonSizeSmall {

			//if button.image != nil {
			//buttonSmallSize = 24
			//}

			button.rect.Left = sec.rect.Left + buttonLeft
			button.rect.Right = button.rect.Left + buttonSmallSize
			button.rect.Top = sec.rect.Top + buttonTop
			button.rect.Bottom = button.rect.Top + buttonSmallSize

			if sec.twoRow {
				button.rect.Top += 4
				button.rect.Bottom += 4
			}

			if buttonRowCount == 0 {
				contentWidth += (button.rect.Right - button.rect.Left) + buttonMargin
			}
			maxRow := 3
			if sec.twoRow {
				maxRow = 2
			}
			if buttonRowCount < (maxRow - 1) {
				buttonTop += buttonSmallSize + buttonMargin
				buttonRowCount++
			} else {
				buttonTop = buttonMargin + sectionMargin
				buttonLeft += buttonSmallSize + buttonMargin
				buttonRowCount = 0
			}
		} else if size == RibbonButtonSizeMedium {
			button.rect.Left = sec.rect.Left + buttonLeft
			button.rect.Right = button.rect.Left + buttonSmallSize
			button.rect.Top = sec.rect.Top + buttonTop
			button.rect.Bottom = button.rect.Top + buttonSmallSize

			// Get the text size
			rectText := g.MeasureText(button.caption, win.DT_LEFT, font)

			textWidth := rectText.Width()
			button.rect.Right += textWidth + 6

			if button.dropDownMenu != nil {
				button.rect.Right += 10
			}

			buttonWidth := button.rect.Right - button.rect.Left

			if buttonRowCount == 0 {
				contentWidth += buttonWidth + buttonMargin
			} else if buttonWidth > lastButtonWidth {
				// Might happen because of different caption size
				contentWidth += buttonWidth - lastButtonWidth
			}
			lastButtonWidth = buttonWidth
			if buttonRowCount < 2 {
				buttonTop += buttonSmallSize + buttonMargin
				buttonRowCount++
			} else {
				buttonTop = buttonMargin + sectionMargin
				buttonLeft += buttonSmallSize + buttonMargin
				buttonRowCount = 0
			}
		} else if size == RibbonButtonSizeBig {
			// Get the text size
			rectText := g.MeasureText(button.caption, win.DT_CENTER|win.DT_TOP, font)
			textWidth := rectText.Width()

			buttonWidth := textWidth + 12

			if buttonWidth < 42 {
				buttonWidth = 42
			}

			// Reset top for big button
			buttonTop = buttonMargin + sectionMargin

			button.rect.Left = sec.rect.Left + buttonLeft
			button.rect.Right = button.rect.Left + buttonWidth
			button.rect.Top = sec.rect.Top + buttonTop
			button.rect.Bottom = button.rect.Top + bigButtonHeight

			button.rect.Top--
			button.rect.Bottom--

			// Update the split rects too

			// Upper rect
			button.splitRects[0] = button.rect
			button.splitRects[0].Bottom -= 28

			// Lower rect
			button.splitRects[1] = button.rect
			button.splitRects[1].Top = button.splitRects[1].Bottom - 29

			buttonLeft += buttonWidth + buttonMargin
			contentWidth += buttonWidth + buttonMargin
		}
	}
	contentWidth += sectionMargin
	// return the amount of width need for all the buttons
	return contentWidth
}

var textColor = Rgb(65, 65, 65)

func (sec *RibbonSection) drawButtons(g *Graphics) {
	ribbon := sec.tab.ribbon
	font := ribbon.GetFont()
	classButtonUtf16, _ := syscall.UTF16PtrFromString("BUTTON")
	htheme := win.OpenThemeData(ribbon.GetHandle(), classButtonUtf16)

	for _, ibutton := range sec.buttons {
		button := ibutton.(*RibbonButtonData)
		size := button.size
		rectButton := button.rect
		// Special case for a color button
		if !button.checkbox {
			// Small color button
			if size == RibbonButtonSizeSmall && button.image == nil {
				if button.hasState(RibbonButtonStateHot) {
					g.DrawFillRectangle(&rectButton, NewRgb(100, 165, 231), NewRgb(255, 255, 255))
				} else {
					g.DrawFillRectangle(&rectButton, &button.border, NewRgb(255, 255, 255))
				}
			} else {
				// Draw hover highlight/toggled
				if button.hasState(RibbonButtonStateChecked) {
					g.DrawFillRectangle(&rectButton, NewRgb(98, 162, 228), NewRgb(201, 224, 247))
					if button.splitDropDown {
						rectUpper := button.splitRects[0]
						g.DrawFillRectangle(&rectUpper, NewRgb(98, 162, 228), NewRgb(201, 224, 247))
					}
				}
				if button.hasState(RibbonButtonStateHot) {
					if button.hasState(RibbonButtonStateChecked) {
						if button.splitDropDown {
							// if mouse on upper
							if button.splitStateIndex == 0 {
								rectUpper := button.splitRects[0]
								rectLower := button.splitRects[1]
								g.DrawFillRectangle(&rectLower, NewRgb(147, 190, 234), NewRgb(245, 246, 247))
								g.DrawFillRectangle(&rectUpper, NewRgb(122, 176, 231), NewRgb(213, 230, 247))
							} else {
								rectUpper := button.splitRects[0]
								rectLower := button.splitRects[1]
								g.DrawFillRectangle(&rectLower, NewRgb(164, 206, 249), NewRgb(232, 239, 247))
								g.DrawFillRectangle(&rectUpper, NewRgb(122, 176, 231), NewRgb(213, 230, 247))
							}
						} else {
							g.DrawFillRectangle(&rectButton, NewRgb(122, 176, 231), NewRgb(213, 230, 247))
						}
					} else {
						if button.splitDropDown {
							if button.enabled {
								// if mouse on upper
								if button.splitStateIndex == 0 {
									rectUpper := button.splitRects[0]
									rectLower := button.splitRects[1]
									g.DrawFillRectangle(&rectLower, NewRgb(164, 206, 249), NewRgb(245, 246, 247))
									g.DrawFillRectangle(&rectUpper, NewRgb(164, 206, 249), NewRgb(232, 239, 247))
								} else {
									rectUpper := button.splitRects[0]
									rectLower := button.splitRects[1]
									g.DrawFillRectangle(&rectUpper, NewRgb(164, 206, 249), NewRgb(245, 246, 247))
									g.DrawFillRectangle(&rectLower, NewRgb(164, 206, 249), NewRgb(232, 239, 247))
								}
							} else {
								rectUpper := button.splitRects[0]
								rectLower := button.splitRects[1]
								if button.splitStateIndex == 0 {
									g.DrawFillRectangle(&rectLower, NewRgb(164, 206, 249), NewRgb(245, 246, 247))
								} else {
									g.DrawFillRectangle(&rectLower, NewRgb(164, 206, 249), NewRgb(232, 239, 247))
								}
								g.DrawRectangle(&rectUpper, NewRgb(196, 211, 226))
							}
						} else {
							if button.enabled {
								g.DrawFillRectangle(&rectButton, NewRgb(164, 206, 249), NewRgb(232, 239, 247))
							}
						}
					}
				}
				if button.hasState(RibbonButtonStatePressed) {
					if button.splitDropDown {
						// if clicked on upper
						if button.splitStateIndex == 0 {
							g.DrawFillRectangle(&button.splitRects[0], NewRgb(98, 162, 228), NewRgb(201, 224, 247))
						} else {
							g.DrawFillRectangle(&button.splitRects[1], NewRgb(98, 162, 228), NewRgb(201, 224, 247))
						}
					} else {
						g.DrawFillRectangle(&rectButton, NewRgb(98, 162, 228), NewRgb(201, 224, 247))
					}
				}
			}
		}
		// Draw the actual button
		if size == RibbonButtonSizeSmall {
			if button.image != nil {
				g.DrawBitmapImageCenter(button.image, &rectButton, !button.enabled)
			} else {
				rectButton.Inflate(-2, -2)
				g.DrawFillRectangle(&rectButton, &button.color, &button.color)
			}
		} else if size == RibbonButtonSizeBig {
			// Draw button with a image
			if button.image != nil {
				rectUpper := button.splitRects[0]
				g.DrawBitmapImageCenter(button.image, &rectUpper, !button.enabled)
			} else {
				// Draw the button (No image, Only bordered rectangle)
				rectButton.Left += 6
				rectButton.Top += 4
				rectButton.Right -= 6
				width := rectButton.Width()
				rectButton.Bottom = rectButton.Top + width // height = width, keep it square
				g.DrawFillRectangle(&rectButton, &button.border, NewRgb(255, 255, 255))
				rectButton.Inflate(-2, -2)
				g.DrawFillRectangle(&rectButton, &button.color, &button.color)
			}
			//  draw the caption at bottom
			rectLower := button.splitRects[1]
			rectLower.Top += 1
			if button.enabled || button.splitDropDown {
				g.DrawText(button.caption, &rectLower, win.DT_CENTER|win.DT_TOP, &textColor, font)
				// Draw drop down arrow if there is a drop menu
				if button.dropDownMenu != nil {
					g.DrawDropDownArrow(rectLower.CenterX(), rectLower.Bottom-10, NewRgb(0, 0, 0))
				}
			} else {
				g.DrawText(button.caption, &rectLower, win.DT_CENTER|win.DT_TOP, NewRgb(140, 140, 140), font)
				// Draw drop down arrow if there is a drop menu
				if button.dropDownMenu != nil {
					g.DrawDropDownArrow(rectLower.CenterX(), rectLower.Bottom-10, NewRgb(180, 180, 180))
				}
			}
		} else if size == RibbonButtonSizeMedium {
			if button.checkbox {
				//checkboxSize := 13
				//centerY := rectButton.CenterY() - (checkboxSize / 2)
				//rectCheckbox := Rect{Left: rectButton.Left + 4, Top: centerY}
				//rectCheckbox.Right = rectCheckbox.Left + checkboxSize
				//rectCheckbox.Bottom = rectCheckbox.Top + checkboxSize

				var rectBtn win.RECT
				cx := win.GetSystemMetrics(win.SM_CXMENUCHECK)
				cy := win.GetSystemMetrics(win.SM_CYMENUCHECK)

				rectBtn.Left = int32(rectButton.Left)
				rectBtn.Top = int32(rectButton.CenterY() - (int(cy) / 2))
				rectBtn.Right = rectBtn.Left + cx
				rectBtn.Bottom = rectBtn.Top + cy

				if !button.IsEnabled() {
					if button.hasState(RibbonButtonStateChecked) {
						win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_CHECKEDDISABLED, &rectBtn, nil)
					} else {
						win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_UNCHECKEDDISABLED, &rectBtn, nil)
					}
				} else {
					if button.hasState(RibbonButtonStatePressed) {
						if button.hasState(RibbonButtonStateChecked) {
							win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_CHECKEDPRESSED, &rectBtn, nil)
						} else {
							win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_UNCHECKEDPRESSED, &rectBtn, nil)
						}
					} else if button.hasState(RibbonButtonStateHot) {
						if button.hasState(RibbonButtonStateChecked) {
							win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_CHECKEDHOT, &rectBtn, nil)
						} else {
							win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_UNCHECKEDHOT, &rectBtn, nil)
						}
						//g.DrawFillRectangle(&rectCheckbox, Rgb(51, 153, 255), Rgb(255, 255, 255))
					} else {
						if button.hasState(RibbonButtonStateChecked) {
							win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_CHECKEDNORMAL, &rectBtn, nil)
						} else {
							win.DrawThemeBackground(htheme, g.GetHDC(), win.BP_CHECKBOX, win.CBS_UNCHECKEDNORMAL, &rectBtn, nil)
						}
						//g.DrawFillRectangle(&rectCheckbox, Rgb(97, 121, 160), Rgb(255, 255, 255))
					}
				}
				/*
					if button.hasState(RibbonButtonStatePressed) {
						g.DrawFillRectangle(&rectCheckbox, Rgb(0, 124, 222), Rgb(217, 236, 255))
					}
					if button.hasState(RibbonButtonStateChecked) {
						g.DrawBitmapImage(imageCheck, rectCheckbox.Left+1, rectCheckbox.Top+1, !button.enabled)
					}
					rectButton.Left += 24
				*/
				rectButton.Left = int(rectBtn.Right) + CheckboxAndTextGap
				if button.enabled {
					g.DrawText(button.caption, &rectButton, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, &textColor, font)
				} else {
					g.DrawText(button.caption, &rectButton, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, NewRgb(140, 140, 140), font)
				}

			} else {
				if button.image != nil {
					rectButton.Left += 2
					centeredTop := rectButton.CenterY() - (button.image.Height / 2)
					g.DrawBitmapImage(button.image, rectButton.Left, centeredTop, !button.enabled)
					rectButton.Left += button.image.Width + 4
				}
				if button.enabled {
					g.DrawText(button.caption, &rectButton, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, &textColor, font)
				} else {
					g.DrawText(button.caption, &rectButton, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, NewRgb(140, 140, 140), font)
				}
				// Draw drop down arrow if there is a drop menu
				if button.dropDownMenu != nil {
					g.DrawDropDownArrow(rectButton.Right-8, rectButton.CenterY()-2, NewRgb(0, 0, 0))
				}
			}
		}
	}
	win.CloseThemeData(htheme)
}

func (ribbon *ribbonData) paint(gOrg *Graphics, rc *Rect) {
	var db DoubleBuffer
	g := db.BeginDoubleBuffer(gOrg, rc, NewRgb(245, 245, 245), NewRgb(245, 245, 245))
	defer db.EndDoubleBuffer()
	const textHeight = 16
	borderColor := Rgb(220, 220, 220)

	if !ribbon.initialized {
		return
	}

	font := ribbon.GetFont()

	// Update
	ribbon.rectAppMenu = Rect{Left: 0, Top: 0, Right: tabHeaderWidth, Bottom: tabHeaderHeight}
	ribbon.measureTabHeaders()

	// Draw the application menu
	appMenuColor := Rgb(25, 121, 202)
	if ribbon.appMenuState == 1 {
		appMenuColor = Rgb(41, 140, 225)
	}
	g.DrawFillRectangle(&ribbon.rectAppMenu, &appMenuColor, &appMenuColor)
	g.DrawText(ribbon.appMenuCaption, &ribbon.rectAppMenu, win.DT_CENTER|win.DT_VCENTER|win.DT_SINGLELINE, NewRgb(255, 255, 255), font)

	// Draw all the tab headers with their captions
	for _, tab := range ribbon.tabs {
		rectTabHead := tab.rectHeader
		g.DrawText(tab.caption, &rectTabHead, win.DT_CENTER|win.DT_VCENTER|win.DT_SINGLELINE, &textColor, font)

		if tab == ribbon.tabSelected {
			// top border
			g.DrawLine(tab.x, 0, tab.x+tab.width, 0, &borderColor)
			// left border
			g.DrawLine(tab.x, 0, tab.x, tabHeaderHeight, &borderColor)
			// right border
			g.DrawLine(tab.x+tab.width, 0, tab.x+tab.width, tabHeaderHeight, &borderColor)
			// left bottom line
			g.DrawLine(tab.x, tabHeaderHeight, 0, tabHeaderHeight, &borderColor)
			// right bottom line
			g.DrawLine(tab.x+tab.width, tabHeaderHeight, rc.Right, tabHeaderHeight, &borderColor)
		} else if tab.highlight {
			tabHeaderHighlightColor := Rgb(225, 225, 225)
			// top border
			g.DrawLine(tab.x, 0, tab.x+tab.width, 0, &tabHeaderHighlightColor)
			// left border
			g.DrawLine(tab.x+2, 0, tab.x+2, tabHeaderHeight, &tabHeaderHighlightColor)
			// right border
			g.DrawLine((tab.x+tab.width)-2, 0, (tab.x+tab.width)-2, tabHeaderHeight, &tabHeaderHighlightColor)
		}
	}

	rcTab := Rect{Left: 0, Top: tabHeaderHeight + 4, Right: rc.Right, Bottom: rc.Bottom - textHeight}
	ribbon.measureRibbonSections(g, &rcTab)

	for _, section := range ribbon.tabSelected.sections {
		seperatorX := section.rect.Right
		g.DrawLine(seperatorX, rcTab.Top, seperatorX, rc.Bottom-4, &borderColor)
		if len(section.caption) > 0 {
			g.DrawText(section.caption,
				&Rect{Left: section.rect.Left, Top: section.rect.Bottom - 1,
					Right: section.rect.Right, Bottom: section.rect.Bottom + textHeight},
				win.DT_CENTER|win.DT_TOP|win.DT_SINGLELINE, NewRgb(90, 90, 90), font)
		}
		section.drawButtons(g)
	}
	// Bottom border
	g.DrawLine(0, rc.Bottom-1, rc.Right, rc.Bottom-1, &borderColor)
}
