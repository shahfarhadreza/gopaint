package gdiplus

const (
	UnitWorld      = 0 // 0 -- World coordinate (non-physical unit)
	UnitDisplay    = 1 // 1 -- Variable -- for PageTransform only
	UnitPixel      = 2 // 2 -- Each unit is one device pixel.
	UnitPoint      = 3 // 3 -- Each unit is a printer's point, or 1/72 inch.
	UnitInch       = 4 // 4 -- Each unit is 1 inch.
	UnitDocument   = 5 // 5 -- Each unit is 1/300 inch.
	UnitMillimeter = 6 // 6 -- Each unit is 1 millimeter.
)

const (
	AlphaShift = 24
	RedShift   = 16
	GreenShift = 8
	BlueShift  = 0
)

const (
	AlphaMask = 0xff000000
	RedMask   = 0x00ff0000
	GreenMask = 0x0000ff00
	BlueMask  = 0x000000ff
)

const (
	FontStyleRegular    = 0
	FontStyleBold       = 1
	FontStyleItalic     = 2
	FontStyleBoldItalic = 3
	FontStyleUnderline  = 4
	FontStyleStrikeout  = 8
)

const LANG_NEUTRAL = 0x00

const (
	QualityModeInvalid = iota - 1
	QualityModeDefault
	QualityModeLow  // Best performance
	QualityModeHigh // Best rendering quality
)

//--------------------------------------------------------------------------
// Alpha Compositing mode constants
//--------------------------------------------------------------------------

const (
	CompositingModeSourceOver = iota // 0
	CompositingModeSourceCopy        // 1
)

//--------------------------------------------------------------------------
// Alpha Compositing quality constants
//--------------------------------------------------------------------------

const (
	CompositingQualityInvalid = iota + QualityModeInvalid
	CompositingQualityDefault
	CompositingQualityHighSpeed
	CompositingQualityHighQuality
	CompositingQualityGammaCorrected
	CompositingQualityAssumeLinear
)

const (
	InterpolationModeInvalid = iota + QualityModeInvalid
	InterpolationModeDefault
	InterpolationModeLowQuality
	InterpolationModeHighQuality
	InterpolationModeBilinear
	InterpolationModeBicubic
	InterpolationModeNearestNeighbor
	InterpolationModeHighQualityBilinear
	InterpolationModeHighQualityBicubic
)

const (
	SmoothingModeInvalid = iota + QualityModeInvalid
	SmoothingModeDefault
	SmoothingModeHighSpeed
	SmoothingModeHighQuality
	SmoothingModeNone
	SmoothingModeAntiAlias

/*
#if (GDIPVER >= 0x0110)
    SmoothingModeAntiAlias8x4 = SmoothingModeAntiAlias,
    SmoothingModeAntiAlias8x8
#endif //(GDIPVER >= 0x0110)
*/
)

//---------------------------------------------------------------------------
// Pixel Format Mode
//---------------------------------------------------------------------------

const (
	PixelOffsetModeInvalid = iota + QualityModeInvalid
	PixelOffsetModeDefault
	PixelOffsetModeHighSpeed
	PixelOffsetModeHighQuality
	PixelOffsetModeNone // No pixel offset
	PixelOffsetModeHalf // Offset by -0.5, -0.5 for fast anti-alias perf
)

//---------------------------------------------------------------------------
// Text Rendering Hint
//---------------------------------------------------------------------------

const (
	TextRenderingHintSystemDefault            = iota // Glyph with system default rendering hint
	TextRenderingHintSingleBitPerPixelGridFit        // Glyph bitmap with hinting
	TextRenderingHintSingleBitPerPixel               // Glyph bitmap without hinting
	TextRenderingHintAntiAliasGridFit                // Glyph anti-alias bitmap with hinting
	TextRenderingHintAntiAlias                       // Glyph anti-alias bitmap without hinting
	TextRenderingHintClearTypeGridFit                // Glyph CT bitmap with hinting
)

//--------------------------------------------------------------------------
// Fill mode constants
//--------------------------------------------------------------------------

const (
	FillModeAlternate = iota // 0
	FillModeWinding          // 1
)
