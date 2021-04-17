package gdiplus

import (
	"log"
	"syscall"

	win "github.com/lxn/win"
)

type Graphics struct {
	nativeGraphics *GpGraphics
}

func NewGraphicsFromHDC(hdc win.HDC) *Graphics {
	g := &Graphics{}
	GdipCreateFromHDC(hdc, &g.nativeGraphics)
	return g
}

func NewGraphicsFromBitmap(bitmap *Bitmap) *Graphics {
	g := &Graphics{}
	GdipGetImageGraphicsContext(bitmap.nativeImage, &g.nativeGraphics)
	return g
}

func (g *Graphics) SetCompositingMode(mode int32) {
	GdipSetCompositingMode(g.nativeGraphics, mode)
}

func (g *Graphics) SetCompositingQuality(quality int32) {
	GdipSetCompositingQuality(g.nativeGraphics, quality)
}

func (g *Graphics) SetInterpolationMode(mode int32) {
	GdipSetInterpolationMode(g.nativeGraphics, mode)
}

func (g *Graphics) SetPixelOffsetMode(mode int32) {
	GdipSetPixelOffsetMode(g.nativeGraphics, mode)
}

func (g *Graphics) SetSmoothingMode(mode int32) {
	GdipSetSmoothingMode(g.nativeGraphics, mode)
}

func (g *Graphics) SetTextRenderingHint(hint int32) {
	GdipSetTextRenderingHint(g.nativeGraphics, hint)
}

func (g *Graphics) Clear(color *Color) {
	GdipGraphicsClear(g.nativeGraphics, color.GetValue())
}

func (g *Graphics) DrawString(text string, font *Font, origin *PointF, brush *Brush) {
	rect := &RectF{X: origin.X, Y: origin.Y, Width: 0, Height: 0}
	textUTF16, _ := syscall.UTF16FromString(text)
	GdipDrawString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)), font.nativeFont, rect, nil, brush.nativeBrush)
}

func (g *Graphics) DrawStringEx(text string, font *Font, rect *RectF, format *StringFormat, brush *Brush) {
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat = nil
	if format != nil {
		nativeFormat = format.nativeFormat
	}
	GdipDrawString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)), font.nativeFont, rect, nativeFormat, brush.nativeBrush)
}

func (g *Graphics) MeasureStringEx(text string, font *Font, layoutRect *RectF,
	stringFormat *StringFormat, boundingBox *RectF,
	codepointsFitted *int32, linesFilled *int32) {
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat = nil
	if stringFormat != nil {
		nativeFormat = stringFormat.nativeFormat
	}
	GdipMeasureString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)), font.nativeFont, layoutRect, nativeFormat, boundingBox, codepointsFitted, linesFilled)
}

func (g *Graphics) MeasureString(text string, font *Font, layoutRect *RectF,
	stringFormat *StringFormat, boundingBox *RectF) {
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat = nil
	if stringFormat != nil {
		nativeFormat = stringFormat.nativeFormat
	}
	GdipMeasureString(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)), font.nativeFont, layoutRect, nativeFormat, boundingBox, nil, nil)
}

func (g *Graphics) MeasureCharacterRanges(text string, font *Font,
	layoutRect *RectF,
	stringFormat *StringFormat, regionCount int32,
	regions **GpRegion) {
	textUTF16, _ := syscall.UTF16FromString(text)
	var nativeFormat *GpStringFormat = nil
	if stringFormat != nil {
		nativeFormat = stringFormat.nativeFormat
	}
	GdipMeasureCharacterRanges(g.nativeGraphics, &textUTF16[0], int32(len(textUTF16)), font.nativeFont, layoutRect, nativeFormat, regionCount, regions)
}

func (g *Graphics) DrawLine(pen *Pen, x1, y1, x2, y2 float32) {
	GdipDrawLine(g.nativeGraphics, pen.nativePen, x1, y1, x2, y2)
}

func (g *Graphics) DrawLineI(pen *Pen, x1, y1, x2, y2 int32) {
	GdipDrawLineI(g.nativeGraphics, pen.nativePen, x1, y1, x2, y2)
}

func (g *Graphics) DrawRectangle(pen *Pen, x, y, width, height float32) {
	GdipDrawRectangle(g.nativeGraphics, pen.nativePen, x, y, width, height)
}

func (g *Graphics) DrawRectangleI(pen *Pen, x, y, width, height int32) {
	GdipDrawRectangleI(g.nativeGraphics, pen.nativePen, x, y, width, height)
}

func (g *Graphics) DrawEllipse(pen *Pen, x, y, width, height float32) {
	GdipDrawEllipse(g.nativeGraphics, pen.nativePen, x, y, width, height)
}

func (g *Graphics) DrawEllipseI(pen *Pen, x, y, width, height int32) {
	GdipDrawEllipseI(g.nativeGraphics, pen.nativePen, x, y, width, height)
}

func (g *Graphics) DrawPolygon(pen *Pen, points []PointF) {
	GdipDrawPolygon(g.nativeGraphics, pen.nativePen, &points[0], int32(len(points)))
}

func (g *Graphics) DrawPolygonI(pen *Pen, points []Point) {
	GdipDrawPolygonI(g.nativeGraphics, pen.nativePen, &points[0], int32(len(points)))
}

func (g *Graphics) DrawPath(pen *Pen, path *GraphicsPath) {
	GdipDrawPath(g.nativeGraphics, pen.nativePen, path.nativePath)
}

func (g *Graphics) DrawImage(image *Image, x, y float32) {
	status := GdipDrawImage(g.nativeGraphics, image.nativeImage, x, y)
	if status != win.Ok {
		log.Panicln(status.String())
	}
}

func (g *Graphics) DrawImageI(image *Image, x, y int32) {
	status := GdipDrawImageI(g.nativeGraphics, image.nativeImage, x, y)
	if status != win.Ok {
		log.Panicln(status.String())
	}
}

func (g *Graphics) DrawImageRect(image *Image, x, y, width, height float32) {
	GdipDrawImageRect(g.nativeGraphics, image.nativeImage, x, y, width, height)
}

func (g *Graphics) DrawImageRectI(image *Image, x, y, width, height int32) {
	GdipDrawImageRectI(g.nativeGraphics, image.nativeImage, x, y, width, height)
}

func (g *Graphics) FillRectangle(brush *Brush, x, y, width, height float32) {
	GdipFillRectangle(g.nativeGraphics, brush.nativeBrush, x, y, width, height)
}

func (g *Graphics) FillRectangleI(brush *Brush, x, y, width, height int32) {
	GdipFillRectangleI(g.nativeGraphics, brush.nativeBrush, x, y, width, height)
}

func (g *Graphics) FillEllipse(brush *Brush, x, y, width, height float32) {
	GdipFillEllipse(g.nativeGraphics, brush.nativeBrush, x, y, width, height)
}

func (g *Graphics) FillEllipseI(brush *Brush, x, y, width, height int32) {
	GdipFillEllipseI(g.nativeGraphics, brush.nativeBrush, x, y, width, height)
}

func (g *Graphics) FillPolygon(brush *Brush, points []PointF, fillMode int32) {
	GdipFillPolygon(g.nativeGraphics, brush.nativeBrush, &points[0], int32(len(points)), fillMode)
}

func (g *Graphics) FillPolygonI(brush *Brush, points []Point, fillMode int32) {
	GdipFillPolygonI(g.nativeGraphics, brush.nativeBrush, &points[0], int32(len(points)), fillMode)
}

func (g *Graphics) FillPath(brush *Brush, path *GraphicsPath) {
	GdipFillPath(g.nativeGraphics, brush.nativeBrush, path.nativePath)
}

func (g *Graphics) Dispose() {
	GdipDeleteGraphics(g.nativeGraphics)
}
