package reza

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	win "github.com/lxn/win"
)

const AC_SRC_OVER = 0

// Color defines RGBA color values
type Color struct {
	color.RGBA
}

// Rect is data structure that holds data for a rectangle (left, top, right, bottom)
type Rect struct {
	Left, Top, Right, Bottom int
}

// Size is a type that holds width, height data
type Size struct {
	Width, Height int
}

// Point is structure consists of two variable names x and y hold information location of..something..i guess
type Point struct {
	X, Y int
}

// Graphics is a GDI32 graphics context
type Graphics struct {
	hdc win.HDC
}

// BitmapImage is struct contains data for a image
type BitmapImage struct {
	Width, Height int
	bitmapGray    win.HBITMAP
	bitmap        win.HBITMAP
}

// AsCOLORREF returns all the color values together as type 'win.COLORREF'
func (c Color) AsCOLORREF() win.COLORREF {
	return win.RGB(c.R, c.G, c.B)
}

func (c *Color) AsRGBA() color.RGBA {
	return c.RGBA
}

// AsString returns
func (c *Color) AsString() string {
	return "(" + strconv.Itoa(int(c.R)) + ", " + strconv.Itoa(int(c.G)) + ",  " + strconv.Itoa(int(c.B)) + ")"
}

// IsEqualTo compares two colors
func (c *Color) IsEqualTo(c2 Color) bool {
	return (c.R == c2.R && c.G == c2.G && c.B == c2.B && c.A == c2.A)
}

func IsEqualColor(c Color, c2 Color) bool {
	return (c.R == c2.R && c.G == c2.G && c.B == c2.B && c.A == c2.A)
}

// FromCOLORREF converts win.COLORREF to Color
func FromCOLORREF(rgba win.COLORREF) Color {
	r, g, b, a := (rgba & 0xff), (rgba>>8)&0xff, (rgba>>16)&0xff, (rgba>>24)&0xff
	return Color{RGBA: color.RGBA{R: byte(r), G: byte(g), B: byte(b), A: byte(a)}}
}

func NewRect(left, top, right, bottom int) *Rect {
	return &Rect{Left: left, Top: top, Right: right, Bottom: bottom}
}

// AsRECT returns the rect as type 'win.RECT'
func (rc *Rect) AsRECT() win.RECT {
	return win.RECT{Left: int32(rc.Left), Top: int32(rc.Top), Right: int32(rc.Right), Bottom: int32(rc.Bottom)}
}

func (rc *Rect) AsPoints() (leftTop, rightBottom Point) {
	return Point{X: rc.Left, Y: rc.Top}, Point{X: rc.Right, Y: rc.Bottom}
}

func FromRECT(rect *win.RECT) Rect {
	return Rect{Left: int(rect.Left), Top: int(rect.Top), Right: int(rect.Right), Bottom: int(rect.Bottom)}
}

// Width returns width of the rectangle
func (rc *Rect) Width() int {
	return rc.Right - rc.Left
}

// Height returns height of the rectangle
func (rc *Rect) Height() int {
	return rc.Bottom - rc.Top
}

// CenterX returns the center x coordinate of the rectangle
func (rc *Rect) CenterX() int {
	return rc.Left + (rc.Width() / 2)
}

// CenterY returns the center y coordinate of the rectangle
func (rc *Rect) CenterY() int {
	return rc.Top + (rc.Height() / 2)
}

// Center returns the center point of the reactangle
func (rc *Rect) Center() Point {
	return Point{rc.CenterX(), rc.CenterY()}
}

const (
	RectPointTopLeft = iota
	RectPointTop
	RectPointTopRight
	RectPointLeft
	RectPointRight
	RectPointBottomLeft
	RectPointBottom
	RectPointBottomRight
)

func (rc *Rect) GetEightPoints() [8]Point {
	halfWidth := rc.Width() / 2
	halfHeight := rc.Height() / 2
	return [8]Point{
		{X: rc.Left, Y: rc.Top},
		{X: rc.Left + halfWidth, Y: rc.Top},
		{X: rc.Right, Y: rc.Top},
		{X: rc.Left, Y: rc.Top + halfHeight},
		{X: rc.Right, Y: rc.Top + halfHeight},
		{X: rc.Left, Y: rc.Bottom},
		{X: rc.Left + halfWidth, Y: rc.Bottom},
		{X: rc.Right, Y: rc.Bottom},
	}
}

func (rc *Rect) GetCornerPoints() [4]Point {
	return [4]Point{
		{X: rc.Left, Y: rc.Top},
		{X: rc.Right, Y: rc.Top},
		{X: rc.Left, Y: rc.Bottom},
		{X: rc.Right, Y: rc.Bottom},
	}
}

// Inflate inflates
func (rc *Rect) Inflate(dx, dy int) {
	rc.Left -= dx
	rc.Right += dx
	rc.Top -= dy
	rc.Bottom += dy
}

func (pt *Point) XY(from *Point) (x, y int) {
	return pt.X, pt.Y
}

func (pt *Point) Distance(from *Point) Point {
	return Point{pt.X - from.X, pt.Y - from.Y}
}

func (pt *Point) DistanceF(from *Point) float64 {
	// Calculating distance
	return math.Sqrt(math.Pow(float64(pt.X)-float64(from.X), 2) + math.Pow(float64(pt.Y)-float64(from.Y), 2))
}

func (pt *Point) DistanceI(from *Point) int {
	// Calculating distance
	return int(math.Sqrt(math.Pow(float64(pt.X)-float64(from.X), 2) + math.Pow(float64(pt.Y)-float64(from.Y), 2)))
}

// IsInsideRect determines whether the point is inside of the given rectangle or not
func (pt *Point) IsInsideRect(rect *Rect) bool {
	if pt.X >= rect.Left && pt.X <= rect.Right {
		if pt.Y >= rect.Top && pt.Y <= rect.Bottom {
			return true
		}
	}
	return false
}

// IsPointInside determines....
func (rc *Rect) IsPointInside(pt *Point) bool {
	if pt.X >= rc.Left && pt.X <= rc.Right {
		if pt.Y >= rc.Top && pt.Y <= rc.Bottom {
			return true
		}
	}
	return false
}

func Rgb(r byte, g byte, b byte) Color {
	return Color{RGBA: color.RGBA{R: byte(r), G: byte(g), B: byte(b), A: byte(255)}}
}

func NewRgb(r byte, g byte, b byte) *Color {
	return &Color{RGBA: color.RGBA{R: byte(r), G: byte(g), B: byte(b), A: byte(255)}}
}

func Rgba(r byte, g byte, b byte, a byte) Color {
	return Color{RGBA: color.RGBA{R: byte(r), G: byte(g), B: byte(b), A: byte(a)}}
}

type BitmapGraphics struct {
	HBitmap  win.HBITMAP
	Hdc      win.HDC
	Graphics *Graphics
	Data     []uint8
}

func NewBitmapGraphics(width, height int) *BitmapGraphics {
	var bi win.BITMAPV5HEADER
	const bitsPerPixel = 32 // We only work with 32 bit (BGRA)
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = int32(width)
	bi.BiHeight = -int32(height)
	bi.BiPlanes = 1
	bi.BiBitCount = uint16(bitsPerPixel)
	bi.BiCompression = win.BI_RGB

	hdc := win.GetDC(0)
	defer win.ReleaseDC(0, hdc)

	bg := &BitmapGraphics{}
	bg.Hdc = win.CreateCompatibleDC(hdc)
	bg.Graphics = NewGraphics(bg.Hdc)

	var lpBits unsafe.Pointer

	// Create the DIB section with an alpha channel.
	bg.HBitmap = win.CreateDIBSection(hdc, &bi.BITMAPINFOHEADER, win.DIB_RGB_COLORS, &lpBits, 0, 0)
	switch bg.HBitmap {
	case 0, win.ERROR_INVALID_PARAMETER:
		log.Println("CreateDIBSection failed")
		bg.Dispose()
		return nil
	}
	win.SelectObject(bg.Hdc, win.HGDIOBJ(bg.HBitmap))
	brushBack := win.GetStockObject(win.WHITE_BRUSH)
	wrect := win.RECT{
		Left:   0,
		Top:    0,
		Right:  int32(width),
		Bottom: int32(height),
	}
	FillRect(bg.Hdc, &wrect, win.HBRUSH(brushBack))

	length := ((width*bitsPerPixel + 31) / 32) * 4 * height
	// Slice memory layout
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(lpBits), length, length}
	// Use unsafe to turn sl into a []uint8
	bg.Data = *(*[]uint8)(unsafe.Pointer(&sl))
	return bg
}

func (bg *BitmapGraphics) GetBPP() int {
	return 32
}

func (bg *BitmapGraphics) Dispose() {
	if bg.HBitmap != 0 {
		win.DeleteObject(win.HGDIOBJ(bg.HBitmap))
	}
	if bg.Hdc != 0 {
		win.DeleteDC(bg.Hdc)
	}
}

func CreateHBitmap(width, height int) (win.HBITMAP, []uint8) {
	var bi win.BITMAPV5HEADER
	const bitsPerPixel = 32 // We only work with 32 bit (BGRA)
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = int32(width)
	bi.BiHeight = -int32(height)
	bi.BiPlanes = 1
	bi.BiBitCount = uint16(bitsPerPixel)
	bi.BiCompression = win.BI_RGB
	// The following mask specification specifies a supported 32 BPP
	// alpha format for Windows XP.
	bi.BV4RedMask = 0x00FF0000
	bi.BV4GreenMask = 0x0000FF00
	bi.BV4BlueMask = 0x000000FF
	bi.BV4AlphaMask = 0xFF000000

	tempdc := win.GetDC(0)
	defer win.ReleaseDC(0, tempdc)

	var lpBits unsafe.Pointer

	// Create the DIB section with an alpha channel.
	hbitmap := win.CreateDIBSection(tempdc, &bi.BITMAPINFOHEADER, win.DIB_RGB_COLORS, &lpBits, 0, 0)
	switch hbitmap {
	case 0, win.ERROR_INVALID_PARAMETER:
		log.Println("CreateDIBSection failed")
		return 0, nil
	}

	length := ((width*bitsPerPixel + 31) / 32) * 4 * height
	// Slice memory layout
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(lpBits), length, length}
	// Use unsafe to turn sl into a []uint8
	bitmapArray := *(*[]uint8)(unsafe.Pointer(&sl))
	return hbitmap, bitmapArray
}

func DeleteHBitmap(hbitmap win.HBITMAP) {
	if hbitmap != 0 {
		win.DeleteObject(win.HGDIOBJ(hbitmap))
	}
}

func HBitmapFromImage(im image.Image) (win.HBITMAP, []byte) {
	var bi win.BITMAPV5HEADER
	bi.BiSize = uint32(unsafe.Sizeof(bi))
	bi.BiWidth = int32(im.Bounds().Dx())
	bi.BiHeight = -int32(im.Bounds().Dy())
	bi.BiPlanes = 1
	bi.BiBitCount = 32
	bi.BiCompression = win.BI_RGB
	// The following mask specification specifies a supported 32 BPP
	// alpha format for Windows XP.
	bi.BV4RedMask = 0x00FF0000
	bi.BV4GreenMask = 0x0000FF00
	bi.BV4BlueMask = 0x000000FF
	bi.BV4AlphaMask = 0xFF000000

	hdc := win.GetDC(0)
	defer win.ReleaseDC(0, hdc)

	var lpBits unsafe.Pointer

	// Create the DIB section with an alpha channel.
	hBitmap := win.CreateDIBSection(hdc, &bi.BITMAPINFOHEADER, win.DIB_RGB_COLORS, &lpBits, 0, 0)
	switch hBitmap {
	case 0, win.ERROR_INVALID_PARAMETER:
		log.Println("CreateDIBSection failed")
		return 0, nil
	}

	bitsPerPixel := 32
	length := ((im.Bounds().Dx()*bitsPerPixel + 31) / 32) * 4 * im.Bounds().Dy()
	// Slice memory layout
	var sl = struct {
		addr uintptr
		len  int
		cap  int
	}{uintptr(lpBits), length, length}

	// Use unsafe to turn sl into a []byte.
	bitmapArray := *(*[]byte)(unsafe.Pointer(&sl))

	i := 0
	for y := im.Bounds().Min.Y; y != im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x != im.Bounds().Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			bitmapArray[i+3] = byte(a >> 8)
			bitmapArray[i+2] = byte(r >> 8)
			bitmapArray[i+1] = byte(g >> 8)
			bitmapArray[i+0] = byte(b >> 8)
			i += 4
		}
	}

	return hBitmap, bitmapArray
}

func CreateBitmapImage(filename string, withGrayScale bool) (*BitmapImage, error) {
	logInfo("create bitmap image...")
	catFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer catFile.Close()
	imData, _, err := image.Decode(catFile)
	if err != nil {
		fmt.Println(err)
	}
	bounds := imData.Bounds()

	bitmapImage := &BitmapImage{Width: int(bounds.Size().X), Height: int(bounds.Size().Y)}
	bitmapImage.bitmap, _ = HBitmapFromImage(imData)
	if withGrayScale {
		// Converting image to grayscale
		grayImg := image.NewRGBA(imData.Bounds())
		for y := imData.Bounds().Min.Y; y < imData.Bounds().Max.Y; y++ {
			for x := imData.Bounds().Min.X; x < imData.Bounds().Max.X; x++ {
				R, G, B, oA := imData.At(x, y).RGBA()
				Y := (0.2126*float64(R) + 0.7152*float64(G) + 0.0752*float64(B)) * (255.0 / 65535)
				grayPix := color.Gray{uint8(Y)}
				grayR, grayG, grayB, _ := grayPix.RGBA()
				// Keep the original alpha
				grayAlphaPix := color.RGBA{R: uint8(grayR), G: uint8(grayG), B: uint8(grayB), A: uint8(oA)}
				grayImg.Set(x, y, grayAlphaPix)
			}
		}
		bitmapImage.bitmapGray, _ = HBitmapFromImage(grayImg)
	}
	return bitmapImage, nil
}

func (img *BitmapImage) Dispose() {
	if img.bitmapGray != 0 {
		win.DeleteObject(win.HGDIOBJ(img.bitmapGray))
	}
	if img.bitmap != 0 {
		win.DeleteObject(win.HGDIOBJ(img.bitmap))
	}
	img.bitmap = 0
	img.bitmapGray = 0
}

type GdiObject interface {
	GetGdiObject() win.HGDIOBJ
}

type Brush struct {
	GdiObject
	hbrush win.HBRUSH
}

func NewSolidBrush(color *Color) *Brush {
	gdiBrush := &Brush{}
	lb := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: color.AsCOLORREF()}
	gdiBrush.hbrush = win.CreateBrushIndirect(lb)
	return gdiBrush
}

func (b *Brush) Dispose() {
	win.DeleteObject(win.HGDIOBJ(b.hbrush))
}

func (b *Brush) GetGdiObject() win.HGDIOBJ {
	return win.HGDIOBJ(b.hbrush)
}

type Pen struct {
	GdiObject
	NativePen win.HPEN
}

func NewPen(style int, width int, color *Color) *Pen {
	gdiPen := &Pen{}
	gdiPen.NativePen = CreatePen(int32(style), uint32(width), color.AsCOLORREF())
	return gdiPen
}

func NewSolidPen(width int, color *Color) *Pen {
	gdiPen := &Pen{}
	gdiPen.NativePen = CreatePen(win.PS_SOLID, uint32(width), color.AsCOLORREF())
	return gdiPen
}

func NewDashPen(width int, color *Color) *Pen {
	gdiPen := &Pen{}
	pbrush := &win.LOGBRUSH{LbStyle: win.PS_SOLID, LbColor: color.AsCOLORREF()}
	gdiPen.NativePen = win.ExtCreatePen(win.PS_COSMETIC|win.PS_ALTERNATE, uint32(width), pbrush, 0, nil)
	return gdiPen
}

func NewUserStylePen(width int, color *Color, userstyle []uint32) *Pen {
	gdiPen := &Pen{}
	pbrush := &win.LOGBRUSH{LbStyle: win.PS_SOLID, LbColor: color.AsCOLORREF()}
	gdiPen.NativePen = win.ExtCreatePen(win.PS_GEOMETRIC|win.PS_USERSTYLE, 1, pbrush, uint32(len(userstyle)), &userstyle[0])
	return gdiPen
}

func (b *Pen) Dispose() {
	win.DeleteObject(win.HGDIOBJ(b.NativePen))
}

func (pen *Pen) GetGdiObject() win.HGDIOBJ {
	return win.HGDIOBJ(pen.NativePen)
}

func NewGraphics(hdc win.HDC) *Graphics {
	return &Graphics{hdc: hdc}
}

func (g *Graphics) Dispose() {

}

func (g *Graphics) GetHDC() win.HDC {
	return g.hdc
}

func (g *Graphics) BitBlt(x, y, width, height int, hdcSrc win.HDC, xSrc, ySrc int, op uint) {
	win.BitBlt(g.GetHDC(),
		int32(x), int32(y),
		int32(width), int32(height),
		hdcSrc,
		int32(xSrc), int32(ySrc), uint32(op))
}

func (g *Graphics) AlphaBlend(x, y, width, height int, hdcSrc win.HDC, xSrc, ySrc, widthSrc, heightSrc int) {
	var bf win.BLENDFUNCTION
	bf.BlendOp = AC_SRC_OVER
	bf.BlendFlags = 0
	bf.SourceConstantAlpha = 255
	bf.AlphaFormat = win.AC_SRC_ALPHA
	win.AlphaBlend(g.GetHDC(), int32(x), int32(y), int32(width), int32(height), hdcSrc, int32(xSrc), int32(ySrc), int32(widthSrc), int32(heightSrc), bf)
}

// DrawBitmapImage draws image
func (g *Graphics) DrawBitmapImage(image *BitmapImage, x, y int, gray bool) {
	var bf win.BLENDFUNCTION
	bf.BlendOp = AC_SRC_OVER
	bf.BlendFlags = 0
	bf.SourceConstantAlpha = 255
	bf.AlphaFormat = win.AC_SRC_ALPHA

	bitmap := image.bitmap
	if gray {
		bitmap = image.bitmapGray
		if bitmap == 0 {
			bitmap = image.bitmap
		}
		bf.SourceConstantAlpha = 180
	}

	memDC := win.CreateCompatibleDC(0)
	prevObj := win.SelectObject(memDC, win.HGDIOBJ(bitmap))

	//win.BitBlt(hdc, 0, 0, int32(canvas.image.width), int32(canvas.image.height), memDC, 0, 0, win.SRCCOPY)

	win.AlphaBlend(g.hdc, int32(x), int32(y), int32(image.Width), int32(image.Height), memDC,
		0, 0, int32(image.Width), int32(image.Height), bf)

	win.SelectObject(memDC, prevObj)
	win.DeleteDC(memDC)
}

// DrawBitmapImageCenter draws the image at the center of the given rectangle
func (g *Graphics) DrawBitmapImageCenter(image *BitmapImage, rect *Rect, gray bool) {
	center := rect.Center()
	centeredLeft := center.X - (image.Width / 2)
	centeredTop := center.Y - (image.Height / 2)
	g.DrawBitmapImage(image, centeredLeft, centeredTop, gray)
}

// MeasureText destermines the size of the given text according to the current font set to this graphics context
func (g *Graphics) MeasureText(text string, format uint32, font win.HFONT) *Rect {
	previousFont := win.SelectObject(g.hdc, win.HGDIOBJ(font))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousFont))
	var rectText win.RECT
	utf16, _ := syscall.UTF16PtrFromString(text)
	win.DrawTextEx(g.hdc, utf16,
		int32(len(text)), &rectText, win.DT_CALCRECT|win.DT_EXPANDTABS|format, nil)
	return &Rect{Left: int(rectText.Left), Top: int(rectText.Top), Right: int(rectText.Right), Bottom: int(rectText.Bottom)}
}

// DrawText draws text
func (g *Graphics) DrawText(text string, rect *Rect, format uint32, textColor *Color, font win.HFONT) {
	previousBkMode := win.SetBkMode(g.hdc, win.TRANSPARENT)
	defer win.SetBkMode(g.hdc, previousBkMode)
	if font != 0 {
		previousFont := win.SelectObject(g.hdc, win.HGDIOBJ(font))
		defer win.SelectObject(g.hdc, win.HGDIOBJ(previousFont))
	}
	previousTextColor := win.SetTextColor(g.hdc, textColor.AsCOLORREF())
	defer win.SetTextColor(g.hdc, previousTextColor)
	wrect := rect.AsRECT()
	utf16, _ := syscall.UTF16FromString(text)
	win.DrawTextEx(g.hdc, &utf16[0], int32(len(utf16)), &wrect, format, nil)
}

// DrawDropDownArrow draws a simple drop down arrow
func (g *Graphics) DrawDropDownArrow(x, y int, color *Color) {
	lineX := x
	lineY := y
	g.DrawLine(lineX-2, lineY, lineX+3, lineY, color)
	lineY++
	g.DrawLine(lineX-1, lineY, lineX+2, lineY, color)
	lineY++
	g.DrawLine(lineX, lineY, lineX+1, lineY, color)
}

func (g *Graphics) DrawCheckerBoard(rect *Rect, size int, color1, color2 *Color) {
	checkSize := size
	xCount := rect.Width() / checkSize
	yCount := rect.Height() / checkSize
	pen1 := NewSolidPen(1, color1)
	defer pen1.Dispose()
	pen2 := NewSolidPen(1, color2)
	defer pen2.Dispose()
	brush1 := NewSolidBrush(color1)
	defer brush1.Dispose()
	brush2 := NewSolidBrush(color2)
	defer brush2.Dispose()
	for x := 0; x <= xCount; x++ {
		for y := 0; y <= yCount; y++ {
			cx := x * checkSize
			cy := y * checkSize
			rectCheck := Rect{
				Left:   cx,
				Top:    cy,
				Right:  cx + checkSize,
				Bottom: cy + checkSize,
			}
			if (x+y)%2 == 0 {
				g.DrawFillRectangleEx(&rectCheck, pen2, brush2)
			} else {
				g.DrawFillRectangleEx(&rectCheck, pen1, brush1)
			}
		}
	}
}

// DrawLine draws line
func (g *Graphics) DrawLine(x, y, x2, y2 int, color *Color) {
	pen := NewSolidPen(1, color)
	defer pen.Dispose()
	g.SelectObject(pen)
	win.MoveToEx(g.hdc, x, y, nil)
	win.LineTo(g.hdc, int32(x2), int32(y2))
}

func (g *Graphics) DrawLineEx(x, y, x2, y2 int, pen *Pen) {
	if pen == nil {
		win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_PEN)))
	} else {
		win.SelectObject(g.hdc, win.HGDIOBJ(pen.NativePen))
	}
	win.MoveToEx(g.hdc, x, y, nil)
	win.LineTo(g.hdc, int32(x2), int32(y2))
}

func (g *Graphics) DrawLineOnly(x, y, x2, y2 int) {
	win.MoveToEx(g.hdc, x, y, nil)
	win.LineTo(g.hdc, int32(x2), int32(y2))
}

// FillRect fills the given rectangle with the specified color
func (g *Graphics) FillRect(rc *Rect, color *Color) {
	brush := NewSolidBrush(color)
	defer brush.Dispose()
	wrect := rc.AsRECT()
	FillRect(g.GetHDC(), &wrect, brush.hbrush)
}

// DrawFillRectangle draws a rectangle filled and bordered with the given colors
func (g *Graphics) DrawFillRectangleEx(rc *Rect, pen *Pen, brush *Brush) {
	win.SelectObject(g.hdc, win.HGDIOBJ(pen.NativePen))
	win.SelectObject(g.hdc, win.HGDIOBJ(brush.hbrush))
	win.Rectangle_(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom))
}

// DrawFillRectangle draws a rectangle filled and bordered with the given colors
func (g *Graphics) DrawFillRectangle(rc *Rect, borderColor, fillColor *Color) {
	if borderColor == fillColor {
		g.FillRect(rc, fillColor)
		return
	}
	pbrush := &win.LOGBRUSH{LbStyle: win.PS_SOLID, LbColor: borderColor.AsCOLORREF()}
	pen := win.ExtCreatePen(win.PS_SOLID, 1, pbrush, 0, nil)
	defer win.DeleteObject(win.HGDIOBJ(pen))
	previousPen := win.SelectObject(g.hdc, win.HGDIOBJ(pen))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousPen))

	lb := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: fillColor.AsCOLORREF()}
	brush := win.CreateBrushIndirect(lb)
	defer win.DeleteObject(win.HGDIOBJ(brush))
	previousBrush := win.SelectObject(g.hdc, win.HGDIOBJ(brush))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousBrush))

	win.Rectangle_(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom))
}

// DrawRectangle draws a rectangle bordered with the given color
func (g *Graphics) DrawRectangle(rc *Rect, color *Color) {
	pbrush := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: color.AsCOLORREF()}
	pen := win.ExtCreatePen(win.PS_SOLID, 1, pbrush, 0, nil)
	defer win.DeleteObject(win.HGDIOBJ(pen))
	previousPen := win.SelectObject(g.hdc, win.HGDIOBJ(pen))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousPen))
	previousBrush := win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_BRUSH)))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousBrush))
	win.Rectangle_(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom))
}

func (g *Graphics) DrawRoundRectangle(rc *Rect, cornerWidth, cornerHeight int, color *Color) {
	pbrush := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: color.AsCOLORREF()}
	pen := win.ExtCreatePen(win.PS_SOLID, 1, pbrush, 0, nil)
	defer win.DeleteObject(win.HGDIOBJ(pen))
	previousPen := win.SelectObject(g.hdc, win.HGDIOBJ(pen))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousPen))
	previousBrush := win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_BRUSH)))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousBrush))
	win.RoundRect(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom), int32(cornerWidth), int32(cornerHeight))
}

func (g *Graphics) DrawDashedRectangle(rc *Rect, color *Color) {
	pbrush := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: color.AsCOLORREF()}
	pen := win.ExtCreatePen(win.PS_COSMETIC|win.PS_ALTERNATE, 1, pbrush, 0, nil)
	defer win.DeleteObject(win.HGDIOBJ(pen))
	previousPen := win.SelectObject(g.hdc, win.HGDIOBJ(pen))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousPen))

	previousBrush := win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_BRUSH)))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousBrush))

	win.Rectangle_(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom))
}

func (g *Graphics) DrawEllipse(rc *Rect, color *Color) {
	pbrush := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: color.AsCOLORREF()}
	pen := win.ExtCreatePen(win.PS_SOLID, 1, pbrush, 0, nil)
	defer win.DeleteObject(win.HGDIOBJ(pen))
	previousPen := win.SelectObject(g.hdc, win.HGDIOBJ(pen))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousPen))
	previousBrush := win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_BRUSH)))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousBrush))
	win.Ellipse(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom))
}

func (g *Graphics) DrawPolygon(points []Point, color *Color) {
	pbrush := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: color.AsCOLORREF()}
	pen := win.ExtCreatePen(win.PS_SOLID, 1, pbrush, 0, nil)
	defer win.DeleteObject(win.HGDIOBJ(pen))
	previousPen := win.SelectObject(g.hdc, win.HGDIOBJ(pen))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousPen))
	previousBrush := win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_BRUSH)))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousBrush))
	wpoints := make([]win.POINT, len(points))
	for i := range points {
		wpoints[i].X = int32(points[i].X)
		wpoints[i].Y = int32(points[i].Y)
	}
	Polygon(g.hdc, &wpoints[0], int32(len(wpoints)))
}

func (g *Graphics) SelectObject(obj GdiObject) (prev win.HGDIOBJ) {
	prev = win.SelectObject(g.hdc, obj.GetGdiObject())
	return
}

func (g *Graphics) SelectBrushAndPen(pen *Pen, brush *Brush) {
	if pen == nil {
		win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_PEN)))
	} else {
		g.SelectObject(pen)
	}
	if brush == nil {
		win.SelectObject(g.hdc, win.HGDIOBJ(win.GetStockObject(win.NULL_BRUSH)))
	} else {
		g.SelectObject(brush)
	}
}

func (g *Graphics) DrawPolygonEx(points []Point, pen *Pen, brush *Brush) {
	g.SelectBrushAndPen(pen, brush)
	wpoints := make([]win.POINT, len(points))
	for i := range points {
		wpoints[i].X = int32(points[i].X)
		wpoints[i].Y = int32(points[i].Y)
	}
	Polygon(g.hdc, &wpoints[0], int32(len(wpoints)))
}

func (g *Graphics) DrawEllipseEx(rc *Rect, pen *Pen, brush *Brush) {
	g.SelectBrushAndPen(pen, brush)
	win.Ellipse(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom))
}

func (g *Graphics) DrawRectangleEx(rc *Rect, pen *Pen, brush *Brush) {
	g.SelectBrushAndPen(pen, brush)
	win.Rectangle_(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom))
}

func (g *Graphics) DrawRoundRectangleEx(rc *Rect, w, h int, pen *Pen, brush *Brush) {
	g.SelectBrushAndPen(pen, brush)
	win.RoundRect(g.hdc, int32(rc.Left), int32(rc.Top), int32(rc.Right), int32(rc.Bottom), int32(w), int32(h))
}

func (g *Graphics) DrawFillPolygon(points []Point, color, fillColor *Color) {
	pbrush := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: color.AsCOLORREF()}
	pen := win.ExtCreatePen(win.PS_SOLID, 1, pbrush, 0, nil)
	defer win.DeleteObject(win.HGDIOBJ(pen))
	previousPen := win.SelectObject(g.hdc, win.HGDIOBJ(pen))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousPen))

	lb := &win.LOGBRUSH{LbStyle: win.BS_SOLID, LbColor: fillColor.AsCOLORREF()}
	brush := win.CreateBrushIndirect(lb)
	defer win.DeleteObject(win.HGDIOBJ(brush))
	previousBrush := win.SelectObject(g.hdc, win.HGDIOBJ(brush))
	defer win.SelectObject(g.hdc, win.HGDIOBJ(previousBrush))

	wpoints := make([]win.POINT, len(points))
	for i := range points {
		wpoints[i].X = int32(points[i].X)
		wpoints[i].Y = int32(points[i].Y)
	}
	Polygon(g.hdc, &wpoints[0], int32(len(wpoints)))
}
