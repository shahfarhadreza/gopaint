package reza

import (
	win "github.com/lxn/win"
)

const FlowNone = 0
const FlowLeftToRight = 1
const FlowRightToLeft = 2

// A Flow container
type Container interface {
	// Promote Window interface methods
	Window
	SetFlowDirection(direction int)
	GetFlowDirection() int
}

type ContainerData struct {
	// Embed Window interface
	Window
	flowDirection int
}

func CreateContainer(flowDirection int, x, y, width, height int, parent Window) Container {
	c := &ContainerData{Window: NewWindow()}
	c.Create("", win.WS_CHILD|win.WS_CLIPCHILDREN|win.WS_VISIBLE, x, y, width, height, parent)
	/*
		c.SetPaintEventHandler(func(g *Graphics, rect *Rect) {
			g.DrawFillRectangle(rect, Rgb(255, 0, 0), Rgb(70, 70, 70))
		})*/
	c.SetResizeEventHandler(func(client *Rect) {
		c.UpdateLayout(client)
	})
	c.SetFlowDirection(flowDirection)
	return c
}

func (c *ContainerData) SetFlowDirection(direction int) {
	c.flowDirection = direction
}
func (c *ContainerData) GetFlowDirection() int {
	return c.flowDirection
}

func (c *ContainerData) UpdateLayout(client *Rect) {
	// Layout the childs
	if c.flowDirection == FlowRightToLeft {
		startX := client.Width()
		startY := 0
		for _, item := range c.GetChildrens() {
			if item.GetDockType() != DockNone {
				continue
			}
			_, marginRight, marginTop, _ := item.GetMargin()
			startY = marginTop
			startX -= marginRight
			itemSize := item.GetSize()
			newX := startX - itemSize.Width
			if newX < 0 {
				startY += itemSize.Height
			}
			item.SetPosition(newX, startY)
			startX -= itemSize.Width
		}
	} else if c.flowDirection == FlowLeftToRight {
		startX := 0
		startY := 0
		maxHeight := 0
		for _, item := range c.GetChildrens() {
			if item.GetDockType() != DockNone {
				continue
			}
			marginLeft, marginRight, marginTop, marginBottom := item.GetMargin()
			startX += marginLeft
			itemSize := item.GetSize()
			if (startX + itemSize.Width) > client.Width() {
				startX = marginLeft
				startY += maxHeight
			}
			itemTotalHeight := marginTop + itemSize.Height + marginBottom
			if itemTotalHeight > maxHeight {
				maxHeight = itemTotalHeight
			}
			item.SetPosition(startX, startY+marginTop)
			startX += itemSize.Width + marginRight
		}
	}
}
