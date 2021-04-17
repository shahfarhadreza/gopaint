package reza

import (
	win "github.com/lxn/win"
)

// Statusbar is the statusbar at the bottom ;)
type Statusbar interface {
	// Inherit data from the window type
	Window
	AddStatus(iconpath string, text string) Status
}

type statusbarData struct {
	// Inherit data from the window type
	Window
	sections []*statusData
}

// Status is a status of a statusbar
type Status interface {
	Update(text string)
	SetVisible(visible bool)
}

type statusData struct {
	statusbar *statusbarData
	icon      *BitmapImage
	text      string
	visible   bool
}

func NewStatusbar(parent Window) Statusbar {
	bar := &statusbarData{Window: NewWindow()}
	bar.init(parent)
	return bar
}

// UpdateText updates the text of the section
func (section *statusData) Update(text string) {
	section.text = text
	section.statusbar.Repaint()
}

func (section *statusData) SetVisible(visible bool) {
	section.visible = visible
	section.statusbar.Repaint()
}

func (status *statusbarData) init(parent Window) Statusbar {
	logInfo("init Statusbar...")
	status.Window.Create("", win.WS_CHILD|win.WS_CLIPCHILDREN|win.WS_VISIBLE, 10, 10, 10, 10, parent)
	status.SetPaintEventHandler(status.paint)
	status.SetFont(app.GetGuiFont().GetHandle())
	status.sections = make([]*statusData, 0)
	logInfo("Done initializing Statusbar")
	return status
}

// AddSection adds new section to the statusbar
func (status *statusbarData) AddStatus(iconpath string, text string) Status {
	section := &statusData{}
	section.statusbar = status
	section.icon, _ = CreateBitmapImage(iconpath, false)
	section.text = text
	section.visible = true
	status.sections = append(status.sections, section)
	return section
}

func (status *statusbarData) paint(gOrg *Graphics, rect *Rect) {
	var db DoubleBuffer
	g := db.BeginDoubleBuffer(gOrg, rect, NewRgb(240, 240, 240), NewRgb(240, 240, 240))
	defer db.EndDoubleBuffer()
	g.DrawLine(rect.Left, rect.Top, rect.Right, rect.Top, NewRgb(215, 215, 215))
	// draw the grip
	font := status.GetFont()
	const sectionWidth = 140
	sectionLeft := 8
	for _, section := range status.sections {
		if section.icon != nil {
			if section.visible {
				centerY := rect.CenterY() - (section.icon.Height / 2)
				g.DrawBitmapImage(section.icon, sectionLeft, centerY, false)
			}
			sectionLeft += section.icon.Width + 5
		}
		if section.visible {
			if len(section.text) > 0 {
				rcText := Rect{Left: sectionLeft, Top: rect.Top - 1, Right: sectionLeft + sectionWidth, Bottom: rect.Bottom - 1}
				g.DrawText(section.text, &rcText, win.DT_LEFT|win.DT_VCENTER|win.DT_SINGLELINE, NewRgb(0, 0, 0), font)
			}
		}
		sectionLeft += sectionWidth
		g.DrawLine(sectionLeft, rect.Top+3, sectionLeft, rect.Bottom-3, NewRgb(215, 215, 215))
		sectionLeft += 8
	}
}
