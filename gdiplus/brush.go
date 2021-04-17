package gdiplus

type Brush struct {
	nativeBrush *GpBrush
}

type SolidBrush struct {
	Brush
}

func NewSolidBrush(color *Color) *SolidBrush {
	b := &SolidBrush{}
	var solidFill *GpSolidFill
	GdipCreateSolidFill(color.GetValue(), &solidFill)
	b.nativeBrush = &solidFill.GpBrush
	return b
}

func (b *SolidBrush) AsBrush() *Brush {
	return &b.Brush
}

func (b *Brush) Dispose() {
	GdipDeleteBrush(b.nativeBrush)
}
