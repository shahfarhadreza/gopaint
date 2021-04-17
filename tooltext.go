package main

import (
	"gopaint/gdiplus"
	. "gopaint/reza"
	"log"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type ToolText struct {
	ToolBasic
	textColor Color
	textEdit  *TextEdit
	typing    bool
	textArea  Rect
	resizer   *SelectionRect
	resizing  bool
}

func (tool *ToolText) initialize() {
	tool.textEdit = NewTextEdit()
	tool.textColor = Rgb(255, 0, 0)
	tool.resizer = NewSelectionRect()
	tool.resizing = false
}

func (tool *ToolText) Dispose() {
	if tool.resizer != nil {
		tool.resizer.Dispose()
	}
	if tool.textEdit != nil {
		tool.textEdit.Dispose()
	}
}

func enumFontFamExProc(lpelfe *win.LOGFONT, lpntme *TEXTMETRIC, FontType uint32, lParam uintptr) uintptr {
	fontName := syscall.UTF16ToString(lpelfe.LfFaceName[:])
	fontList := (*map[string]bool)(unsafe.Pointer(lParam))
	(*fontList)[fontName] = true
	log.Println(fontName)
	return 1
}

func (tool *ToolText) prepare() {
	canvas := mainWindow.workspace.canvas
	win.SetFocus(canvas.GetHandle())
	tool.resizing = false
	//tool.newTextEditAt(&Point{X: 150, Y: 50})
	//tool.textEdit.AppendText("I am Reza.\nThanks!")
	/*
		fontList := make(map[string]bool)
		var lf win.LOGFONT
		lf.LfFaceName[0] = 0
		lf.LfCharSet = win.DEFAULT_CHARSET
		hwnd := mainWindow.workspace.canvas.GetHandle()
		hdc := win.GetDC(hwnd)
		EnumFontFamiliesEx(hdc, &lf, syscall.NewCallback(enumFontFamExProc), uintptr(unsafe.Pointer(&fontList)), 0)
		win.ReleaseDC(hwnd, hdc)

		//for i := range fontList {
		//log.Println(i)
		//}
	*/
}

func (tool *ToolText) getCursor(ptMouse *Point) win.HCURSOR {
	if tool.typing {
		if onpoint, point := tool.resizer.GetClosestRectPoint(ptMouse, 6); onpoint {
			switch point {
			case RectPointTop, RectPointBottom:
				return mainWindow.hCursorSizeNS
			case RectPointLeft, RectPointRight:
				return mainWindow.hCursorSizeWE
			case RectPointTopLeft, RectPointBottomRight:
				return mainWindow.hCursorSizeNWSE
			case RectPointTopRight, RectPointBottomLeft:
				return mainWindow.hCursorSizeNESW
			}
		}
	}
	return mainWindow.hCursorIBeam
}

func (tool *ToolText) keyPressEvent(e *ToolKeyEvent) {
	canvas := e.canvas
	keycode := e.keycode
	//log.Printf("key %d\n", keycode)
	if tool.typing {
		if keycode == win.VK_ESCAPE {
			textEdit := tool.textEdit
			textEdit.Clear()
			tool.typing = false
		} else {
			tool.textEdit.KeyPressEvent(keycode)
		}
		canvas.RepaintVisible()
	}
}

func (tool *ToolText) newTextEditAt(pt *Point) {
	tool.textEdit.x = pt.X
	tool.textEdit.y = pt.Y
	tool.textEdit.UpdateLines()
	tool.typing = true
}

func (tool *ToolText) draw(e *ToolDrawEvent) {
	g := e.gdi32
	graphics := e.graphics
	gdicolor := getColorForeground()
	color := fromGdiplusColor(&gdicolor)
	textEdit := tool.textEdit

	if textEdit == nil {
		return
	}

	if tool.typing {
		if textEdit.Length() > 0 {
			tool.textArea = textEdit.GetTextArea()
		} else {
			textArea := Rect{
				Left:   textEdit.x,
				Top:    textEdit.y,
				Right:  textEdit.x + 150,
				Bottom: textEdit.y + 20,
			}
			outRect := &gdiplus.RectF{}
			graphics.MeasureStringEx("AAAAAAAAAAAAAAA", textEdit.font, &gdiplus.RectF{}, textEdit.format, outRect, nil, nil)
			//textBasic := *g.MeasureText("AAAAAAAAAAAAAAA", win.DT_LEFT, textEdit.font.GetHandle())
			textArea.Right = textArea.Left + int(outRect.Width)  //textBasic.Width()
			textArea.Bottom = textArea.Top + int(outRect.Height) //textBasic.Height()
			tool.textArea = textArea
		}
		rectBorder := tool.textArea
		rectBorder.Inflate(7, 5)

		//g.DrawFillRectangle(&tool.textArea, Rgb(10, 0, 255), Rgb(10, 0, 255))
		textEdit.Draw(g, &color)

		tool.resizer.SetRect(&rectBorder)
		tool.resizer.Draw(g)
	}
}

func (tool *ToolText) leave() {
	if tool.typing {
		tool.finalizeText()
	}
}

func (tool *ToolText) finalizeText() {
	g := mainWindow.workspace.canvas.image.context
	textEdit := tool.textEdit
	if !textEdit.IsEmpty() {
		gdicolor := getColorForeground()
		//textColor := fromGdiplusColor(&gdicolor)
		text := textEdit.GetText()
		//RenderText(g, text, &tool.textArea, &textColor, tool.textEdit.font)
		//gc.DrawText(text, &tool.textArea, win.DT_LEFT|win.DT_EXPANDTABS, color, tool.textEdit.font.GetHandle())
		brush := gdiplus.NewSolidBrush(&gdicolor)
		lrect := &gdiplus.RectF{X: float32(tool.textArea.Left), Y: float32(tool.textArea.Top), Width: 0, Height: 0}
		g.DrawStringEx(text, textEdit.font, lrect, textEdit.format, brush.AsBrush())
		brush.Dispose()
		textEdit.Clear()
	}
	tool.typing = false
}

func (tool *ToolText) mouseDownEvent(e *ToolMouseEvent) {
	canvas := e.canvas
	win.SetFocus(canvas.GetHandle())
	if onpoint, _ := tool.resizer.GetClosestRectPoint(&e.pt, 6); onpoint {
		tool.resizing = true
	}
}

func (tool *ToolText) mouseMoveEvent(e *ToolMouseEvent) {

}

func (tool *ToolText) mouseUpEvent(e *ToolMouseEvent) {
	mbutton := e.mbutton
	pt := e.pt
	if mbutton == MouseButtonLeft || mbutton == MouseButtonRight {
		if !tool.resizing {
			if tool.typing {
				if !tool.textArea.IsPointInside(&pt) {
					tool.finalizeText()
				}
			} else {
				tool.newTextEditAt(&pt)
			}
		} else {
			tool.resizing = false
		}
	}
}
