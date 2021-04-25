package main

import (
	"github.com/shahfarhadreza/go-gdiplus"
)

type TriangleDrawer struct {
	ShapeDrawer
	points [3]gdiplus.Point
}

func (tool *TriangleDrawer) draw(args *ToolDrawShapeArgs) {
	g := args.context
	pen := args.pen
	brush := args.brush
	startPoint := args.startPoint
	endPoint := args.endPoint
	tool.points[0] = gdiplus.Point{X: startPoint.X, Y: endPoint.Y}
	tool.points[1] = gdiplus.Point{X: endPoint.X, Y: endPoint.Y}
	tool.points[2] = gdiplus.Point{X: startPoint.X + ((endPoint.X - startPoint.X) / 2), Y: startPoint.Y}
	if brush != nil {
		g.FillPolygonI(brush, tool.points[:], gdiplus.FillModeAlternate)
	}
	if pen != nil {
		g.DrawPolygonI(pen, tool.points[:])
	}
}
