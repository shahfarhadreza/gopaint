package reza

// RibbonTab contains data for each tab of a ribbon
type RibbonTab struct {
	caption    string
	x, width   int
	rectHeader Rect
	highlight  bool
	sections   []*RibbonSection
	ribbon     Ribbon // pointer to the ribbon this tab belongs to
}

// RibbonSection contains data for each section of a ribbon tab
type RibbonSection struct {
	caption string
	rect    Rect
	buttons []RibbonButton
	twoRow  bool
	tab     *RibbonTab // pointer to the ribbon tab this section belongs to
}

const (
	RibbonButtonSizeSmall  = 210
	RibbonButtonSizeMedium = 211
	RibbonButtonSizeBig    = 212
)

const (
	RibbonButtonStateHot     = 1 << 0
	RibbonButtonStatePressed = 1 << 1
	RibbonButtonStateChecked = 1 << 2
)

// RibbonButton is a ribbon button item
type RibbonButton interface {
	GetText() string
	SetText(text string)
	SetClickEvent(f func(e *RibbonButtonEvent))
	SetDropdownMenu(menu PopupWindow, split bool)
	SetEnabled(enabled bool)
	IsEnabled() bool
	GetColor() Color
	SetColor(c Color)
	SetIcon(image *BitmapImage)
	GetIcon() *BitmapImage
	SetToggled(toggled bool)
	IsToggled() bool
	Dispose()
}

// RibbonButtonData is a ribbon button item data
type RibbonButtonData struct {
	rect            Rect    // defines the entire button
	splitRects      [2]Rect // splitted rects for a split drop down button
	state           uint8
	splitStateIndex int // 0 for uppper rect, 1 for lower rect
	size            int
	caption         string
	border          Color
	color           Color
	image           *BitmapImage
	checkbox        bool
	enabled         bool
	dropDownMenu    PopupWindow
	splitDropDown   bool
	clickEvent      func(e *RibbonButtonEvent)
	section         *RibbonSection
}

// RibbonColorButton defines a colored button
type RibbonColorButton struct {
	Name    string
	Color   Color
	OnClick func(e *RibbonButtonEvent)
}

// RibbonButtonEvent holds information about click event of a ribbon button
type RibbonButtonEvent struct {
	Button RibbonButton
}

func (b *RibbonButtonData) Dispose() {
	if b.image != nil {
		b.image.Dispose()
	}
}

func (b *RibbonButtonData) GetText() string {
	return b.caption
}
func (b *RibbonButtonData) SetText(text string) {
	b.caption = text
	if !b.section.tab.ribbon.IsRepaintSuspended() {
		b.section.tab.ribbon.Repaint()
	}
}

func (b *RibbonButtonData) SetClickEvent(f func(e *RibbonButtonEvent)) {
	b.clickEvent = f
}

func (b *RibbonButtonData) SetIcon(image *BitmapImage) {
	b.image = image
}
func (b *RibbonButtonData) GetIcon() *BitmapImage {
	return b.image
}

func (b *RibbonButtonData) GetColor() Color {
	return b.color
}

func (b *RibbonButtonData) SetColor(c Color) {
	b.color = c
	if !b.section.tab.ribbon.IsRepaintSuspended() {
		b.section.tab.ribbon.Repaint()
	}
}

func (b *RibbonButtonData) SetToggled(toggled bool) {
	if toggled {
		b.addState(RibbonButtonStateChecked)
	} else {
		b.removeState(RibbonButtonStateChecked)
	}
	if !b.section.tab.ribbon.IsRepaintSuspended() {
		b.section.tab.ribbon.Repaint()
	}
}
func (b *RibbonButtonData) IsToggled() bool {
	return b.hasState(RibbonButtonStateChecked)
}

func (b *RibbonButtonData) addState(s uint8) {
	b.state |= s
}

func (b *RibbonButtonData) removeState(s uint8) {
	b.state = b.state &^ s
}

func (b *RibbonButtonData) ToggleState(s uint8) {
	b.state = b.state ^ s
}

func (b *RibbonButtonData) hasState(s uint8) bool {
	return (b.state & s) != 0
}

func (b *RibbonButtonData) SetEnabled(enabled bool) {
	b.enabled = enabled
	if !b.section.tab.ribbon.IsRepaintSuspended() {
		b.section.tab.ribbon.Repaint()
	}
}

func (b *RibbonButtonData) IsEnabled() bool {
	return b.enabled
}

// SetDropdownMenu sets the menu
func (b *RibbonButtonData) SetDropdownMenu(menu PopupWindow, split bool) {
	b.dropDownMenu = menu
	b.splitDropDown = split
}

func (sec *RibbonSection) init(tab *RibbonTab, text string) {
	sec.buttons = make([]RibbonButton, 0)
	sec.caption = text
	sec.twoRow = false
	sec.tab = tab
}

func (sec *RibbonSection) SetTwoRow(twoRow bool) {
	sec.twoRow = twoRow
}

func (sec *RibbonSection) AddCheckButton(text string, checked bool) RibbonButton {
	button := &RibbonButtonData{caption: text, size: RibbonButtonSizeMedium}
	button.section = sec
	button.checkbox = true
	button.enabled = true
	if checked {
		button.addState(RibbonButtonStateChecked)
	}
	sec.buttons = append(sec.buttons, button)
	return button
}

func (sec *RibbonSection) AddColorButtons(cbuttons []RibbonColorButton) []RibbonButton {
	border := Rgb(160, 160, 160)
	buttons := make([]RibbonButton, len(cbuttons))
	for i, b := range cbuttons {
		button := sec.AddButton(b.Name, border, b.Color, RibbonButtonSizeSmall)
		button.SetClickEvent(b.OnClick)
		buttons[i] = button
	}
	return buttons
}

func (sec *RibbonSection) AddButton(text string, borderColor, c Color, sz int) RibbonButton {
	button := &RibbonButtonData{caption: text, border: borderColor, color: c, size: sz}
	button.section = sec
	button.enabled = true
	sec.buttons = append(sec.buttons, button)
	return button
}

func (sec *RibbonSection) AddImageButton(text string, filename string, tsize int) RibbonButton {
	button := &RibbonButtonData{caption: text, size: tsize}
	button.section = sec
	if len(filename) > 0 {
		button.image, _ = CreateBitmapImage(filename, true)
	}
	button.enabled = true
	sec.buttons = append(sec.buttons, button)
	return button
}

func (tab *RibbonTab) init(ribbon Ribbon, text string) {
	tab.sections = make([]*RibbonSection, 0)
	tab.ribbon = ribbon
	tab.caption = text
	tab.highlight = false
}

func (tab *RibbonTab) AddSection(text string) *RibbonSection {
	section := &RibbonSection{}
	section.init(tab, text)
	tab.sections = append(tab.sections, section)
	return section
}
