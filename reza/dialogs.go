package reza

import (
	"syscall"
	"unsafe"

	win "github.com/lxn/win"
)

type Dialog interface {
	// Embedd the window interface
	Window
	Initialize(parent Window, caption string, width, height int)
	Show(modal bool, fOnAccept func())
	AddWidgets(widgets []Widget)
}

type dialogData struct {
	// Embedd the window interface
	Window
	modal     bool
	btnOK     Button
	btnCancel Button
	fOnAccept func()
}

func NewDialog() Dialog {
	return &dialogData{Window: NewWindow()}
}

func (dlg *dialogData) Initialize(parent Window, caption string, width, height int) {
	logInfo("init dialog...")

	dlg.Create(caption, win.WS_CAPTION|win.WS_SYSMENU, 600, 400, width, height, parent)
	dlg.SetCloseEventHandler(func() bool {
		dlg.Hide()
		// Don't destroy just keep it hidden
		return false
	})

	const buttonWidth = 86
	const buttonHeight = 24

	container := CreateContainer(FlowRightToLeft, 0, 0, 0, 50, dlg) // only height counts for 'DockBottom'
	container.SetDockType(DockBottom)

	dlg.btnOK = CreateButton("OK", 0, 0, buttonWidth, buttonHeight, win.BS_PUSHBUTTON, container)
	dlg.btnCancel = CreateButton("Cancel", 0, 0, buttonWidth, buttonHeight, win.BS_DEFPUSHBUTTON, container)
	dlg.btnCancel.SetMargin(10, 10, 10, 10)
	dlg.btnOK.SetMargin(10, 10, 10, 10)
	dlg.btnOK.SetClickEventHandler(func(sender Button) {
		if dlg.fOnAccept != nil {
			dlg.fOnAccept()
		}
		dlg.Hide()
	})
	dlg.btnCancel.SetClickEventHandler(func(sender Button) {
		dlg.Hide()
	})
	logInfo("Done initializing dialog")
}

func (dlg *dialogData) AddWidgets(widgets []Widget) {
	for _, widget := range widgets {
		widget.Create(dlg)
	}
}

func (dlg *dialogData) Hide() {
	if dlg.modal && dlg.HasParent() {
		win.EnableWindow(dlg.GetParent().GetHandle(), true)
	}
	logInfo("Done with the dialog")
	dlg.SetVisible(false)
}

func (dlg *dialogData) Show(modal bool, fOnAccept func()) {
	logInfo("Show dialog")
	if modal && dlg.HasParent() {
		win.EnableWindow(dlg.GetParent().GetHandle(), false)
	}
	dlg.modal = modal
	dlg.fOnAccept = fOnAccept
	dlg.RequestLayout()
	dlg.SetVisible(true)
}

func fileDialogHook(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_NOTIFY:
		logInfo("WM_NOTIFY")
	}
	return 0
}

func fileDialog(owner Window, filePath, filter string, filterIndex int, fun func(ofn *win.OPENFILENAME) bool) (filepath string, accepted bool) {
	var ofn win.OPENFILENAME // common dialog box structure
	filterFinalUTF16 := make([]uint16, len(filter)+2)
	_filterUTF16, _ := syscall.UTF16FromString(filter)
	copy(filterFinalUTF16, _filterUTF16)
	// Replace '|' with the expected '\0'.
	for i, c := range filter {
		if byte(c) == '|' {
			filterFinalUTF16[i] = uint16(0)
		}
	}
	// Initialize OPENFILENAME
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))
	ofn.HwndOwner = owner.GetHandle()
	filePathFinalUTF16 := make([]uint16, 1024)
	_filePathUTF16, _ := syscall.UTF16FromString(filePath)
	copy(filePathFinalUTF16, _filePathUTF16)
	ofn.LpstrFile = &filePathFinalUTF16[0]
	ofn.NMaxFile = uint32(len(filePathFinalUTF16))
	ofn.LpstrFilter = &filterFinalUTF16[0]
	ofn.NFilterIndex = uint32(filterIndex)
	ofn.LpstrFileTitle = nil
	ofn.NMaxFileTitle = 0
	ofn.LpstrInitialDir = nil
	ofn.Flags = win.OFN_PATHMUSTEXIST | win.OFN_FILEMUSTEXIST // | win.OFN_ENABLEHOOK | win.OFN_EXPLORER | win.OFN_ENABLESIZING
	//ofn.LpfnHook = win.LPOFNHOOKPROC(syscall.NewCallback(fileDialogHook))
	// Display the Open dialog box.
	if fun(&ofn) {
		return syscall.UTF16ToString(filePathFinalUTF16), true
	}
	return "", false
}

// OpenFileDialog lets browse files through a dialog
func OpenFileDialog(owner Window, filter string, filterIndex int) (filepath string, accepted bool) {
	return fileDialog(owner, "", filter, filterIndex, win.GetOpenFileName)
}

// SaveFileDialog lets browse files through a dialog
func SaveFileDialog(owner Window, filename, filter string, filterIndex int) (filepath string, accepted bool) {
	return fileDialog(owner, filename, filter, filterIndex, win.GetSaveFileName)
}

// ChoseColorDialog lets choose colors from a color dialog
func ChoseColorDialog(owner Window, color Color, customColors [16]Color) (result Color, ccolors [16]Color) {
	var cc win.CHOOSECOLOR
	var acrCustClr [16]win.COLORREF
	var newCustomColors [16]Color
	for i := range acrCustClr {
		acrCustClr[i] = customColors[i].AsCOLORREF()
	}
	cc.LStructSize = uint32(unsafe.Sizeof(cc))
	cc.HwndOwner = owner.GetHandle()
	cc.RgbResult = color.AsCOLORREF()
	cc.Flags = win.CC_FULLOPEN | win.CC_RGBINIT
	cc.LpCustColors = &acrCustClr
	if win.ChooseColor(&cc) {
		rgba := FromCOLORREF(cc.RgbResult)
		for i := range acrCustClr {
			newCustomColors[i] = FromCOLORREF(acrCustClr[i])
		}
		return rgba, newCustomColors
	}
	return color, newCustomColors
}
