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
}

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
