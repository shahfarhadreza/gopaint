package gdiplus

type GraphicsPath struct {
	nativePath *GpPath
}

func NewPath(fillMode int32) *GraphicsPath {
	p := &GraphicsPath{}
	GdipCreatePath(fillMode, &p.nativePath)
	return p
}

func (p *GraphicsPath) AddArcRect(rect *Rect, startAngle, sweepAngle float32) {
	GdipAddPathArcI(p.nativePath, rect.X, rect.Y, rect.Width, rect.Height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddArcRectF(rect *RectF, startAngle, sweepAngle float32) {
	GdipAddPathArc(p.nativePath, rect.X, rect.Y, rect.Width, rect.Height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddArc(x, y, width, height, startAngle, sweepAngle float32) {
	GdipAddPathArc(p.nativePath, x, y, width, height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddArcI(x, y, width, height int32, startAngle, sweepAngle float32) {
	GdipAddPathArcI(p.nativePath, x, y, width, height, startAngle, sweepAngle)
}

func (p *GraphicsPath) AddLine(x1, y1, x2, y2 float32) {
	GdipAddPathLine(p.nativePath, x1, y1, x2, y2)
}

func (p *GraphicsPath) AddLineI(x1, y1, x2, y2 int32) {
	GdipAddPathLineI(p.nativePath, x1, y1, x2, y2)
}

func (p *GraphicsPath) CloseAllFigures() {
	GdipClosePathFigures(p.nativePath)
}

func (p *GraphicsPath) CloseFigure() {
	GdipClosePathFigure(p.nativePath)
}

func (p *GraphicsPath) Dispose() {
	GdipDeletePath(p.nativePath)
}
