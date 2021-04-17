package gdiplus

import (
	"log"
	"syscall"

	win "github.com/lxn/win"
)

type Bitmap struct {
	Image
}

func NewBitmap(width, height int32, format PixelFormat) *Bitmap {
	bitmap := &Bitmap{}
	var nativeBitmap *win.GpBitmap
	status := GdipCreateBitmapFromScan0(width, height, 0, format, nil, &nativeBitmap)
	if status != win.Ok {
		log.Panicln(status.String())
	}
	bitmap.nativeImage = (*win.GpImage)(nativeBitmap)
	return bitmap
}

func NewBitmapEx(width, height, stride int32, format PixelFormat, scan0 *byte) *Bitmap {
	bitmap := &Bitmap{}
	var nativeBitmap *win.GpBitmap
	GdipCreateBitmapFromScan0(width, height, stride, format, scan0, &nativeBitmap)
	bitmap.nativeImage = (*win.GpImage)(nativeBitmap)
	return bitmap
}

func NewBitmapFromHBITMAP(hbitmap win.HBITMAP) *Bitmap {
	bitmap := &Bitmap{}
	var nativeBitmap *win.GpBitmap
	win.GdipCreateBitmapFromHBITMAP(hbitmap, 0, &nativeBitmap)
	bitmap.nativeImage = (*win.GpImage)(nativeBitmap)
	return bitmap
}

func NewBitmapFromFile(fileName string) *Bitmap {
	bitmap := &Bitmap{}
	fileNameUTF16, _ := syscall.UTF16PtrFromString(fileName)
	var nativeBitmap *win.GpBitmap
	win.GdipCreateBitmapFromFile(fileNameUTF16, &nativeBitmap)
	bitmap.nativeImage = (*win.GpImage)(nativeBitmap)
	return bitmap
}

func (bitmap *Bitmap) Dispose() {
	win.GdipDisposeImage(bitmap.nativeImage)
}
