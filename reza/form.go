package reza

import (
	win "github.com/lxn/win"
)

type Form interface {
	Window
	Initialize()
	Show() bool
	Hide()
}

type FormData struct {
	Window
}

func NewForm() Form {
	fd := &FormData{Window: NewWindow()}
	//fd.quickAccessToolbar = true
	return fd
}

func (f *FormData) Initialize() {

}

func (f *FormData) Hide() {
	f.SetVisible(false)
}

func (f *FormData) Show() bool {
	ret := win.ShowWindow(f.GetHandle(), win.SW_SHOWDEFAULT)
	if ret {
		ret = win.UpdateWindow(f.GetHandle())
		return ret
	}
	return false
}
