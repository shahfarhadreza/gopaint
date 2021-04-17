package reza

import (
	"image"
	"image/color"
)

// BGRA is like RGBA but in blue, green, red, alpha byte order
type BGRA struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

// NewBGRA returns a new BGRA image with the given bounds.
func NewBGRA(r image.Rectangle) *BGRA {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 4*w*h)
	return &BGRA{buf, 4 * w, r}
}

func NewBGRAWithData(r image.Rectangle, data []uint8) *BGRA {
	w, h := r.Dx(), r.Dy()
	if len(data) < 4*w*h {
		panic("not enough data supplied")
	}
	return &BGRA{data, 4 * w, r}
}

func (p *BGRA) ColorModel() color.Model { return color.RGBAModel }

func (p *BGRA) Bounds() image.Rectangle { return p.Rect }

func (p *BGRA) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *BGRA) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := p.PixOffset(x, y)
	return color.RGBA{p.Pix[i+2], p.Pix[i+1], p.Pix[i+0], p.Pix[i+3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *BGRA) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

func (p *BGRA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	p.Pix[i+0] = c1.B
	p.Pix[i+1] = c1.G
	p.Pix[i+2] = c1.R
	p.Pix[i+3] = c1.A
}

func (p *BGRA) SetRGBA(x, y int, c color.RGBA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i+0] = c.B
	p.Pix[i+1] = c.G
	p.Pix[i+2] = c.R
	p.Pix[i+3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *BGRA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not
	// guaranteed to be inside either r1 or r2 if the intersection
	// is empty. Without explicitly checking for this, the Pix[i:]
	// expression below can panic.
	if r.Empty() {
		return &BGRA{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &BGRA{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *BGRA) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0 := 0
	i1 := p.Rect.Dx() * 4
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 4 {
			if p.Pix[i+3] != 0xff {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}
