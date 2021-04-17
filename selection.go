package main

import (
	. "gopaint/reza"
)

type SelectionRect struct {
	rect         Rect
	penWhite     *Pen
	penBorder    *Pen
	penResizer   *Pen
	brushResizer *Brush
}

func NewSelectionRect() *SelectionRect {
	sr := &SelectionRect{}
	sr.penWhite = NewSolidPen(1, NewRgb(255, 255, 255))
	sr.penBorder = NewUserStylePen(1, NewRgb(0, 120, 215), []uint32{3, 4})
	sr.penResizer = NewSolidPen(1, NewRgb(85, 85, 85))
	sr.brushResizer = NewSolidBrush(NewRgb(255, 255, 255))
	return sr
}

func (sr *SelectionRect) Dispose() {
	sr.penBorder.Dispose()
	sr.brushResizer.Dispose()
	sr.penResizer.Dispose()
	sr.penWhite.Dispose()
}

func (sr *SelectionRect) Clear() {
	sr.rect = Rect{}
}

func (sr *SelectionRect) SetRect(rect *Rect) {
	sr.rect = *rect
}

func (sr *SelectionRect) GetRect() Rect {
	return sr.rect
}

func (sr *SelectionRect) IsEmpty() bool {
	return !(sr.rect.Width() >= 1 && sr.rect.Height() >= 1)
}

func (sr *SelectionRect) GetClosestRectPoint(mouse *Point, dist int) (onPoint bool, rectPoint int) {
	points := sr.rect.GetEightPoints()
	for i := range points {
		point := &points[i]
		pdist := point.DistanceI(mouse)
		if pdist < dist {
			return true, i
		}
	}
	return false, 0
}

func (sr *SelectionRect) Draw(g *Graphics) {
	points := sr.rect.GetEightPoints()
	const halfSize = 3

	g.DrawRectangleEx(&sr.rect, sr.penWhite, nil)
	g.DrawRectangleEx(&sr.rect, sr.penBorder, nil)

	for i := range points {
		point := &points[i]
		rcHandle := Rect{
			Left:   point.X - halfSize,
			Right:  point.X + halfSize,
			Top:    point.Y - halfSize,
			Bottom: point.Y + halfSize,
		}
		g.DrawRectangleEx(&rcHandle, sr.penResizer, sr.brushResizer)
	}
}
