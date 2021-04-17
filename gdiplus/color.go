package gdiplus

type Color struct {
	Argb ARGB
}

func MakeARGB(a, r, g, b byte) ARGB {
	return ((ARGB(b) << BlueShift) | (ARGB(g) << GreenShift) | (ARGB(r) << RedShift) | (ARGB(a) << AlphaShift))
}

func NewColor(r, g, b, a byte) *Color {
	c := &Color{}
	c.Argb = MakeARGB(a, r, g, b)
	return c
}

func (c *Color) GetAlpha() byte {
	return byte(c.Argb >> AlphaShift)
}

func (c *Color) GetA() byte {
	return c.GetAlpha()
}

func (c *Color) GetRed() byte {
	return byte(c.Argb >> RedShift)
}

func (c *Color) GetR() byte {
	return c.GetRed()
}

func (c *Color) GetGreen() byte {
	return byte(c.Argb >> GreenShift)
}

func (c *Color) GetG() byte {
	return c.GetGreen()
}

func (c *Color) GetBlue() byte {
	return byte(c.Argb >> BlueShift)
}

func (c *Color) GetB() byte {
	return c.GetBlue()
}

func (c *Color) GetValue() ARGB {
	return c.Argb
}
