package main

import (
	"gopaint/gdiplus"
	. "gopaint/reza"
	"log"

	win "github.com/lxn/win"
)

type TextEdit struct {
	x             int
	y             int
	font          *gdiplus.Font
	format        *gdiplus.StringFormat
	lines         []*TextLine
	buffer        *TextBuffer
	caratPosition int
	penCarat      *Pen
}

type TextLine struct {
	chunk      *TextBuffer
	startIndex int
	rect       Rect
}

type TextBuffer struct {
	chars []rune
}

func NewTextBuffer() *TextBuffer {
	tb := &TextBuffer{}
	tb.chars = make([]rune, 0)
	return tb
}

func (tb *TextBuffer) Clear() {
	tb.chars = tb.chars[:0]
}

func (tb *TextBuffer) AsString() string {
	return string(tb.chars)
}

func (tb *TextBuffer) Append(char rune) {
	tb.chars = append(tb.chars, char)
}

func (tb *TextBuffer) Insert(pos int, char rune) {
	tb.chars = append(tb.chars, char)
}

func (tb *TextBuffer) InsertText(pos int, text string) {
	runes := []rune(text)
	tb.chars = append(tb.chars, runes...)
}

func (tb *TextBuffer) Delete(position int, length int) {
	// perform bounds checking first
	bufferLength := len(tb.chars)
	if position >= bufferLength || position < 0 {
		log.Panicf("Index is out of range. Index is %d with slice length %d", position, bufferLength)
	}
	tb.chars = append(tb.chars[:position], tb.chars[position+length:]...)
}

func (tb *TextBuffer) Length() int {
	return len(tb.chars)
}

func NewTextLine() *TextLine {
	tl := &TextLine{}
	tl.chunk = NewTextBuffer()
	return tl
}

func NewTextEdit() *TextEdit {
	te := &TextEdit{}
	te.buffer = NewTextBuffer()
	te.caratPosition = 0
	te.penCarat = NewSolidPen(1, NewRgb(0, 0, 0))
	te.font = gdiplus.NewFont("Arial Black", 20, gdiplus.FontStyleRegular, gdiplus.UnitPoint, nil) //CreateDPIAwareFont("Arial Black", 20)
	te.format = gdiplus.NewGenericTypographicStringFormat()
	return te
}

func (te *TextEdit) Dispose() {
	if te.format != nil {
		te.format.Dispose()
	}
	if te.font != nil {
		te.font.Dispose()
	}
	te.penCarat.Dispose()
}

func (te *TextEdit) Clear() {
	te.buffer.Clear()
	te.caratPosition = 0
	te.UpdateLines()
}

func (te *TextEdit) GetText() string {
	return te.buffer.AsString()
}

func (te *TextEdit) IsEmpty() bool {
	return te.buffer.Length() < 1
}

func (te *TextEdit) AppendText(text string) {
	length := te.buffer.Length()
	te.buffer.InsertText(length, text)
	te.caratPosition = te.caratPosition + len(text)
	te.UpdateLines()
}

func (te *TextEdit) Insert(position int, char rune) {
	te.buffer.Insert(position, char)
	te.caratPosition = te.caratPosition + 1
	te.UpdateLines()
}

func (te *TextEdit) InsertText(position int, text string) {
	te.buffer.InsertText(position, text)
	te.caratPosition = te.caratPosition + len(text)
	te.UpdateLines()
}

func (te *TextEdit) Delete(position int, length int) {
	te.buffer.Delete(position, length)
	te.caratPosition = te.caratPosition - length
	te.UpdateLines()
}

func (te *TextEdit) DeleteBack() {
	len := te.buffer.Length()
	if len > 0 {
		te.Delete(len-1, 1)
	}
}

func (te *TextEdit) Length() int {
	return te.buffer.Length()
}

func (te *TextEdit) UpdateLines() {
	te.lines = make([]*TextLine, 0)

	canvas := mainWindow.workspace.canvas
	hdc := canvas.GetDC()
	defer canvas.ReleaseDC(hdc)
	//g2 := NewGraphics(hdc)

	g := gdiplus.NewGraphicsFromHDC(hdc)
	g.SetTextRenderingHint(gdiplus.TextRenderingHintAntiAlias)

	newLine := NewTextLine()
	newLine.startIndex = 0
	te.lines = append(te.lines, newLine)
	for i := range te.buffer.chars {
		ch := te.buffer.chars[i]
		if ch == '\r' || ch == '\n' {
			newLine = NewTextLine()
			newLine.startIndex = i + 1 // skip the newline char
			te.lines = append(te.lines, newLine)
		} else {
			newLine.chunk.Append(ch)
		}
	}

	lineX := te.x
	lineY := te.y
	textHeight := 0
	textWidth := 0
	for i := range te.lines {
		line := te.lines[i]
		charCount := line.chunk.Length()
		if charCount > 0 {
			lineText := line.chunk.AsString()
			//line.rect = *g2.MeasureText(lineText, win.DT_LEFT, te.font.GetHandle())

			lrect := &gdiplus.RectF{}
			outRect := &gdiplus.RectF{}
			g.MeasureStringEx(lineText, te.font, lrect, te.format, outRect, nil, nil)
			line.rect.Right = int(outRect.Width)
			line.rect.Bottom = int(outRect.Height)

			textWidth = line.rect.Width()
			textHeight = line.rect.Height()
		} else {
			// Give a basic height
			//line.rect = *g.MeasureText("A", win.DT_LEFT, te.font.GetHandle())

			lrect := &gdiplus.RectF{}
			outRect := &gdiplus.RectF{}
			g.MeasureStringEx("A", te.font, lrect, te.format, outRect, nil, nil)
			line.rect.Right = int(outRect.Width)
			line.rect.Bottom = int(outRect.Height)

			textHeight = line.rect.Height()
		}

		line.rect.Top = lineY
		line.rect.Bottom = line.rect.Top + textHeight
		line.rect.Left = lineX
		line.rect.Right = line.rect.Left + textWidth

		lineY += textHeight
	}
	g.Dispose()
}

func (te *TextEdit) GetTextArea() Rect {
	if len(te.lines) < 1 {
		log.Panicln("BUGGG!!!!")
	}
	textArea := te.lines[0].rect
	for i := range te.lines {
		line := te.lines[i]
		rect := &line.rect
		if rect.Left < textArea.Left {
			textArea.Left = rect.Left
		}
		if rect.Right > textArea.Right {
			textArea.Right = rect.Right
		}
		if rect.Top < textArea.Top {
			textArea.Top = rect.Top
		}
		if rect.Bottom > textArea.Bottom {
			textArea.Bottom = rect.Bottom
		}
	}
	return textArea
}

func (te *TextEdit) KeyPressEvent(keycode int) {
	if keycode == win.VK_BACK {
		te.DeleteBack()
	} else {
		te.Insert(te.caratPosition, rune(keycode))
	}
}

func (te *TextEdit) Draw(g *Graphics, color *Color) {
	graphics := gdiplus.NewGraphicsFromHDC(g.GetHDC())
	//area := te.GetTextArea()
	//RenderText(g, "I Am Shuvo", &area, color, te.font)
	brush := gdiplus.NewSolidBrush(asGdiplusColor(color))
	graphics.SetTextRenderingHint(gdiplus.TextRenderingHintAntiAlias)
	// Draw each line
	for i := range te.lines {
		line := te.lines[i]
		lineText := line.chunk.AsString()

		//RenderText(g, lineText, &line.rect, color, te.font)

		textPosition := gdiplus.PointF{X: float32(line.rect.Left), Y: float32(line.rect.Top)}
		lrect := &gdiplus.RectF{X: textPosition.X, Y: textPosition.Y, Width: 0, Height: 0}
		graphics.DrawStringEx(lineText, te.font, lrect, te.format, brush.AsBrush())

		//g.DrawText(lineText, &line.rect, win.DT_LEFT|win.DT_EXPANDTABS, color, te.font.GetHandle())
		// Draw the carat if matches with the postion
		//log.Printf("pos %d, carat %d\n", line.startIndex, te.caratPosition)
		if line.startIndex == te.caratPosition {
			caratY := line.rect.Top
			caratX := line.rect.Left
			g.DrawLineEx(caratX, caratY, caratX, caratY+line.rect.Height(), te.penCarat)
		} else {
			for c := range line.chunk.chars {
				index := line.startIndex + c + 1
				//log.Printf("pos %d, carat %d\n", index, te.caratPosition)
				if index == te.caratPosition {
					caratY := line.rect.Top
					caratX := line.rect.Left
					subText := string(line.chunk.chars[:c+1])

					//subRect := *g.MeasureText(subText, win.DT_LEFT, te.font.GetHandle())

					lrect := &gdiplus.RectF{}
					outRect := &gdiplus.RectF{}
					graphics.MeasureStringEx(subText, te.font, lrect, te.format, outRect, nil, nil)

					//log.Printf("width %f vs %d\n", outRect.Width, subRect.Width())

					caratX += int(outRect.Width + 1) //subRect.Width()
					g.DrawLineEx(caratX, caratY, caratX, caratY+line.rect.Height(), te.penCarat)
					break
				}
			}
		}
	}
	brush.Dispose()
	graphics.Dispose()
}

func RenderText(g *Graphics, text string, rect *Rect, color__ *Color, font *GdiFont) {

	width, height := rect.Width(), rect.Height()

	bg := NewBitmapGraphics(width, height)
	hdc := bg.Hdc
	context := bg.Graphics

	color := Rgb(255, 0, 0)

	context.DrawText(text, rect, win.DT_LEFT|win.DT_EXPANDTABS, &color, font.GetHandle())

	context.DrawLine(0, 0, 50, 50, &color)
	/*
		win.SetBkMode(hdc, win.OPAQUE)
		win.SelectObject(hdc, win.HGDIOBJ(font.GetHandle()))
		win.SetTextColor(hdc, Rgb(255, 0, 0).AsCOLORREF())
		utf16, _ := syscall.UTF16FromString(text)
		ExtTextOutW(hdc, int32(rect.Left), int32(rect.Top), 0, nil, &utf16[0], uint32(len(utf16)), nil)
	*/

	pixel := bg.Data
	pixelCount := width * height

	index := 0
	for c := 0; c != pixelCount; c++ {
		// Set alpha
		alpha := 255 - pixel[0]
		pixel[index+3] = alpha
		// Set colour
		pixel[index+0] = (color.B * alpha) / 255
		pixel[index+1] = (color.G * alpha) / 255
		pixel[index+2] = (color.R * alpha) / 255
		index += 4
	}

	//g.BitBlt(rect.Left, rect.Top, width, height, hdc, 0, 0, win.SRCCOPY)
	/*
		utf16, _ := syscall.UTF16FromString(text)
		win.SetBkMode(hdc, win.TRANSPARENT)
		win.SelectObject(hdc, win.HGDIOBJ(font.GetHandle()))
		win.SetTextColor(hdc, color.AsCOLORREF())

		ExtTextOutW(hdc, int32(rect.Left), int32(rect.Top), 0, nil, &utf16[0], uint32(len(utf16)), nil)
	*/

	var bf win.BLENDFUNCTION
	bf.BlendOp = AC_SRC_OVER
	bf.BlendFlags = 0
	bf.SourceConstantAlpha = 255
	bf.AlphaFormat = win.AC_SRC_ALPHA
	win.AlphaBlend(g.GetHDC(), int32(rect.Left), int32(rect.Top), int32(width), int32(height), hdc, 0, 0, int32(width), int32(height), bf)

	/*
		hdc := bg.Hdc
		utf16, _ := syscall.UTF16FromString(text)
		win.SetBkMode(hdc, win.TRANSPARENT)
		win.SelectObject(hdc, win.HGDIOBJ(font.GetHandle()))
		win.SetTextColor(hdc, color.AsCOLORREF())

		ExtTextOutW(hdc, int32(rect.Left), int32(rect.Top), 0, nil, &utf16[0], uint32(len(utf16)), nil)
	*/
	bg.Dispose()
}
