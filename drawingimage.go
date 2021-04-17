package main

import (
	"gopaint/gdiplus"
	. "gopaint/reza"
	"image"
	"image/color"
	"strconv"

	"github.com/fogleman/gg"
	win "github.com/lxn/win"
)

type DrawingImage struct {
	BGRA
	hbitmap    win.HBITMAP
	memdc      win.HDC
	context    *gdiplus.Graphics
	context2   *gg.Context
	context3   *Graphics
	filepath   string
	sizeOnDisk int64 // in bytes
	lastSaved  string
}

func NewDrawingImage(width, height int) *DrawingImage {
	this := &DrawingImage{}
	this.Rect = image.Rect(0, 0, width, height)
	this.Stride = 4 * width // 4 bytes per pixel
	this.Pix = make([]uint8, this.Stride*height)

	this.hbitmap, this.Pix = CreateHBitmap(width, height)
	hScreenDC := win.GetDC(0)
	this.memdc = win.CreateCompatibleDC(hScreenDC)
	win.SelectObject(this.memdc, win.HGDIOBJ(this.hbitmap))
	win.ReleaseDC(0, hScreenDC)

	this.context = gdiplus.NewGraphicsFromHDC(this.memdc)
	this.context.SetTextRenderingHint(gdiplus.TextRenderingHintAntiAlias)
	//this.context.SetSmoothingMode(gdiplus.SmoothingModeAntiAlias)

	rgba := (*image.RGBA)(&this.BGRA)
	this.context2 = gg.NewContextForRGBA(rgba)

	this.context3 = NewGraphics(this.memdc)

	//wrect := NewRect(0, 0, width, height).AsRECT()
	//FillRect(this.memdc, &wrect, win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH)))

	//this.bitmap = gdiplus.NewBitmapEx(int32(width), int32(height), int32(this.Stride), gdiplus.PixelFormat32bppPARGB, &this.Pix[0])

	//this.bitmap = gdiplus.NewBitmapFromHBITMAP(this.hbitmap)

	return this
}

func (image *DrawingImage) Dispose() {
	if image.context != nil {
		image.context.Dispose()
	}
	if image.memdc != 0 {
		win.DeleteDC(image.memdc)
	}
	DeleteHBitmap(image.hbitmap)
	image.hbitmap = 0

}

func (di *DrawingImage) Width() int {
	return di.Bounds().Dx()
}

func (di *DrawingImage) Height() int {
	return di.Bounds().Dy()
}

func (di *DrawingImage) GetColorAt(x, y int) Color {
	r, g, b, a := di.At(x, y).RGBA()
	return Color{RGBA: color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}}
}

func (image *DrawingImage) HasFilePath() bool {
	return len(image.filepath) > 0
}

func (image *DrawingImage) Clear(color *gdiplus.Color) {
	image.context.Clear(color)
	/*
		//i := 0
		//bitmapArray := di.Pix
		for y := di.Bounds().Min.Y; y != di.Bounds().Max.Y; y++ {
			for x := di.Bounds().Min.X; x != di.Bounds().Max.X; x++ {
				di.Set(x, y, c.AsRGBA())

					//bitmapArray[i+0] = 255
					//bitmapArray[i+1] = 0
					//bitmapArray[i+2] = 0
					//bitmapArray[i+3] = 255
					//i += 4
			}
		}
	*/
}

/*

func (image *DrawingImage) AquireData(im image.Image) {
	i := 0
	bitmapArray := image.Pix
	for y := im.Bounds().Min.Y; y != im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x != im.Bounds().Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			bitmapArray[i+0] = byte(b >> 8)
			bitmapArray[i+1] = byte(g >> 8)
			bitmapArray[i+2] = byte(r >> 8)
			bitmapArray[i+3] = byte(a >> 8)
			i += 4
		}
	}
}
*/

func (image *DrawingImage) SizeOnDisk() (asString string, available bool) {
	if image.sizeOnDisk > 0 {
		if image.sizeOnDisk < 1024 {
			return strconv.FormatInt(image.sizeOnDisk, 10) + " Bytes", true
		}
		kb := float64(image.sizeOnDisk) / 1024.0
		if kb > 1024 {
			mb := kb / 1024.0
			return strconv.FormatFloat(mb, 'f', 1, 64) + "MB", true
		} else {
			return strconv.FormatFloat(kb, 'f', 1, 64) + "KB", true
		}
	}
	return "", false
}

func (image *DrawingImage) LastSaved() (asString string, available bool) {
	if len(image.lastSaved) > 0 {
		return image.lastSaved, true
	}
	return "", false
}

// DrawImage draws image
func DrawImage(g *Graphics, image *DrawingImage, x, y int) {
	/*
		width := 100
		height := 100

		var bi win.BITMAPV5HEADER
		bi.BiSize = uint32(unsafe.Sizeof(bi))
		bi.BiWidth = int32(width)
		bi.BiHeight = -int32(height)
		bi.BiPlanes = 1
		bi.BiBitCount = 32
		bi.BiCompression = win.BI_BITFIELDS
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
		hbitmap := win.CreateDIBSection(hdc, &bi.BITMAPINFOHEADER, win.DIB_RGB_COLORS, &lpBits, 0, 0)

		bitsPerPixel := 32
		length := ((width*bitsPerPixel + 31) / 32) * 4 * height
		// Slice memory layout
		var sl = struct {
			addr uintptr
			len  int
			cap  int
		}{uintptr(lpBits), length, length}

		// Use unsafe to turn sl into a []byte.
		bitmapArray := *(*[]byte)(unsafe.Pointer(&sl))

		i := 0
		for y := 0; y != height; y++ {
			for x := 0; x != width; x++ {
				bitmapArray[i+3] = 255
				bitmapArray[i+2] = 255
				bitmapArray[i+1] = 0
				bitmapArray[i+0] = 0
				i += 4
			}
		}

		memDC := win.CreateCompatibleDC(g.GetHDC())
		win.SelectObject(memDC, win.HGDIOBJ(hbitmap))

		gbitmap := NewGraphics(memDC)
		gbitmap.DrawLine(10, 10, 100, 100, Rgb(0, 255, 0))

		win.BitBlt(g.GetHDC(), 0, 0, int32(width), int32(height), memDC, 0, 0, win.SRCCOPY)

		win.DeleteObject(win.HGDIOBJ(hbitmap))
	*/

	var bf win.BLENDFUNCTION
	bf.BlendOp = AC_SRC_OVER
	bf.BlendFlags = 0
	bf.SourceConstantAlpha = 255
	bf.AlphaFormat = win.AC_SRC_ALPHA

	memDC := win.CreateCompatibleDC(0)
	prevObj := win.SelectObject(memDC, win.HGDIOBJ(image.hbitmap))

	width := int32(image.Width())
	height := int32(image.Height())

	//win.BitBlt(g.GetHDC(), 0, 0, int32(width), int32(height), memDC, 0, 0, win.SRCCOPY)
	win.AlphaBlend(g.GetHDC(), int32(x), int32(y), width, height, memDC, 0, 0, width, height, bf)

	win.SelectObject(memDC, prevObj)
	win.DeleteDC(memDC)

}
