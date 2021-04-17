package gdiplus

type Pen struct {
	nativePen *GpPen
}

func NewPen(color *Color, width float32) *Pen {
	p := &Pen{}
	GdipCreatePen1(color.GetValue(), width, UnitWorld, &p.nativePen)
	return p
}

func (p *Pen) Dispose() {
	GdipDeletePen(p.nativePen)
}
