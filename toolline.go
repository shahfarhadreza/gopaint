package main

type LineDrawer struct {
	ShapeDrawer
}

func (tool *LineDrawer) draw(args *ToolDrawShapeArgs) {
	g := args.context
	pen := args.pen
	if pen != nil {
		g.DrawLineI(pen, args.startPointOrg.X, args.startPointOrg.Y, args.endPointOrg.X, args.endPointOrg.Y)
	}
}
