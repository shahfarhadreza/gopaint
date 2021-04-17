package reza

import "github.com/lxn/win"

type Widget interface {
	Create(parent Window) Window
}

type WGroup struct {
	DockType int
	Margins  Margins
	X, Y     int
	Width    int
	Height   int
	Text     string
	Widgets  []Widget
}

func (w *WGroup) Create(parent Window) Window {
	window := CreateGroup(w.Text, w.X, w.Y, w.Width, w.Height, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	for _, widget := range w.Widgets {
		widget.Create(window)
	}
	return window
}

type WFlowContainer struct {
	DockType      int
	Margins       Margins
	X, Y          int
	Width         int
	Height        int
	FlowDirection int
	Widgets       []Widget
}

func (w *WFlowContainer) Create(parent Window) Window {
	window := CreateContainer(w.FlowDirection, w.X, w.Y, w.Width, w.Height, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	for _, widget := range w.Widgets {
		widget.Create(window)
	}
	return window
}

type WLabel struct {
	DockType int
	Margins  Margins
	X, Y     int
	Width    int
	Height   int
	Text     string
	AssignTo *Label
}

func (w *WLabel) Create(parent Window) Window {
	window := CreateLabel(w.Text, w.X, w.Y, w.Width, w.Height, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	if w.AssignTo != nil {
		*w.AssignTo = window
	}
	return window
}

type WButton struct {
	DockType int
	Margins  Margins
	X, Y     int
	Width    int
	Height   int
	Text     string
	OnClick  func(sender Button)
	AssignTo *Button
}

func (w *WButton) Create(parent Window) Window {
	window := CreateButton(w.Text, w.X, w.Y, w.Width, w.Height, win.BS_PUSHBUTTON, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	window.SetClickEventHandler(w.OnClick)
	if w.AssignTo != nil {
		*w.AssignTo = window
	}
	return window
}

type WRadioButton struct {
	DockType int
	Margins  Margins
	X, Y     int
	Width    int
	Height   int
	Text     string
	Checked  bool
	AssignTo *Button
}

func (w *WRadioButton) Create(parent Window) Window {
	style := uint(win.BS_AUTORADIOBUTTON)
	window := CreateButton(w.Text, w.X, w.Y, w.Width, w.Height, style, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	if w.Checked {
		win.SendMessage(window.GetHandle(), win.BM_SETCHECK, win.BST_CHECKED, 0)
	}
	if w.AssignTo != nil {
		*w.AssignTo = window
	}
	return window
}

type WCheckButton struct {
	DockType int
	Margins  Margins
	X, Y     int
	Width    int
	Height   int
	Text     string
	Checked  bool
	AssignTo *Button
}

func (w *WCheckButton) Create(parent Window) Window {
	window := CreateButton(w.Text, w.X, w.Y, w.Width, w.Height, win.BS_AUTOCHECKBOX, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	if w.Checked {
		win.SendMessage(window.GetHandle(), win.BM_SETCHECK, win.BST_CHECKED, 0)
	}
	if w.AssignTo != nil {
		*w.AssignTo = window
	}
	return window
}

type WTextBox struct {
	DockType int
	Margins  Margins
	X, Y     int
	Width    int
	Height   int
	Text     string
	AssignTo *TextBox
}

func (w *WTextBox) Create(parent Window) Window {
	window := CreateTextBox(w.Text, w.X, w.Y, w.Width, w.Height, win.WS_BORDER|win.ES_LEFT, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	if w.AssignTo != nil {
		*w.AssignTo = window
	}
	return window
}

type WImageViewer struct {
	Path     string
	X, Y     int
	DockType int
	Margins  Margins
	AssignTo *ImageViewer
}

func (w *WImageViewer) Create(parent Window) Window {
	window := CreateImageViewer(w.Path, w.X, w.Y, parent)
	window.SetDockType(w.DockType)
	window.SetMargin(w.Margins.Left, w.Margins.Right, w.Margins.Top, w.Margins.Bottom)
	if w.AssignTo != nil {
		*w.AssignTo = window
	}
	return window
}
