package reza

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	win "github.com/lxn/win"
	"golang.org/x/sys/windows"
)

//---------------------------------------------------------

type MARGINS struct {
	CxLeftWidth    int32 // width of left border that retains its size
	CxRightWidth   int32 // width of right border that retains its size
	CyTopHeight    int32 // height of top border that retains its size
	CyBottomHeight int32 // height of bottom border that retains its size
}

var (
	getCapture                 *windows.LazyProc
	mapWindowPoints            *windows.LazyProc
	setLayeredWindowAttributes *windows.LazyProc
	fillRect                   *windows.LazyProc
	loadCursorFromFileW        *windows.LazyProc
)

var (
	dwmIsCompositionEnabled      *windows.LazyProc
	dwmExtendFrameIntoClientArea *windows.LazyProc
	dwmDefWindowProc             *windows.LazyProc
)

const (
	SM_CXPADDEDBORDER = 92
)

// 'WINDOW' class parts
const (
	WP_CAPTION                        = 1
	WP_SMALLCAPTION                   = 2
	WP_MINCAPTION                     = 3
	WP_SMALLMINCAPTION                = 4
	WP_MAXCAPTION                     = 5
	WP_SMALLMAXCAPTION                = 6
	WP_FRAMELEFT                      = 7
	WP_FRAMERIGHT                     = 8
	WP_FRAMEBOTTOM                    = 9
	WP_SMALLFRAMELEFT                 = 10
	WP_SMALLFRAMERIGHT                = 11
	WP_SMALLFRAMEBOTTOM               = 12
	WP_SYSBUTTON                      = 13
	WP_MDISYSBUTTON                   = 14
	WP_MINBUTTON                      = 15
	WP_MDIMINBUTTON                   = 16
	WP_MAXBUTTON                      = 17
	WP_CLOSEBUTTON                    = 18
	WP_SMALLCLOSEBUTTON               = 19
	WP_MDICLOSEBUTTON                 = 20
	WP_RESTOREBUTTON                  = 21
	WP_MDIRESTOREBUTTON               = 22
	WP_HELPBUTTON                     = 23
	WP_MDIHELPBUTTON                  = 24
	WP_HORZSCROLL                     = 25
	WP_HORZTHUMB                      = 26
	WP_VERTSCROLL                     = 27
	WP_VERTTHUMB                      = 28
	WP_DIALOG                         = 29
	WP_CAPTIONSIZINGTEMPLATE          = 30
	WP_SMALLCAPTIONSIZINGTEMPLATE     = 31
	WP_FRAMELEFTSIZINGTEMPLATE        = 32
	WP_SMALLFRAMELEFTSIZINGTEMPLATE   = 33
	WP_FRAMERIGHTSIZINGTEMPLATE       = 34
	WP_SMALLFRAMERIGHTSIZINGTEMPLATE  = 35
	WP_FRAMEBOTTOMSIZINGTEMPLATE      = 36
	WP_SMALLFRAMEBOTTOMSIZINGTEMPLATE = 37
	WP_FRAME                          = 38
)

// 'WINDOW' class states
const (
	FS_ACTIVE   = 1
	FS_INACTIVE = 2
)

const (
	CS_ACTIVE   = 1
	CS_INACTIVE = 2
	CS_DISABLED = 3
)

const (
	MXCS_ACTIVE   = 1
	MXCS_INACTIVE = 2
	MXCS_DISABLED = 3
)

const (
	MNCS_ACTIVE   = 1
	MNCS_INACTIVE = 2
	MNCS_DISABLED = 3
)

const (
	HSS_NORMAL   = 1
	HSS_HOT      = 2
	HSS_PUSHED   = 3
	HSS_DISABLED = 4
)

const (
	HTS_NORMAL   = 1
	HTS_HOT      = 2
	HTS_PUSHED   = 3
	HTS_DISABLED = 4
)

const (
	VSS_NORMAL   = 1
	VSS_HOT      = 2
	VSS_PUSHED   = 3
	VSS_DISABLED = 4
)

const (
	VTS_NORMAL   = 1
	VTS_HOT      = 2
	VTS_PUSHED   = 3
	VTS_DISABLED = 4
)

const (
	SBS_NORMAL   = 1
	SBS_HOT      = 2
	SBS_PUSHED   = 3
	SBS_DISABLED = 4
)

const (
	MINBS_NORMAL   = 1
	MINBS_HOT      = 2
	MINBS_PUSHED   = 3
	MINBS_DISABLED = 4
)

const (
	RBS_NORMAL   = 1
	RBS_HOT      = 2
	RBS_PUSHED   = 3
	RBS_DISABLED = 4
)

const (
	HBS_NORMAL   = 1
	HBS_HOT      = 2
	HBS_PUSHED   = 3
	HBS_DISABLED = 4
)

const (
	CBS_NORMAL   = 1
	CBS_HOT      = 2
	CBS_PUSHED   = 3
	CBS_DISABLED = 4
)

var (
	getThemeSysSize *windows.LazyProc
	getThemeRect    *windows.LazyProc
)

var (
	coCreateInstance *windows.LazyProc
)

var (
	qiSearch *windows.LazyProc
)

/*
const (
	COINITBASE_MULTITHREADED = 0x0
	COINIT_APARTMENTTHREADED = 0x2 // Apartment model
	COINIT_MULTITHREADED     = COINITBASE_MULTITHREADED
	COINIT_DISABLE_OLE1DDE   = 0x4 // Don't use DDE for Ole1 support.
	COINIT_SPEED_OVER_MEMORY = 0x8 // Trade memory for speed.
)
*/
const (
	S_OK    = win.HRESULT(0)
	S_FALSE = win.HRESULT(1)
	//RPC_E_CHANGED_MODE = win.HRESULT(0x80010106)
)

type QITAB struct {
	Piid     *ole.GUID
	DwOffset uint32
}

const LF_FACESIZE = 32
const LF_FULLFACESIZE = 64

type ENUMLOGFONTEX struct {
	ElfLogFont  win.LOGFONT
	ElfFullName [LF_FULLFACESIZE]uint16
	ElfStyle    [LF_FACESIZE]uint16
	ElfScript   [LF_FACESIZE]uint16
}

type TEXTMETRIC struct {
	TmHeight           int32
	TmAscent           int32
	TmDescent          int32
	TmInternalLeading  int32
	TmExternalLeading  int32
	TmAveCharWidth     int32
	TmMaxCharWidth     int32
	TmWeight           int32
	TmOverhang         int32
	TmDigitizedAspectX int32
	TmDigitizedAspectY int32
	TmFirstChar        uint16
	TmLastChar         uint16
	TmDefaultChar      uint16
	TmBreakChar        uint16
	TmItalic           uint8
	TmUnderlined       uint8
	TmStruckOut        uint8
	TmPitchAndFamily   uint8
	TmCharSet          uint8
}

type NEWTEXTMETRIC struct {
	TmHeight           int32
	TmAscent           int32
	TmDescent          int32
	TmInternalLeading  int32
	TmExternalLeading  int32
	TmAveCharWidth     int32
	TmMaxCharWidth     int32
	TmWeight           int32
	TmOverhang         int32
	TmDigitizedAspectX int32
	TmDigitizedAspectY int32
	TmFirstChar        uint16
	TmLastChar         uint16
	TmDefaultChar      uint16
	TmBreakChar        uint16
	TmItalic           uint8
	TmUnderlined       uint8
	TmStruckOut        uint8
	TmPitchAndFamily   uint8
	TmCharSet          uint8
	NtmFlags           int32
	NtmSizeEM          uint32
	NtmCellHeight      uint32
	NtmAvgWidth        uint32
}

type FONTSIGNATURE struct {
	FsUsb [4]uint32
	FsCsb [2]uint32
}

type NEWTEXTMETRICEX struct {
	NtmTm      NEWTEXTMETRIC
	NtmFontSig FONTSIGNATURE
}

var (
	polygon            *windows.LazyProc
	createPen          *windows.LazyProc
	enumFontFamiliesEx *windows.LazyProc
	extTextOutW        *windows.LazyProc
	createRectRgn      *windows.LazyProc
	selectClipRgn      *windows.LazyProc
)

func loadApiFunctions() {
	// Library
	libgdi32 := windows.NewLazySystemDLL("gdi32.dll")
	libdwmapi := windows.NewLazySystemDLL("dwmapi.dll")
	libuser32 := windows.NewLazySystemDLL("user32.dll")
	libuxtheme := windows.NewLazySystemDLL("uxtheme.dll")
	libole32 := windows.NewLazySystemDLL("ole32.dll")
	libshlwapi := windows.NewLazySystemDLL("shlwapi.dll")
	if libdwmapi == nil {
		log.Panicln("'dwmapi.dll' Not Found!!!")
	}
	// Gdi32 functions
	polygon = libgdi32.NewProc("Polygon")
	createPen = libgdi32.NewProc("CreatePen")
	enumFontFamiliesEx = libgdi32.NewProc("EnumFontFamiliesExW")
	extTextOutW = libgdi32.NewProc("ExtTextOutW")
	createRectRgn = libgdi32.NewProc("CreateRectRgn")
	selectClipRgn = libgdi32.NewProc("SelectClipRgn")
	// User32 Functions
	fillRect = libuser32.NewProc("FillRect")
	getCapture = libuser32.NewProc("GetCapture")
	mapWindowPoints = libuser32.NewProc("MapWindowPoints")
	setLayeredWindowAttributes = libuser32.NewProc("SetLayeredWindowAttributes")
	loadCursorFromFileW = libuser32.NewProc("LoadCursorFromFileW")
	// DWM API Functions
	dwmIsCompositionEnabled = libdwmapi.NewProc("DwmIsCompositionEnabled")
	dwmExtendFrameIntoClientArea = libdwmapi.NewProc("DwmExtendFrameIntoClientArea")
	dwmDefWindowProc = libdwmapi.NewProc("DwmDefWindowProc")
	// Theme functions
	getThemeSysSize = libuxtheme.NewProc("GetThemeSysSize")
	getThemeRect = libuxtheme.NewProc("GetThemeRect")
	// Ole32 Functions
	coCreateInstance = libole32.NewProc("CoCreateInstance")
	//coUninitialize = libole32.NewProc("CoUninitialize")
	// Shell API
	qiSearch = libshlwapi.NewProc("QISearch")
}

func CreateRectRgn(x1, y1, x2, y2 int32) win.HRGN {
	ret, _, _ := createRectRgn.Call(
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return win.HRGN(ret)
}

func SelectClipRgn(hdc win.HDC, hrgn win.HRGN) int32 {
	ret, _, _ := selectClipRgn.Call(
		uintptr(hdc),
		uintptr(hrgn))
	return int32(ret)
}

func ExtTextOutW(hdc win.HDC, x, y int32, options uint32, rect *Rect, text *uint16, textCount uint32, lpDx *int32) int32 {
	ret, _, _ := extTextOutW.Call(
		uintptr(hdc),
		uintptr(x),
		uintptr(y),
		uintptr(options),
		uintptr(unsafe.Pointer(rect)),
		uintptr(unsafe.Pointer(text)),
		uintptr(textCount),
		uintptr(unsafe.Pointer(lpDx)))
	return int32(ret)
}

func EnumFontFamiliesEx(hdc win.HDC, lpLogfont *win.LOGFONT, lpProc uintptr, lParam uintptr, dwFlags uint32) int32 {
	ret, _, _ := enumFontFamiliesEx.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(lpLogfont)),
		lpProc,
		lParam,
		uintptr(dwFlags))
	return int32(ret)
}

func Polygon(hdc win.HDC, points *win.POINT, count int32) int32 {
	ret, _, _ := polygon.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return int32(ret)
}

func CreatePen(style int32, width uint32, color win.COLORREF) win.HPEN {
	ret, _, _ := createPen.Call(
		uintptr(style),
		uintptr(width),
		uintptr(color))
	return win.HPEN(ret)
}

func QISearch(that unsafe.Pointer, pqit *QITAB, riid *ole.GUID, ppv *unsafe.Pointer) int32 {
	ret, _, _ := qiSearch.Call(
		uintptr(that),
		uintptr(unsafe.Pointer(pqit)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppv)))
	return int32(ret)
}

func CoCreateInstance(clsid *ole.GUID, outer *ole.IUnknown, clsContext uint32, iid *ole.GUID, object *uintptr) uint32 {
	r1, _, _ := coCreateInstance.Call(
		uintptr(unsafe.Pointer(clsid)),
		uintptr(unsafe.Pointer(outer)),
		uintptr(clsContext),
		uintptr(unsafe.Pointer(iid)),
		uintptr(unsafe.Pointer(object)))
	return uint32(r1)
}

/*
func CoUninitialize() win.HRESULT {
	ret, _, _ := coUninitialize.Call()
	return win.HRESULT(ret)
}
*/
func GetThemeSysSize(hTheme win.HTHEME, iSizeId int32) int32 {
	ret, _, _ := getThemeSysSize.Call(uintptr(hTheme),
		uintptr(iSizeId))
	return int32(ret)
}

func GetThemeRect(hTheme win.HTHEME, iPartId, iStateId, iPropId int32, pRect *win.RECT) int32 {
	ret, _, _ := getThemeRect.Call(
		uintptr(hTheme),
		uintptr(iPartId),
		uintptr(iStateId),
		uintptr(iPropId),
		uintptr(unsafe.Pointer(pRect)))
	return int32(ret)
}

func DwmExtendFrameIntoClientArea(hwnd win.HWND, pMarInset *MARGINS) int32 {
	ret, _, _ := dwmExtendFrameIntoClientArea.Call(uintptr(hwnd), uintptr(unsafe.Pointer(pMarInset)))
	return int32(ret)
}

func DwmDefWindowProc(hWnd win.HWND, msg uint32, wParam, lParam uintptr, lpResult *uintptr) int32 {
	ret, _, _ := dwmDefWindowProc.Call(uintptr(hWnd),
		uintptr(msg),
		wParam,
		lParam,
		uintptr(unsafe.Pointer(lpResult)))
	return int32(ret)
}

func DwmIsCompositionEnabled(enabled *bool) int32 {
	var pEnabled int32
	ret, _, _ := dwmIsCompositionEnabled.Call(uintptr(unsafe.Pointer(&pEnabled)))
	if pEnabled == 1 {
		*enabled = true
	} else {
		*enabled = false
	}
	return int32(ret)
}

func FillRect(hdc win.HDC, prect *win.RECT, hbrush win.HBRUSH) int32 {
	ret, _, _ := fillRect.Call(uintptr(hdc),
		uintptr(unsafe.Pointer(prect)),
		uintptr(hbrush))
	return int32(ret)
}

func LoadCursorFromFile(lpFileName *uint16) win.HCURSOR {
	ret, _, _ := loadCursorFromFileW.Call(uintptr(unsafe.Pointer(lpFileName)))
	return win.HCURSOR(ret)
}

func GetCapture() win.HWND {
	ret, _, _ := syscall.Syscall(getCapture.Addr(), 0,
		0,
		0,
		0)
	return win.HWND(ret)
}

func MapWindowPoints(hWndFrom, hWndTo win.HWND, lpPoints *win.POINT, cPoints uint32) int32 {
	ret, _, _ := syscall.Syscall6(mapWindowPoints.Addr(), 4,
		uintptr(hWndFrom),
		uintptr(hWndTo),
		uintptr(unsafe.Pointer(lpPoints)),
		uintptr(cPoints),
		0,
		0)

	return int32(ret)
}

const (
	LWA_ALPHA    = 0x00000002
	LWA_COLORKEY = 0x00000001
)

func SetLayeredWindowAttributes(hwnd win.HWND, crKey win.COLORREF, bAlpha uint8, dwFlags uint32) int32 {
	ret, _, _ := setLayeredWindowAttributes.Call(
		uintptr(hwnd),
		uintptr(crKey),
		uintptr(bAlpha),
		uintptr(dwFlags))
	return int32(ret)
}
