package main

import (
	. "gopaint/reza"
	"log"
	"strconv"
)

type ResizeDialog struct {
	Dialog
}

type PropertiesDialog struct {
	Dialog
	lastSaved  Label
	sizeOnDisk Label
	width      TextBox
	height     TextBox
}

func NewResizeDialog(parent Window) *ResizeDialog {
	dlg := &ResizeDialog{Dialog: NewDialog()}
	dlg.Init(parent)
	return dlg
}

func (dlg *ResizeDialog) Init(parent Window) {
	logInfo("Initialize Resize dialog...")
	dlg.Dialog.Initialize(parent, "Resize", 300, 320)

	dlg.AddWidgets([]Widget{
		&WGroup{Text: "Resize", DockType: DockFill,
			Margins: Margins{Left: 10, Top: 10, Right: 10, Bottom: 10}, Widgets: []Widget{
				&WFlowContainer{FlowDirection: FlowLeftToRight, DockType: DockFill,
					Margins: Margins{Left: 20, Top: 30, Right: 10, Bottom: 10}, Widgets: []Widget{
						&WLabel{Text: "By:\t"},
						&WRadioButton{Text: "Parcentage", Margins: Margins{Right: 20}},
						&WRadioButton{Text: "Pixels", Margins: Margins{Right: 20, Bottom: 20}, Checked: true},

						&WImageViewer{Path: ".\\icons\\horizintal.png", Margins: Margins{Right: 30, Bottom: 20}},
						&WLabel{Text: "Horizintal:\t", Margins: Margins{Top: 5}},
						&WTextBox{Text: "100", Width: 60, Height: 24, Margins: Margins{Right: 20}},

						&WImageViewer{Path: ".\\icons\\vertical.png", Margins: Margins{Right: 30, Bottom: 20}},
						&WLabel{Text: "Vertical:\t\t", Margins: Margins{Top: 5}},
						&WTextBox{Text: "100", Width: 60, Height: 24, Margins: Margins{Right: 20}},

						&WCheckButton{Text: "Maintain aspect ratio", Checked: true},
					}},
			}},
	})
}

func NewPropertiesDialog(parent Window) *PropertiesDialog {
	dlg := &PropertiesDialog{Dialog: NewDialog()}
	dlg.Init(parent)
	return dlg
}

func (dlg *PropertiesDialog) Init(parent Window) {
	logInfo("Initialize Properties dialog...")
	dlg.Dialog.Initialize(parent, "Image Properties", 340, 360)

	dlg.AddWidgets([]Widget{
		&WGroup{Text: "File Attributes", DockType: DockTop, Height: 100,
			Margins: Margins{Left: 10, Top: 10, Right: 10, Bottom: 10}, Widgets: []Widget{
				&WFlowContainer{FlowDirection: FlowLeftToRight, DockType: DockFill,
					Margins: Margins{Left: 15, Top: 15, Right: 5, Bottom: 5}, Widgets: []Widget{
						&WLabel{Text: "Last Saved:\tNot Available", Margins: Margins{Top: 10, Right: 20}, AssignTo: &dlg.lastSaved},
						&WLabel{Text: "Size on disk:\tNot Available", Margins: Margins{Top: 10, Right: 40}, AssignTo: &dlg.sizeOnDisk},
						&WLabel{Text: "Resolution:\t96 DPI", Margins: Margins{Top: 10}},
					}},
			}},
		&WFlowContainer{FlowDirection: FlowLeftToRight, DockType: DockBottom, Height: 25,
			Margins: Margins{Left: 10, Right: 10, Bottom: 10}, Widgets: []Widget{
				&WLabel{Text: "Width:\t", Margins: Margins{Top: 4}},
				&WTextBox{Text: "1152", Width: 60, Height: 24, Margins: Margins{Right: 10}, AssignTo: &dlg.width},
				&WLabel{Text: "Height:\t", Margins: Margins{Top: 4}},
				&WTextBox{Text: "648", Width: 60, Height: 24, Margins: Margins{Right: 10}, AssignTo: &dlg.height},
				&WButton{Text: "Default", Width: 70, Height: 24, OnClick: func(sender Button) {
					dlg.width.SetText("1152")
					dlg.height.SetText("648")
				}},
			}},
		&WFlowContainer{FlowDirection: FlowNone, DockType: DockFill,
			Margins: Margins{Left: 10, Right: 10, Bottom: 10}, Widgets: []Widget{
				&WGroup{Text: "Units", DockType: DockLeft, Width: 140, Margins: Margins{Right: 5}, Widgets: []Widget{
					&WFlowContainer{FlowDirection: FlowLeftToRight, DockType: DockFill,
						Margins: Margins{Left: 5, Top: 20, Right: 10, Bottom: 10}, Widgets: []Widget{
							&WRadioButton{Text: "Inches", Margins: Margins{Left: 10, Top: 10}},
							&WRadioButton{Text: "Centimeters", Margins: Margins{Left: 10, Top: 10}},
							&WRadioButton{Text: "Pixels", Margins: Margins{Left: 10, Top: 10}, Checked: true},
						}},
				}},
				&WGroup{Text: "Colors", DockType: DockFill, Margins: Margins{Left: 5}, Widgets: []Widget{
					&WFlowContainer{FlowDirection: FlowLeftToRight, DockType: DockFill,
						Margins: Margins{Left: 5, Top: 20, Right: 10, Bottom: 10}, Widgets: []Widget{
							&WRadioButton{Text: "Black and white", Margins: Margins{Left: 10, Top: 10}},
							&WRadioButton{Text: "Color", Margins: Margins{Left: 10, Top: 10}, Checked: true},
						}},
				}},
			}},
	})
}

func (dlg *PropertiesDialog) Show() {
	canvas := mainWindow.workspace.canvas
	image := canvas.image
	if imageSize, available := image.SizeOnDisk(); available {
		dlg.sizeOnDisk.SetText("Size on disk:\t" + imageSize)
	} else {
		dlg.sizeOnDisk.SetText("Size on disk:\tNot Available")
	}
	if lastSaved, available := image.LastSaved(); available {
		dlg.lastSaved.SetText("Last Saved:\t" + lastSaved)
	} else {
		dlg.lastSaved.SetText("Last Saved:\tNot Available")
	}
	dlg.width.SetText(strconv.Itoa(image.Width()))
	dlg.height.SetText(strconv.Itoa(image.Height()))
	dlg.Dialog.Show(true, func() {
		width := dlg.width.GetText()
		height := dlg.height.GetText()
		newCanvasWidth, err := strconv.Atoi(width)
		if err == nil {
			newCanvasHeight, err := strconv.Atoi(height)
			if err == nil {
				if newCanvasHeight > 0 && newCanvasWidth > 0 {
					canvas.Resize(newCanvasWidth, newCanvasHeight)
					mainWindow.workspace.RequestLayout()
				} else {
					log.Panicln("Enter valid canvas size!!!")
				}
			} else {
				log.Panicln(err)
			}
		} else {
			log.Panicln(err)
		}
	})
}
