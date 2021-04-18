package reza

import win "github.com/lxn/win"

// DoubleBuffer lets us use the double buffered drawing in order to get rid of flickerings during repaint/redraw
type DoubleBuffer struct {
	orgGraphics *Graphics
	newGraphics *Graphics
	hdc         win.HDC
	bitmap      win.HBITMAP
	rect        Rect
}

func NewDoubleBuffer(window Window, rc *Rect, color *Color) *DoubleBuffer {
	db := &DoubleBuffer{}
	db.rect = *rc // copy
	hdc := win.GetDC(window.GetHandle())
	db.hdc = win.CreateCompatibleDC(hdc)
	db.bitmap = win.CreateCompatibleBitmap(hdc, int32(rc.Width()), int32(rc.Height()))
	win.SelectObject(db.hdc, win.HGDIOBJ(db.bitmap))
	win.ReleaseDC(window.GetHandle(), hdc)
	db.newGraphics = &Graphics{db.hdc}
	brushWhite := NewSolidBrush(color)
	wrect := rc.AsRECT()
	FillRect(db.hdc, &wrect, win.HBRUSH(brushWhite.GetGdiObject()))
	brushWhite.Dispose()
	return db
}

func (db *DoubleBuffer) Dispose() {
	if db.bitmap != 0 {
		win.DeleteObject(win.HGDIOBJ(db.bitmap))
	}
	if db.hdc != 0 {
		win.DeleteDC(db.hdc)
	}
}

func (db *DoubleBuffer) GetGraphics() *Graphics {
	return db.newGraphics
}

func (db *DoubleBuffer) BitBlt(hdc win.HDC) {
	win.BitBlt(hdc, 0, 0, int32(db.rect.Width()), int32(db.rect.Height()), db.hdc, 0, 0, win.SRCCOPY)
}

// BeginDoubleBuffer begins
func (db *DoubleBuffer) BeginDoubleBuffer(g *Graphics, rc *Rect, border, fill *Color) *Graphics {
	db.orgGraphics = g
	db.rect = *rc // copy
	db.hdc = win.CreateCompatibleDC(g.hdc)
	db.bitmap = win.CreateCompatibleBitmap(g.hdc, int32(rc.Width()), int32(rc.Height()))
	win.SelectObject(db.hdc, win.HGDIOBJ(db.bitmap))
	db.newGraphics = &Graphics{db.hdc}
	db.newGraphics.FillRectangle(rc, border, fill)
	defFont := win.GetStockObject(win.DEFAULT_GUI_FONT)
	win.SelectObject(db.hdc, win.HGDIOBJ(defFont))
	return db.newGraphics
}

// EndDoubleBuffer ends it
func (db *DoubleBuffer) EndDoubleBuffer() {
	win.BitBlt(db.orgGraphics.hdc, 0, 0, int32(db.rect.Width()), int32(db.rect.Height()), db.hdc, 0, 0, win.SRCCOPY)
	win.DeleteObject(win.HGDIOBJ(db.bitmap))
	win.DeleteDC(db.hdc)
}
