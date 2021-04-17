package main

type RectangleDrawer struct {
	ShapeDrawer
}

func (tool *RectangleDrawer) draw(args *ToolDrawShapeArgs) {
	g := args.context
	pen := args.pen
	brush := args.brush
	rect := args.rect
	if brush != nil {
		g.FillRectangleI(brush, int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height))
	}
	if pen != nil {
		g.DrawRectangleI(pen, int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height))
	}
}
