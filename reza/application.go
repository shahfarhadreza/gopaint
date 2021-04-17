package reza

import (
	"log"
	"syscall"

	"github.com/go-ole/go-ole"
	win "github.com/lxn/win"
)

type Application interface {
	SetMainWindow(mainWindow Form)
	GetMainWindow() Form
	GetGuiFont() *GdiFont
	GetCursorPos() Point
	GetAppInstance() win.HINSTANCE
	Run()
	Exit()
}

type applicationData struct {
	hInstance        win.HINSTANCE
	GUIFont          *GdiFont
	mainWindow       Form
	FontReferenceDPI uint32
	DPI              uint32
}

var app *applicationData = nil

func NewApplication() Application {
	if app != nil {
		log.Panicln("Multiple application instances!!!")
	}
	app = &applicationData{}
	app.hInstance = win.GetModuleHandle(nil)
	app.FontReferenceDPI = 72
	fDwmEnabled := false
	if DwmIsCompositionEnabled(&fDwmEnabled) >= 0 {
		if !fDwmEnabled {
			log.Panicln("DWM NOT ENABLED!")
		}
	} else {
		log.Panicln("DWM FAILED!")
	}
	logInfo("Initializing OLE ....")
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		log.Panicln(err)
	}
	logInfo("Initializing GDI+ ....")
	// GDI+
	var gdiplusStartupInput win.GdiplusStartupInput
	var gdiplusToken win.GdiplusStartupOutput
	gdiplusStartupInput.GdiplusVersion = 1
	if win.GdiplusStartup(&gdiplusStartupInput, &gdiplusToken) != 0 {
		log.Panicln("failed to initialize GDI+")
	}
	logInfo("Done initializing everything!!!")
	return app
}

func (app *applicationData) SetMainWindow(mainWindow Form) {
	app.mainWindow = mainWindow
	// Initialize the main window
	app.mainWindow.Create("Form", win.WS_OVERLAPPEDWINDOW|win.WS_CLIPCHILDREN|win.WS_VISIBLE, 10, 10, 10, 10, nil)
	app.DPI = win.GetDpiForWindow(app.mainWindow.GetHandle())
	if app.DPI == 0 {
		log.Panicln("error: Failed to get dpi of main window")
		return
	}
	app.initFont()
	app.mainWindow.Initialize()
}

type GdiFont struct {
	hfont win.HFONT
}

func CreateDPIAwareFont(fontName string, points int) *GdiFont {
	gdiFont := &GdiFont{}
	var lfFont win.LOGFONT
	src, _ := syscall.UTF16FromString(fontName)
	dest := lfFont.LfFaceName[:]
	copy(dest, src)
	lfFont.LfHeight = -win.MulDiv(int32(points), int32(app.DPI), int32(app.FontReferenceDPI))
	lfFont.LfWeight = win.FW_LIGHT
	lfFont.LfCharSet = win.ANSI_CHARSET
	lfFont.LfOutPrecision = win.OUT_DEFAULT_PRECIS
	lfFont.LfClipPrecision = win.CLIP_DEFAULT_PRECIS
	lfFont.LfQuality = win.CLEARTYPE_QUALITY
	gdiFont.hfont = win.CreateFontIndirect(&lfFont)
	return gdiFont
}

func (f *GdiFont) GetHandle() win.HFONT {
	return f.hfont
}

func (f *GdiFont) Dispose() {
	if f.hfont != 0 {
		win.DeleteObject(win.HGDIOBJ(f.hfont))
	}
}

func (app *applicationData) initFont() {
	const fontSize = 8
	const fontName = "Segoe UI"

	defer func() {
		if app.GUIFont == nil {
			log.Printf("error: Failed to create GUI font '%s' of size %d\n", fontName, fontSize)
			log.Println("Switching to windows default gui font...")
			app.GUIFont = &GdiFont{}
			app.GUIFont.hfont = win.HFONT(win.GetStockObject(win.DEFAULT_GUI_FONT))
		}
	}()
	app.GUIFont = CreateDPIAwareFont(fontName, fontSize)
}

func (app *applicationData) GetGuiFont() *GdiFont {
	return app.GUIFont
}

// GetCursorPos returns cursor position
func (app *applicationData) GetCursorPos() Point {
	var pt win.POINT
	win.GetCursorPos(&pt)
	return Point{X: int(pt.X), Y: int(pt.Y)}
}

// GetMainWindow returns pointer to the main window object
func (app *applicationData) GetMainWindow() Form {
	return app.mainWindow
}

func (app *applicationData) GetAppInstance() win.HINSTANCE {
	return app.hInstance
}

// Exit closes the application
func (app *applicationData) Exit() {
	win.PostQuitMessage(0)
}

func (app *applicationData) Run() {
	logInfo("Show the main window and run application message loop...")
	if app.mainWindow != nil {
		// The main application message loop
		var msg win.MSG
		for {
			if win.GetMessage(&msg, 0, 0, 0) == 0 {
				break
			}
			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg)
		}
		logInfo("Application terminated")
		// Someone closed the window....what the hell *facepalms*
		logInfo("Do cleanups...")
		app.mainWindow.Dispose()
		if app.GUIFont != nil {
			app.GUIFont.Dispose()
		}
	}
	win.GdiplusShutdown()
	ole.CoUninitialize()
}
