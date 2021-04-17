package main

import (
	"gopaint/gdiplus"
)

type RoundRectDrawer struct {
	ShapeDrawer
}

func NewRoundedRectPath(bounds *gdiplus.Rect, radius int32) *gdiplus.GraphicsPath {
	path := gdiplus.NewPath(gdiplus.FillModeAlternate)
	diameter := radius * 2
	arc := gdiplus.NewRect(bounds.X, bounds.Y, diameter, diameter)
	// top left arc
	path.AddArcRect(arc, 180, 90)
	// top right arc
	arc.X = bounds.Right() - diameter
	path.AddArcRect(arc, 270, 90)
	// bottom right arc
	arc.Y = bounds.Bottom() - diameter
	path.AddArcRect(arc, 0, 90)
	// bottom left arc
	arc.X = bounds.X
	path.AddArcRect(arc, 90, 90)
	path.CloseFigure()
	return path
}

func (tool *RoundRectDrawer) draw(args *ToolDrawShapeArgs) {
	g := args.context
	pen := args.pen
	brush := args.brush
	rect := args.rect
	radius := int32(10)
	// oh god this hardcoded fix.....
	if rect.Width < 40 || rect.Height < 40 {
		radius = 8
	}
	if rect.Width < 30 || rect.Height < 30 {
		radius = 5
	}
	if rect.Width < 10 || rect.Height < 10 {
		radius = 1
	}
	round := NewRoundedRectPath(&rect, radius)
	if brush != nil {
		g.FillPath(brush, round)
	}
	if pen != nil {
		g.DrawPath(pen, round)
	}
	round.Dispose()
}
