package gdiplus

type StringFormat struct {
	nativeFormat *GpStringFormat
}

func NewStringFormat() *StringFormat {
	format := &StringFormat{}
	GdipCreateStringFormat(0, LANG_NEUTRAL, &format.nativeFormat)
	return format
}

func NewGenericTypographicStringFormat() *StringFormat {
	format := &StringFormat{}
	GdipStringFormatGetGenericTypographic(&format.nativeFormat)
	return format
}

func (format *StringFormat) Dispose() {
	GdipDeleteStringFormat(format.nativeFormat)
}
