package main

import "log"

type ToolsManager struct {
	toolSelect    *ToolSelect
	toolPencil    *ToolPencil
	toolBucket    *ToolBucket
	toolText      *ToolText
	toolEraser    *ToolEraser
	toolPickColor *ToolPickColor
	toolBrush     *ToolBrush
	// Shape tools
	toolShapeLine      *ToolShape
	toolShapeRect      *ToolShape
	toolShapeRoundRect *ToolShape
	toolShapeEllipse   *ToolShape
	toolShapeTriangle  *ToolShape
	toolShapeDiamond   *ToolShape
	currentTool        Tool
}

func NewToolsManager() *ToolsManager {
	mgr := &ToolsManager{}
	mgr.init()
	return mgr
}

func (tools *ToolsManager) init() {
	tools.toolPencil = &ToolPencil{}
	tools.toolPencil.initialize()
	tools.toolBrush = &ToolBrush{}
	tools.toolBrush.initialize()
	tools.toolPickColor = &ToolPickColor{}
	tools.toolPickColor.initialize()
	tools.toolEraser = &ToolEraser{}
	tools.toolEraser.initialize()
	tools.toolBucket = &ToolBucket{}
	tools.toolBucket.initialize()
	tools.toolText = &ToolText{}
	tools.toolText.initialize()
	tools.toolSelect = &ToolSelect{}
	tools.toolSelect.initialize()
	// Shape tools
	tools.toolShapeLine = &ToolShape{ShapeDrawer: &LineDrawer{}}
	tools.toolShapeLine.initialize()
	tools.toolShapeRect = &ToolShape{ShapeDrawer: &RectangleDrawer{}}
	tools.toolShapeRect.initialize()
	tools.toolShapeRoundRect = &ToolShape{ShapeDrawer: &RoundRectDrawer{}}
	tools.toolShapeRoundRect.initialize()
	tools.toolShapeEllipse = &ToolShape{ShapeDrawer: &EllipseDrawer{}}
	tools.toolShapeEllipse.initialize()
	tools.toolShapeTriangle = &ToolShape{ShapeDrawer: &TriangleDrawer{}}
	tools.toolShapeTriangle.initialize()
	tools.toolShapeDiamond = &ToolShape{ShapeDrawer: &DiamondDrawer{}}
	tools.toolShapeDiamond.initialize()
}

func (tools *ToolsManager) Dispose() {
	logInfo("Disposing ToolsManager...")
	tools.toolPencil.Dispose()
	tools.toolBrush.Dispose()
	tools.toolPickColor.Dispose()
	tools.toolEraser.Dispose()
	tools.toolBucket.Dispose()
	tools.toolText.Dispose()
	tools.toolSelect.Dispose()
	// Shape tools
	tools.toolShapeLine.Dispose()
	tools.toolShapeRect.Dispose()
	tools.toolShapeRoundRect.Dispose()
	tools.toolShapeEllipse.Dispose()
	tools.toolShapeTriangle.Dispose()
	tools.toolShapeDiamond.Dispose()
}

func (tools *ToolsManager) SetCurrentTool(tool Tool) {
	if tool == nil {
		log.Panicln("tool is nil !!!!")
	}
	if tools.currentTool == tool {
		return
	}
	if tools.currentTool != nil {
		tools.currentTool.leave()
	}
	tools.currentTool = tool
	tool.prepare()
}

func (tools *ToolsManager) GetCurrentTool() Tool {
	return tools.currentTool
}
