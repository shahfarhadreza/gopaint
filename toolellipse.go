package main

type EllipseDrawer struct {
	ShapeDrawer
}

func (tool *EllipseDrawer) draw(args *ToolDrawShapeArgs) {
	g := args.context
	pen := args.pen
	brush := args.brush
	rect := args.rect
	if brush != nil {
		g.FillEllipseI(brush, int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height))
	}
	if pen != nil {
		g.DrawEllipseI(pen, int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height))
	}
}
