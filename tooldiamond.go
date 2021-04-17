package main

import (
	"gopaint/gdiplus"
)

type DiamondDrawer struct {
	ShapeDrawer
	points [4]gdiplus.Point
}

func (tool *DiamondDrawer) draw(args *ToolDrawShapeArgs) {
	g := args.context
	pen := args.pen
	brush := args.brush
	startPoint := args.startPoint
	endPoint := args.endPoint
	halfWidth, halfHeight := ((endPoint.X - startPoint.X) / 2), ((endPoint.Y - startPoint.Y) / 2)
	tool.points[0] = gdiplus.Point{X: startPoint.X + halfWidth, Y: startPoint.Y}
	tool.points[1] = gdiplus.Point{X: startPoint.X, Y: startPoint.Y + halfHeight}
	tool.points[2] = gdiplus.Point{X: startPoint.X + halfWidth, Y: endPoint.Y}
	tool.points[3] = gdiplus.Point{X: endPoint.X, Y: startPoint.Y + halfHeight}
	if brush != nil {
		g.FillPolygonI(brush, tool.points[:], gdiplus.FillModeAlternate)
	}
	if pen != nil {
		g.DrawPolygonI(pen, tool.points[:])
	}
}
