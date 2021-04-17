package gdiplus

import (
	"log"
	"syscall"
)

type FontCollection struct {
	nativeFontCollection *GpFontCollection
}

type FontFamily struct {
	nativeFontFamily *GpFontFamily
}

type Font struct {
	nativeFont *GpFont
}

func NewFontFamily(familyName string, fontCollection *FontCollection) *FontFamily {
	fm := &FontFamily{}
	if fontCollection != nil {
		log.Panicln("FontCollection not implemented!")
	}
	familyNameUTF16, _ := syscall.UTF16PtrFromString(familyName)
	GdipCreateFontFamilyFromName(familyNameUTF16, nil, &fm.nativeFontFamily)
	return fm
}

func (fm *FontFamily) Dispose() {
	GdipDeleteFontFamily(fm.nativeFontFamily)
}

func NewFont(familyName string, emSize float32, style int32, unit GpUnit, fontCollection *FontCollection) *Font {
	f := &Font{}
	fontFamily := NewFontFamily(familyName, fontCollection)
	defer fontFamily.Dispose()
	GdipCreateFont(fontFamily.nativeFontFamily, emSize, style, unit, &f.nativeFont)
	return f
}

func (f *Font) Dispose() {
	GdipDeleteFont(f.nativeFont)
}
