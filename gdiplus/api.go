package gdiplus

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	win "github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var (
	// Graphics
	gdipCreateFromHDC          *windows.LazyProc
	gdipDeleteGraphics         *windows.LazyProc
	gdipSetInterpolationMode   *windows.LazyProc
	gdipSetSmoothingMode       *windows.LazyProc
	gdipSetPixelOffsetMode     *windows.LazyProc
	gdipSetCompositingQuality  *windows.LazyProc
	gdipSetCompositingMode     *windows.LazyProc
	gdipSetTextRenderingHint   *windows.LazyProc
	gdipGraphicsClear          *windows.LazyProc
	gdipDrawLine               *windows.LazyProc
	gdipDrawLineI              *windows.LazyProc
	gdipDrawRectangleI         *windows.LazyProc
	gdipDrawRectangle          *windows.LazyProc
	gdipDrawEllipseI           *windows.LazyProc
	gdipDrawEllipse            *windows.LazyProc
	gdipDrawPolygonI           *windows.LazyProc
	gdipDrawPolygon            *windows.LazyProc
	gdipDrawPath               *windows.LazyProc
	gdipDrawString             *windows.LazyProc
	gdipDrawImage              *windows.LazyProc
	gdipDrawImageI             *windows.LazyProc
	gdipDrawImageRect          *windows.LazyProc
	gdipDrawImageRectI         *windows.LazyProc
	gdipFillRectangle          *windows.LazyProc
	gdipFillRectangleI         *windows.LazyProc
	gdipFillPolygon            *windows.LazyProc
	gdipFillPolygonI           *windows.LazyProc
	gdipFillPath               *windows.LazyProc
	gdipFillEllipse            *windows.LazyProc
	gdipFillEllipseI           *windows.LazyProc
	gdipMeasureString          *windows.LazyProc
	gdipMeasureCharacterRanges *windows.LazyProc
	// Pen
	gdipCreatePen1 *windows.LazyProc
	gdipDeletePen  *windows.LazyProc
	// Brush
	gdipCreateSolidFill *windows.LazyProc
	gdipDeleteBrush     *windows.LazyProc
	// Image
	gdipLoadImageFromFile       *windows.LazyProc
	gdipSaveImageToFile         *windows.LazyProc
	gdipGetImageWidth           *windows.LazyProc
	gdipGetImageHeight          *windows.LazyProc
	gdipGetImageGraphicsContext *windows.LazyProc
	// Bitmap
	gdipCreateBitmapFromScan0 *windows.LazyProc
	// Font
	gdipCreateFontFromDC           *windows.LazyProc
	gdipCreateFont                 *windows.LazyProc
	gdipDeleteFont                 *windows.LazyProc
	gdipNewInstalledFontCollection *windows.LazyProc
	gdipCreateFontFamilyFromName   *windows.LazyProc
	gdipDeleteFontFamily           *windows.LazyProc
	// StringFormat
	gdipCreateStringFormat                *windows.LazyProc
	gdipDeleteStringFormat                *windows.LazyProc
	gdipStringFormatGetGenericTypographic *windows.LazyProc
	// Path
	gdipCreatePath       *windows.LazyProc
	gdipDeletePath       *windows.LazyProc
	gdipAddPathArc       *windows.LazyProc
	gdipAddPathArcI      *windows.LazyProc
	gdipAddPathLine      *windows.LazyProc
	gdipAddPathLineI     *windows.LazyProc
	gdipClosePathFigure  *windows.LazyProc
	gdipClosePathFigures *windows.LazyProc
)

func init() {
	// Library
	libgdiplus := windows.NewLazySystemDLL("gdiplus.dll")
	// Graphics
	gdipCreateFromHDC = libgdiplus.NewProc("GdipCreateFromHDC")
	gdipDeleteGraphics = libgdiplus.NewProc("GdipDeleteGraphics")
	gdipSetInterpolationMode = libgdiplus.NewProc("GdipSetInterpolationMode")
	gdipSetSmoothingMode = libgdiplus.NewProc("GdipSetSmoothingMode")
	gdipSetPixelOffsetMode = libgdiplus.NewProc("GdipSetPixelOffsetMode")
	gdipSetCompositingQuality = libgdiplus.NewProc("GdipSetCompositingQuality")
	gdipSetCompositingMode = libgdiplus.NewProc("GdipSetCompositingMode")
	gdipSetTextRenderingHint = libgdiplus.NewProc("GdipSetTextRenderingHint")
	gdipGraphicsClear = libgdiplus.NewProc("GdipGraphicsClear")
	gdipDrawLine = libgdiplus.NewProc("GdipDrawLine")
	gdipDrawLineI = libgdiplus.NewProc("GdipDrawLineI")
	gdipDrawRectangleI = libgdiplus.NewProc("GdipDrawRectangleI")
	gdipDrawRectangle = libgdiplus.NewProc("GdipDrawRectangle")
	gdipDrawEllipseI = libgdiplus.NewProc("GdipDrawEllipseI")
	gdipDrawEllipse = libgdiplus.NewProc("GdipDrawEllipse")
	gdipDrawPolygonI = libgdiplus.NewProc("GdipDrawPolygonI")
	gdipDrawPolygon = libgdiplus.NewProc("GdipDrawPolygon")
	gdipDrawPath = libgdiplus.NewProc("GdipDrawPath")
	gdipDrawString = libgdiplus.NewProc("GdipDrawString")
	gdipDrawImage = libgdiplus.NewProc("GdipDrawImage")
	gdipDrawImageI = libgdiplus.NewProc("GdipDrawImageI")
	gdipDrawImageRect = libgdiplus.NewProc("GdipDrawImageRect")
	gdipDrawImageRectI = libgdiplus.NewProc("GdipDrawImageRectI")
	gdipFillRectangle = libgdiplus.NewProc("GdipFillRectangle")
	gdipFillRectangleI = libgdiplus.NewProc("GdipFillRectangleI")
	gdipFillPolygon = libgdiplus.NewProc("GdipFillPolygon")
	gdipFillPolygonI = libgdiplus.NewProc("GdipFillPolygonI")
	gdipFillPath = libgdiplus.NewProc("GdipFillPath")
	gdipFillEllipse = libgdiplus.NewProc("GdipFillEllipse")
	gdipFillEllipseI = libgdiplus.NewProc("GdipFillEllipseI")
	gdipMeasureString = libgdiplus.NewProc("GdipMeasureString")
	gdipMeasureCharacterRanges = libgdiplus.NewProc("GdipMeasureCharacterRanges")
	// Pen
	gdipCreatePen1 = libgdiplus.NewProc("GdipCreatePen1")
	gdipDeletePen = libgdiplus.NewProc("GdipDeletePen")
	// Brush
	gdipCreateSolidFill = libgdiplus.NewProc("GdipCreateSolidFill")
	gdipDeleteBrush = libgdiplus.NewProc("GdipDeleteBrush")
	// Image
	gdipLoadImageFromFile = libgdiplus.NewProc("GdipLoadImageFromFile")
	gdipSaveImageToFile = libgdiplus.NewProc("GdipSaveImageToFile")
	gdipGetImageWidth = libgdiplus.NewProc("GdipGetImageWidth")
	gdipGetImageHeight = libgdiplus.NewProc("GdipGetImageHeight")
	gdipGetImageGraphicsContext = libgdiplus.NewProc("GdipGetImageGraphicsContext")
	// Bitmap
	gdipCreateBitmapFromScan0 = libgdiplus.NewProc("GdipCreateBitmapFromScan0")
	// Font
	gdipCreateFontFromDC = libgdiplus.NewProc("GdipCreateFontFromDC")
	gdipCreateFont = libgdiplus.NewProc("GdipCreateFont")
	gdipDeleteFont = libgdiplus.NewProc("GdipDeleteFont")
	gdipNewInstalledFontCollection = libgdiplus.NewProc("GdipNewInstalledFontCollection")
	gdipCreateFontFamilyFromName = libgdiplus.NewProc("GdipCreateFontFamilyFromName")
	gdipDeleteFontFamily = libgdiplus.NewProc("GdipDeleteFontFamily")
	// StringFormat
	gdipCreateStringFormat = libgdiplus.NewProc("GdipCreateStringFormat")
	gdipDeleteStringFormat = libgdiplus.NewProc("GdipDeleteStringFormat")
	gdipStringFormatGetGenericTypographic = libgdiplus.NewProc("GdipStringFormatGetGenericTypographic")
	// Path
	gdipCreatePath = libgdiplus.NewProc("GdipCreatePath")
	gdipDeletePath = libgdiplus.NewProc("GdipDeletePath")
	gdipAddPathArc = libgdiplus.NewProc("GdipAddPathArc")
	gdipAddPathArcI = libgdiplus.NewProc("GdipAddPathArcI")
	gdipAddPathLine = libgdiplus.NewProc("GdipAddPathLine")
	gdipAddPathLineI = libgdiplus.NewProc("GdipAddPathLineI")
	gdipClosePathFigure = libgdiplus.NewProc("GdipClosePathFigure")
	gdipClosePathFigures = libgdiplus.NewProc("GdipClosePathFigures")
}

// Graphics
func GdipCreateFromHDC(hdc win.HDC, graphics **GpGraphics) win.GpStatus {
	ret, _, _ := gdipCreateFromHDC.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(graphics)))
	return win.GpStatus(ret)
}

func GdipGetImageGraphicsContext(image *win.GpImage, graphics **GpGraphics) win.GpStatus {
	ret, _, _ := gdipGetImageGraphicsContext.Call(
		uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(graphics)))
	return win.GpStatus(ret)
}

func GdipDeleteGraphics(graphics *GpGraphics) win.GpStatus {
	ret, _, _ := gdipDeleteGraphics.Call(uintptr(unsafe.Pointer(graphics)))
	return win.GpStatus(ret)
}

func GdipSetCompositingMode(graphics *GpGraphics, mode int32) win.GpStatus {
	ret, _, _ := gdipSetCompositingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return win.GpStatus(ret)
}

func GdipSetCompositingQuality(graphics *GpGraphics, quality int32) win.GpStatus {
	ret, _, _ := gdipSetCompositingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(quality))
	return win.GpStatus(ret)
}

func GdipSetInterpolationMode(graphics *GpGraphics, mode int32) win.GpStatus {
	ret, _, _ := gdipSetInterpolationMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return win.GpStatus(ret)
}

func GdipSetPixelOffsetMode(graphics *GpGraphics, mode int32) win.GpStatus {
	ret, _, _ := gdipSetPixelOffsetMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return win.GpStatus(ret)
}

func GdipSetSmoothingMode(graphics *GpGraphics, mode int32) win.GpStatus {
	ret, _, _ := gdipSetSmoothingMode.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(mode))
	return win.GpStatus(ret)
}

func GdipSetTextRenderingHint(graphics *GpGraphics, hint int32) win.GpStatus {
	ret, _, _ := gdipSetTextRenderingHint.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(hint))
	return win.GpStatus(ret)
}

func GdipGraphicsClear(graphics *GpGraphics, color ARGB) win.GpStatus {
	ret, _, _ := gdipGraphicsClear.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(color))
	return win.GpStatus(ret)
}

func GdipDrawLine(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 float32) win.GpStatus {
	ret, _, _ := gdipDrawLine.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return win.GpStatus(ret)
}

func GdipDrawLineI(graphics *GpGraphics, pen *GpPen, x1, y1, x2, y2 int32) win.GpStatus {
	ret, _, _ := gdipDrawLineI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return win.GpStatus(ret)
}

func GdipDrawRectangle(graphics *GpGraphics, pen *GpPen, x, y, width, height float32) win.GpStatus {
	ret, _, _ := gdipDrawRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return win.GpStatus(ret)
}

func GdipDrawRectangleI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) win.GpStatus {
	ret, _, _ := gdipDrawRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return win.GpStatus(ret)
}

func GdipDrawEllipse(graphics *GpGraphics, pen *GpPen, x, y, width, height float32) win.GpStatus {
	ret, _, _ := gdipDrawEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return win.GpStatus(ret)
}

func GdipDrawEllipseI(graphics *GpGraphics, pen *GpPen, x, y, width, height int32) win.GpStatus {
	ret, _, _ := gdipDrawEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return win.GpStatus(ret)
}

func GdipDrawPolygon(graphics *GpGraphics, pen *GpPen, points *PointF, count int32) win.GpStatus {
	ret, _, _ := gdipDrawPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return win.GpStatus(ret)
}

func GdipDrawPolygonI(graphics *GpGraphics, pen *GpPen, points *Point, count int32) win.GpStatus {
	ret, _, _ := gdipDrawPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count))
	return win.GpStatus(ret)
}

func GdipDrawPath(graphics *GpGraphics, pen *GpPen, path *GpPath) win.GpStatus {
	ret, _, _ := gdipDrawPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(pen)),
		uintptr(unsafe.Pointer(path)))
	return win.GpStatus(ret)
}

func GdipDrawString(graphics *GpGraphics, text *uint16, length int32, font *GpFont, layoutRect *RectF, stringFormat *GpStringFormat, brush *GpBrush) win.GpStatus {
	ret, _, _ := gdipDrawString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(unsafe.Pointer(brush)))
	return win.GpStatus(ret)
}

func GdipDrawImage(graphics *GpGraphics, image *win.GpImage, x, y float32) win.GpStatus {
	ret, _, _ := gdipDrawImage.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)))
	return win.GpStatus(ret)
}

func GdipDrawImageI(graphics *GpGraphics, image *win.GpImage, x, y int32) win.GpStatus {
	ret, _, _ := gdipDrawImageI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y))
	return win.GpStatus(ret)
}

func GdipDrawImageRect(graphics *GpGraphics, image *win.GpImage, x, y, width, height float32) win.GpStatus {
	ret, _, _ := gdipDrawImageRect.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return win.GpStatus(ret)
}

func GdipDrawImageRectI(graphics *GpGraphics, image *win.GpImage, x, y, width, height int32) win.GpStatus {
	ret, _, _ := gdipDrawImageRectI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(image)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return win.GpStatus(ret)
}

func GdipFillRectangle(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) win.GpStatus {
	ret, _, _ := gdipFillRectangle.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return win.GpStatus(ret)
}

func GdipFillRectangleI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) win.GpStatus {
	ret, _, _ := gdipFillRectangleI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return win.GpStatus(ret)
}

func GdipFillEllipse(graphics *GpGraphics, brush *GpBrush, x, y, width, height float32) win.GpStatus {
	ret, _, _ := gdipFillEllipse.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)))
	return win.GpStatus(ret)
}

func GdipFillEllipseI(graphics *GpGraphics, brush *GpBrush, x, y, width, height int32) win.GpStatus {
	ret, _, _ := gdipFillEllipseI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height))
	return win.GpStatus(ret)
}

func GdipFillPolygon(graphics *GpGraphics, brush *GpBrush, points *PointF, count int32, fillMode int32) win.GpStatus {
	ret, _, _ := gdipFillPolygon.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return win.GpStatus(ret)
}

func GdipFillPolygonI(graphics *GpGraphics, brush *GpBrush, points *Point, count int32, fillMode int32) win.GpStatus {
	ret, _, _ := gdipFillPolygonI.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(points)),
		uintptr(count),
		uintptr(fillMode))
	return win.GpStatus(ret)
}

func GdipFillPath(graphics *GpGraphics, brush *GpBrush, path *GpPath) win.GpStatus {
	ret, _, _ := gdipFillPath.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(brush)),
		uintptr(unsafe.Pointer(path)))
	return win.GpStatus(ret)
}

func GdipMeasureString(
	graphics *GpGraphics, text *uint16,
	length int32, font *GpFont, layoutRect *RectF,
	stringFormat *GpStringFormat, boundingBox *RectF,
	codepointsFitted *int32, linesFilled *int32) win.GpStatus {

	ret, _, _ := gdipMeasureString.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(unsafe.Pointer(boundingBox)),
		uintptr(unsafe.Pointer(codepointsFitted)),
		uintptr(unsafe.Pointer(linesFilled)))
	return win.GpStatus(ret)
}

func GdipMeasureCharacterRanges(
	graphics *GpGraphics, text *uint16,
	length int32, font *GpFont, layoutRect *RectF,
	stringFormat *GpStringFormat, regionCount int32,
	regions **GpRegion) win.GpStatus {

	ret, _, _ := gdipMeasureCharacterRanges.Call(
		uintptr(unsafe.Pointer(graphics)),
		uintptr(unsafe.Pointer(text)),
		uintptr(length),
		uintptr(unsafe.Pointer(font)),
		uintptr(unsafe.Pointer(layoutRect)),
		uintptr(unsafe.Pointer(stringFormat)),
		uintptr(regionCount),
		uintptr(unsafe.Pointer(regions)))
	return win.GpStatus(ret)
}

// Pen
func GdipCreatePen1(color ARGB, width float32, unit GpUnit, pen **GpPen) win.GpStatus {
	ret, _, _ := gdipCreatePen1.Call(
		uintptr(color),
		uintptr(math.Float32bits(width)),
		uintptr(unit),
		uintptr(unsafe.Pointer(pen)))
	return win.GpStatus(ret)
}

func GdipDeletePen(pen *GpPen) win.GpStatus {
	ret, _, _ := gdipDeletePen.Call(uintptr(unsafe.Pointer(pen)))
	return win.GpStatus(ret)
}

func GdipCreateSolidFill(color ARGB, brush **GpSolidFill) win.GpStatus {
	ret, _, _ := gdipCreateSolidFill.Call(
		uintptr(color),
		uintptr(unsafe.Pointer(brush)))
	return win.GpStatus(ret)
}

func GdipDeleteBrush(brush *GpBrush) win.GpStatus {
	ret, _, _ := gdipDeleteBrush.Call(uintptr(unsafe.Pointer(brush)))
	return win.GpStatus(ret)
}

func GdipCreateFontFromDC(hdc win.HDC, font **GpFont) win.GpStatus {
	ret, _, _ := gdipCreateFontFromDC.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(font)))
	return win.GpStatus(ret)
}

func GdipCreateFont(fontFamily *GpFontFamily, emSize float32, style int32, unit GpUnit, font **GpFont) win.GpStatus {
	ret, _, _ := gdipCreateFont.Call(
		uintptr(unsafe.Pointer(fontFamily)),
		uintptr(math.Float32bits(emSize)),
		uintptr(style),
		uintptr(unit),
		uintptr(unsafe.Pointer(font)))
	return win.GpStatus(ret)
}

func GdipDeleteFont(font *GpFont) win.GpStatus {
	ret, _, _ := gdipDeleteFont.Call(uintptr(unsafe.Pointer(font)))
	return win.GpStatus(ret)
}

func GdipNewInstalledFontCollection(fontCollection **GpFontCollection) win.GpStatus {
	ret, _, _ := gdipNewInstalledFontCollection.Call(uintptr(unsafe.Pointer(fontCollection)))
	return win.GpStatus(ret)
}

func GdipCreateFontFamilyFromName(name *uint16, fontCollection *GpFontCollection, fontFamily **GpFontFamily) win.GpStatus {
	ret, _, _ := gdipCreateFontFamilyFromName.Call(
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(fontCollection)),
		uintptr(unsafe.Pointer(fontFamily)))
	return win.GpStatus(ret)
}

func GdipDeleteFontFamily(fontFamily *GpFontFamily) win.GpStatus {
	ret, _, _ := gdipDeleteFontFamily.Call(uintptr(unsafe.Pointer(fontFamily)))
	return win.GpStatus(ret)
}

func GdipCreateStringFormat(formatAttributes int32, language uint16, format **GpStringFormat) win.GpStatus {
	ret, _, _ := gdipCreateStringFormat.Call(
		uintptr(formatAttributes),
		uintptr(language),
		uintptr(unsafe.Pointer(format)))
	return win.GpStatus(ret)
}

func GdipStringFormatGetGenericTypographic(format **GpStringFormat) win.GpStatus {
	ret, _, _ := gdipStringFormatGetGenericTypographic.Call(uintptr(unsafe.Pointer(format)))
	return win.GpStatus(ret)
}

func GdipDeleteStringFormat(format *GpStringFormat) win.GpStatus {
	ret, _, _ := gdipDeleteStringFormat.Call(uintptr(unsafe.Pointer(format)))
	return win.GpStatus(ret)
}

func GdipCreatePath(brushMode int32, path **GpPath) win.GpStatus {
	ret, _, _ := gdipCreatePath.Call(uintptr(brushMode), uintptr(unsafe.Pointer(path)))
	return win.GpStatus(ret)
}

func GdipDeletePath(path *GpPath) win.GpStatus {
	ret, _, _ := gdipDeletePath.Call(uintptr(unsafe.Pointer(path)))
	return win.GpStatus(ret)
}

func GdipAddPathArc(path *GpPath, x, y, width, height, startAngle, sweepAngle float32) win.GpStatus {
	ret, _, _ := gdipAddPathArc.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x)),
		uintptr(math.Float32bits(y)),
		uintptr(math.Float32bits(width)),
		uintptr(math.Float32bits(height)),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return win.GpStatus(ret)
}

func GdipAddPathArcI(path *GpPath, x, y, width, height int32, startAngle, sweepAngle float32) win.GpStatus {
	ret, _, _ := gdipAddPathArcI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(math.Float32bits(startAngle)),
		uintptr(math.Float32bits(sweepAngle)))
	return win.GpStatus(ret)
}

func GdipAddPathLine(path *GpPath, x1, y1, x2, y2 float32) win.GpStatus {
	ret, _, _ := gdipAddPathLine.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(math.Float32bits(x1)),
		uintptr(math.Float32bits(y1)),
		uintptr(math.Float32bits(x2)),
		uintptr(math.Float32bits(y2)))
	return win.GpStatus(ret)
}

func GdipAddPathLineI(path *GpPath, x1, y1, x2, y2 int32) win.GpStatus {
	ret, _, _ := gdipAddPathLineI.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(x1),
		uintptr(y1),
		uintptr(x2),
		uintptr(y2))
	return win.GpStatus(ret)
}

func GdipClosePathFigure(path *GpPath) win.GpStatus {
	ret, _, _ := gdipClosePathFigure.Call(uintptr(unsafe.Pointer(path)))
	return win.GpStatus(ret)
}

func GdipClosePathFigures(path *GpPath) win.GpStatus {
	ret, _, _ := gdipClosePathFigures.Call(uintptr(unsafe.Pointer(path)))
	return win.GpStatus(ret)
}

// ----

func GdipLoadImageFromFile(filename *uint16, image **win.GpImage) win.GpStatus {
	ret, _, _ := gdipLoadImageFromFile.Call(
		uintptr(unsafe.Pointer(filename)),
		uintptr(unsafe.Pointer(image)))
	return win.GpStatus(ret)
}

func GdipSaveImageToFile(image *win.GpBitmap, filename *uint16, clsidEncoder *ole.GUID, encoderParams *EncoderParameters) win.GpStatus {
	ret, _, _ := gdipSaveImageToFile.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(filename)), uintptr(unsafe.Pointer(clsidEncoder)),
		uintptr(unsafe.Pointer(encoderParams)))
	return win.GpStatus(ret)
}

func GdipGetImageWidth(image *win.GpImage, width *uint32) win.GpStatus {
	ret, _, _ := gdipGetImageWidth.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(width)))
	return win.GpStatus(ret)
}

func GdipGetImageHeight(image *win.GpImage, height *uint32) win.GpStatus {
	ret, _, _ := gdipGetImageHeight.Call(uintptr(unsafe.Pointer(image)),
		uintptr(unsafe.Pointer(height)))
	return win.GpStatus(ret)
}

// Bitmap
func GdipCreateBitmapFromScan0(width, height, stride int32, format PixelFormat, scan0 *byte, bitmap **win.GpBitmap) win.GpStatus {
	ret, _, _ := gdipCreateBitmapFromScan0.Call(
		uintptr(width),
		uintptr(height),
		uintptr(stride),
		uintptr(format),
		uintptr(unsafe.Pointer(scan0)),
		uintptr(unsafe.Pointer(bitmap)))
	return win.GpStatus(ret)
}

func SavePNG(fileName string, newBMP win.HBITMAP) error {
	// HBITMAP
	var bmp *win.GpBitmap
	if win.GdipCreateBitmapFromHBITMAP(newBMP, 0, &bmp) != 0 {
		return fmt.Errorf("failed to create HBITMAP")
	}
	defer win.GdipDisposeImage((*win.GpImage)(bmp))
	clsid, err := ole.CLSIDFromString("{557CF406-1A04-11D3-9A73-0000F81EF32E}")
	if err != nil {
		return err
	}
	fname, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return err
	}
	if GdipSaveImageToFile(bmp, fname, clsid, nil) != 0 {
		return fmt.Errorf("failed to call PNG encoder")
	}
	return nil
}
