package main

import (
	. "gopaint/reza"
	"log"
	"path/filepath"

	win "github.com/lxn/win"
)

// MainWindow is the main frame window of GoPaint
type MainWindow struct {
	// Embed the Form interface
	Form
	// Own data
	workspaceColor   Color
	hCursorArrow     win.HCURSOR
	hCursorSizeNS    win.HCURSOR
	hCursorSizeWE    win.HCURSOR
	hCursorSizeNWSE  win.HCURSOR
	hCursorSizeNESW  win.HCURSOR
	hCursorIBeam     win.HCURSOR
	hCursorMove      win.HCURSOR
	ribbon           Ribbon
	workspace        *Workspace
	statusbar        Statusbar
	statusMousePos   Status
	statusSelSize    Status
	statusCanvasSize Status
	statusFileSize   Status
	color1           RibbonButton
	color2           RibbonButton
	bsize            RibbonButton
	bsizeMenu        PopupSizeMenu
	btnTools         []RibbonButton
	buttonToolPairs  map[RibbonButton]Tool
	menuNoOutline    PopupMenuItem
	menuSolidOutline PopupMenuItem
	menuNoFill       PopupMenuItem
	menuSolidFill    PopupMenuItem
	bShowGridlines   RibbonButton
	tools            *ToolsManager
	resizeDialog     *ResizeDialog
	propertiesDialog *PropertiesDialog
	initDone         bool
}

const appWidth = 1300
const appHeight = 840
const newImageName = "Untitled"

var mainWindow *MainWindow

func NewMainWindow() *MainWindow {
	mainWindow = &MainWindow{Form: NewForm()}
	return mainWindow
}

func (window *MainWindow) Initialize() {
	logInfo("initializing main window...")
	window.initDone = false
	window.workspaceColor = Rgb(199, 208, 224)
	window.SetText(newImageName + " - " + app.Title)
	window.SetPosition(310, 100)
	window.SetSize(appWidth, appHeight)
	window.SetPaintEventHandler(func(g *Graphics, rect *Rect) {
		g.FillRect(rect, &window.workspaceColor)
	})
	window.SetDestroyEventHandler(func() {
		app.Exit()
	})
	window.initResources()

	window.tools = NewToolsManager()

	window.initRibbon()

	statusbar := NewStatusbar(window)
	statusbar.SetDockType(DockBottom)
	statusbar.SetSize(0, 27)
	window.statusMousePos = statusbar.AddStatus(".\\icons\\mouse-position.png", "1, 3px")
	window.statusSelSize = statusbar.AddStatus(".\\icons\\selection-size.png", "10 x 10px")
	window.statusCanvasSize = statusbar.AddStatus(".\\icons\\canvas-size.png", "800 x 600px")
	window.statusFileSize = statusbar.AddStatus(".\\icons\\file-size.png", "0.2KB")
	window.statusbar = statusbar

	window.workspace = NewWorkspace(window)
	window.workspace.SetDockType(DockFill)

	window.SetCurrentTool(window.tools.toolBrush)

	window.resizeDialog = NewResizeDialog(window)
	window.propertiesDialog = NewPropertiesDialog(window)

	logInfo("Done initializing main window")
	window.initDone = true
	window.RequestLayout()

	//TestNewDialog(window)
	//log.Println(fileSaveDialog.IUnknown)
}

func (window *MainWindow) Dispose() {
	logInfo("MainWindow Dispose")
	if window.tools != nil {
		window.tools.Dispose()
	}
	if window.workspace != nil {
		window.workspace.Dispose()
	}
	if window.ribbon != nil {
		window.ribbon.Dispose()
	}
}

func (window *MainWindow) initResources() {
	// Load some resources
	window.hCursorArrow = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	window.hCursorSizeNS = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_SIZENS))
	window.hCursorSizeWE = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_SIZEWE))
	window.hCursorSizeNWSE = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_SIZENWSE))
	window.hCursorSizeNESW = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_SIZENESW))
	window.hCursorIBeam = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_IBEAM))
	window.hCursorMove = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_SIZEALL))
}

func (window *MainWindow) initRibbonApplicationMenu() {
	ribbon := window.ribbon

	fResetColors := func() {
		// Reset the background and foreground colors to default values
		window.color1.SetColor(Rgb(0, 0, 0))
		window.color2.SetColor(Rgb(255, 255, 255))
	}

	funcNew := func(e *PopupItemEvent) {
		workspace := window.workspace
		workspace.canvas.NewImage(600, 480)
		workspace.RequestLayout()
		window.SetText(newImageName + " - " + app.Title)
		fResetColors()
	}

	funcOpen := func(e *PopupItemEvent) {
		filter := GetOpenFileDialogFilters()
		filename, accepted := OpenFileDialog(window,
			filter,
			GetFormatCount()+1)
		if accepted {
			workspace := window.workspace
			if workspace.canvas.OpenImage(filename) {
				workspace.RequestLayout()
				fResetColors()
				fileNameOnly := filepath.Base(filename)
				window.SetText(fileNameOnly + " - " + app.Title)
			}
		}
	}

	fSaveAs := func(filename string) {
		ext := filepath.Ext(filename)
		filterIndex := FindFormatIndexFromExt(ext) + 1 // dialog filter indices are 1 based
		filter := GetSaveFileDialogFilters()
		newfilepath, accepted := SaveFileDialog(window, filename, filter, filterIndex)
		if accepted {
			workspace := window.workspace
			workspace.canvas.SaveImage(newfilepath)
			workspace.Repaint()
			fileNameOnly := filepath.Base(newfilepath)
			window.SetText(fileNameOnly + " - " + app.Title)
		}
	}

	fMenuSave := func(e *PopupItemEvent) {
		workspace := window.workspace
		canvas := workspace.canvas
		image := canvas.image
		// If we already have a path
		if image.HasFilePath() {
			// We just simply overwrite the previous file
			canvas.SaveImage(image.filepath)
		} else {
			fSaveAs("Untitled.png")
		}
	}

	fMenuSaveAs := func(e *PopupItemEvent) {
		workspace := window.workspace
		canvas := workspace.canvas
		image := canvas.image
		if len(image.filepath) > 0 {
			// Make sure we only pass the base file name
			filename := filepath.Base(image.filepath)
			fSaveAs(filename)
		} else {
			fSaveAs("Untitled.png")
		}
	}

	appMenu := NewPopupMenu(ribbon, []MenuItemInfo{
		{Text: "New", IconPath: ".\\icons\\big-new.png", OnClick: funcNew},
		{Text: "Open", IconPath: ".\\icons\\big-open.png", OnClick: funcOpen},
		{Text: "Save", IconPath: ".\\icons\\big-save.png", OnClick: fMenuSave},
		{Text: "Save as", IconPath: ".\\icons\\big-save-as.png", OnClick: fMenuSaveAs},
		//{Sperator: true},
		//{Text: "Print", IconPath: ".\\icons\\big-print.png"},
		//{Sperator: true},
		//{Text: "Set as desktop background", IconPath: ".\\icons\\big-set-as-desktop.png"},
		{Sperator: true},
		{Text: "Properties", IconPath: ".\\icons\\big-properties.png", OnClick: func(e *PopupItemEvent) {
			window.propertiesDialog.Show()
		}},
		{Sperator: true},
		{Text: "Exit", IconPath: ".\\icons\\big-exit.png",
			OnClick: func(e *PopupItemEvent) {
				app.Exit()
			}},
	})
	appMenu.SetLargeItem(true)
	ribbon.SetApplicationMenu("File", appMenu)
}

func (window *MainWindow) initRibbon() {
	logInfo("init Ribbon....")
	window.ribbon = NewRibbon(window)
	ribbon := window.ribbon
	ribbon.SetDockType(DockTop)
	ribbon.SetSize(0, 118)
	ribbon.SuspendRepaint()

	window.btnTools = make([]RibbonButton, 0)

	window.initRibbonApplicationMenu()

	home := ribbon.AddTab("Home")

	clipboard := home.AddSection("Clipboard")

	bpaste := clipboard.AddImageButton("Paste", ".\\icons\\paste.png", RibbonButtonSizeBig)

	var mpaste, mpasteFrom PopupMenuItem

	bpasteMenu := NewPopupMenu(ribbon, []MenuItemInfo{
		{Text: "Paste", IconPath: ".\\icons\\paste-small.png", AssignTo: &mpaste},
		{Text: "Paste from", IconPath: ".\\icons\\paste-from-small.png", AssignTo: &mpasteFrom},
	})
	bpaste.SetDropdownMenu(bpasteMenu, true)
	bpaste.SetEnabled(false)
	mpaste.SetEnabled(false)
	mpasteFrom.SetEnabled(false)

	bcut := clipboard.AddImageButton("Cut", ".\\icons\\cut.png", RibbonButtonSizeMedium)
	bcopy := clipboard.AddImageButton("Copy", ".\\icons\\copy.png", RibbonButtonSizeMedium)

	bcut.SetEnabled(false)
	bcopy.SetEnabled(false)

	imagesec := home.AddSection("Image")

	selectIcon, _ := CreateBitmapImage(".\\icons\\select.png", false)
	lassoIcon, _ := CreateBitmapImage(".\\icons\\select-lasso.png", false)

	bselect := imagesec.AddImageButton("Select", ".\\icons\\select.png", RibbonButtonSizeBig)
	window.btnTools = append(window.btnTools, bselect)

	bselect.SetIcon(selectIcon)

	var regularSel, lassoSel PopupMenuItem
	bselectMenu := NewPopupMenu(ribbon, []MenuItemInfo{
		{Text: "Selection shapes", Sperator: true},
		{Text: "Rectangular selection", IconPath: ".\\icons\\select-small.png", AssignTo: &regularSel,
			OnClick: func(e *PopupItemEvent) {
				bselect.SetIcon(selectIcon)
				regularSel.SetToggled(true)
				lassoSel.SetToggled(false)
				window.SetCurrentTool(window.tools.toolSelect)
			}},
		{Text: "Free-form selection", IconPath: ".\\icons\\select-lasso-small.png", AssignTo: &lassoSel,
			OnClick: func(e *PopupItemEvent) {
				bselect.SetIcon(lassoIcon)
				regularSel.SetToggled(false)
				lassoSel.SetToggled(true)
				window.SetCurrentTool(window.tools.toolSelect)
			}},
		{Text: "Selection options", Sperator: true},
		{Text: "Select all", IconPath: ".\\icons\\select-all-small.png",
			OnClick: func(e *PopupItemEvent) {
				window.SetCurrentTool(window.tools.toolSelect)
				window.tools.toolSelect.SelectAll()
			}},
		{Text: "Deselect", OnClick: func(e *PopupItemEvent) {
			window.SetCurrentTool(window.tools.toolSelect)
			window.tools.toolSelect.Deselect()
		}},
		{Text: "Invert selection", IconPath: ".\\icons\\select-invert-small.png"},
		{Text: "Delete", IconPath: ".\\icons\\delete-small.png",
			OnClick: func(e *PopupItemEvent) {
				window.SetCurrentTool(window.tools.toolSelect)
				window.tools.toolSelect.DeleteSelection()
			}},
		//{Text: "Transparent selection"},
	})
	bselect.SetDropdownMenu(bselectMenu, true)
	regularSel.SetToggled(true)

	bcrop := imagesec.AddImageButton("Crop", ".\\icons\\crop.png", RibbonButtonSizeMedium)
	bcrop.SetEnabled(false)

	imagesec.AddImageButton("Resize", ".\\icons\\resize.png", RibbonButtonSizeMedium).SetClickEvent(func(e *RibbonButtonEvent) {
		window.resizeDialog.Show(true, func() {
			log.Println("Implement resizing!")
		})
	})

	brotate := imagesec.AddImageButton("Rotate", ".\\icons\\rotate.png", RibbonButtonSizeMedium)

	brotateMenu := NewPopupMenu(ribbon, []MenuItemInfo{
		{Text: "Rotate right 90", IconPath: ".\\icons\\rotate-right-small.png"},
		{Text: "Rotate left 90", IconPath: ".\\icons\\rotate-left-small.png"},
		{Text: "Rotate 180", IconPath: ".\\icons\\rotate-180-small.png"},
		{Text: "Flip vertical", IconPath: ".\\icons\\rotate-v-small.png"},
		{Text: "Flip horizontal", IconPath: ".\\icons\\rotate-h-small.png"},
	})
	brotate.SetDropdownMenu(brotateMenu, false)

	tools := home.AddSection("Tools")
	tools.SetTwoRow(true)

	bpencil := tools.AddImageButton("Pencil", ".\\icons\\pencil.png", RibbonButtonSizeSmall)
	beraser := tools.AddImageButton("Eraser", ".\\icons\\eraser.png", RibbonButtonSizeSmall)
	bbucket := tools.AddImageButton("Fill with color", ".\\icons\\fill.png", RibbonButtonSizeSmall)
	bpickcolor := tools.AddImageButton("Color picker", ".\\icons\\pick.png", RibbonButtonSizeSmall)
	btext := tools.AddImageButton("Text", ".\\icons\\text.png", RibbonButtonSizeSmall)
	//bzoom := tools.AddImageButton("Magnifier", ".\\icons\\zoom.png", RibbonButtonSizeSmall)

	window.btnTools = append(window.btnTools, bpencil)
	window.btnTools = append(window.btnTools, beraser)
	window.btnTools = append(window.btnTools, bbucket)
	window.btnTools = append(window.btnTools, bpickcolor)
	window.btnTools = append(window.btnTools, btext)
	//window.btnTools = append(window.btnTools, bzoom)

	brushsec := home.AddSection("")

	bbrush := brushsec.AddImageButton("Brush", ".\\icons\\brush.png", RibbonButtonSizeBig)
	window.btnTools = append(window.btnTools, bbrush)

	shapes := home.AddSection("Shapes")

	bline := shapes.AddImageButton("Line", ".\\icons\\shape-line.png", RibbonButtonSizeSmall)
	brrect := shapes.AddImageButton("Round Rectangle", ".\\icons\\shape-round-rect.png", RibbonButtonSizeSmall)
	//bhexagon := shapes.AddImageButton("Hexagon", ".\\icons\\shape-hexagon.png", RibbonButtonSizeSmall)

	bellipse := shapes.AddImageButton("Ellipse", ".\\icons\\shape-ellipse.png", RibbonButtonSizeSmall)
	//bpentagon := shapes.AddImageButton("Pentagon", ".\\icons\\shape-pentagon.png", RibbonButtonSizeSmall)
	btriangle := shapes.AddImageButton("Triangle", ".\\icons\\shape-triangle.png", RibbonButtonSizeSmall)

	brect := shapes.AddImageButton("Rectangle", ".\\icons\\shape-rectangle.png", RibbonButtonSizeSmall)
	bdiamond := shapes.AddImageButton("Diamond", ".\\icons\\shape-diamond.png", RibbonButtonSizeSmall)
	//brtriangle := shapes.AddImageButton("Right triangle", ".\\icons\\shape-triangle-right.png", RibbonButtonSizeSmall)

	window.btnTools = append(window.btnTools, bline)
	window.btnTools = append(window.btnTools, brrect)
	//window.btnTools = append(window.btnTools, bhexagon)
	window.btnTools = append(window.btnTools, bellipse)
	//window.btnTools = append(window.btnTools, bpentagon)
	window.btnTools = append(window.btnTools, btriangle)
	window.btnTools = append(window.btnTools, brect)
	window.btnTools = append(window.btnTools, bdiamond)
	//window.btnTools = append(window.btnTools, brtriangle)

	window.buttonToolPairs = make(map[RibbonButton]Tool)

	for _, btn := range window.btnTools {
		switch btn {
		case bpencil:
			window.buttonToolPairs[btn] = window.tools.toolPencil
		case bbrush:
			window.buttonToolPairs[btn] = window.tools.toolBrush
		case bpickcolor:
			window.buttonToolPairs[btn] = window.tools.toolPickColor
		case beraser:
			window.buttonToolPairs[btn] = window.tools.toolEraser
		case bbucket:
			window.buttonToolPairs[btn] = window.tools.toolBucket
		case btext:
			window.buttonToolPairs[btn] = window.tools.toolText
		case bselect:
			window.buttonToolPairs[btn] = window.tools.toolSelect
		// Shapes
		case bline:
			window.buttonToolPairs[btn] = window.tools.toolShapeLine
		case brect:
			window.buttonToolPairs[btn] = window.tools.toolShapeRect
		case brrect:
			window.buttonToolPairs[btn] = window.tools.toolShapeRoundRect
		case bellipse:
			window.buttonToolPairs[btn] = window.tools.toolShapeEllipse
		case btriangle:
			window.buttonToolPairs[btn] = window.tools.toolShapeTriangle
		case bdiamond:
			window.buttonToolPairs[btn] = window.tools.toolShapeDiamond
		}
	}

	fbtntools := func(e *RibbonButtonEvent) {
		for _, btn := range window.btnTools {
			btn.SetToggled(false)
		}
		e.Button.SetToggled(true)
		// it will get enabled by the tool if required
		window.bsize.SetEnabled(false)
		tool, found := window.buttonToolPairs[e.Button]
		if found {
			window.tools.SetCurrentTool(tool)
		}
		window.workspace.canvas.Repaint()
	}

	for _, btn := range window.btnTools {
		btn.SetClickEvent(fbtntools)
	}

	boutline := shapes.AddImageButton("Outline", ".\\icons\\outline.png", RibbonButtonSizeMedium)

	boutlineMenu := NewPopupMenu(ribbon, []MenuItemInfo{
		{Text: "No outline", IconPath: ".\\icons\\no-fill-small.png", AssignTo: &window.menuNoOutline},
		{Text: "Solid color", IconPath: ".\\icons\\solid-color-small.png", AssignTo: &window.menuSolidOutline},
	})
	window.menuSolidOutline.SetToggled(true)
	boutline.SetDropdownMenu(boutlineMenu, false)

	window.menuNoOutline.SetClickEvent(func(e *PopupItemEvent) {
		window.menuNoOutline.SetToggled(true)
		window.menuSolidOutline.SetToggled(false)
	})
	window.menuSolidOutline.SetClickEvent(func(e *PopupItemEvent) {
		window.menuSolidOutline.SetToggled(true)
		window.menuNoOutline.SetToggled(false)
	})

	bfill := shapes.AddImageButton("Fill", ".\\icons\\fill-type.png", RibbonButtonSizeMedium)

	bfillMenu := NewPopupMenu(ribbon, []MenuItemInfo{
		{Text: "No fill", IconPath: ".\\icons\\no-fill-small.png", AssignTo: &window.menuNoFill},
		{Text: "Solid color", IconPath: ".\\icons\\solid-color-small.png", AssignTo: &window.menuSolidFill},
	})
	window.menuNoFill.SetToggled(true)
	bfill.SetDropdownMenu(bfillMenu, false)

	window.menuNoFill.SetClickEvent(func(e *PopupItemEvent) {
		window.menuNoFill.SetToggled(true)
		window.menuSolidFill.SetToggled(false)
	})
	window.menuSolidFill.SetClickEvent(func(e *PopupItemEvent) {
		window.menuSolidFill.SetToggled(true)
		window.menuNoFill.SetToggled(false)
	})

	sizesec := home.AddSection("")

	funcSize := func(e *PopupItemEvent) {
		for _, item := range e.PopupWindow.GetItems() {
			(item.(PopupSizeMenuItem)).SetToggled(false)
		}
		sitem := e.Item.(PopupSizeMenuItem)
		sitem.SetToggled(true)
		tool := window.tools.GetCurrentTool()
		tool.changeSize(sitem.GetSize())
		// we should now repaint the popupwindow....nvm
	}

	window.bsize = sizesec.AddImageButton("Size", ".\\icons\\size.png", RibbonButtonSizeBig)
	window.bsizeMenu = NewPopupSizeMenu(ribbon, []SizeMenuItemInfo{
		{Size: 1, Toggled: true, OnClick: funcSize},
		{Size: 2, OnClick: funcSize},
		{Size: 3, OnClick: funcSize},
		{Size: 4, OnClick: funcSize},
	})
	window.bsize.SetDropdownMenu(window.bsizeMenu, false)
	window.bsize.SetEnabled(false)

	// init all the color buttons
	window.initRibbonColorSection(home)

	// view
	view := ribbon.AddTab("View")
	szoom := view.AddSection("Zoom")

	szoom.AddImageButton("Zoom\nin", ".\\icons\\zoom-in.png", RibbonButtonSizeBig).SetEnabled(false)
	szoom.AddImageButton("Zoom\nout", ".\\icons\\zoom-out.png", RibbonButtonSizeBig).SetEnabled(false)
	szoom.AddImageButton("100\n%", ".\\icons\\zoom-100.png", RibbonButtonSizeBig)

	shideshow := view.AddSection("Show or hide")

	shideshow.AddCheckButton("Rulers", false).SetEnabled(false)
	window.bShowGridlines = shideshow.AddCheckButton("Gridlines", false)
	window.bShowGridlines.SetClickEvent(func(e *RibbonButtonEvent) {
		window.workspace.canvas.RepaintVisible()
	})
	bstatusbar := shideshow.AddCheckButton("Status bar", true)
	bstatusbar.SetClickEvent(func(e *RibbonButtonEvent) {
		if bstatusbar.IsToggled() {
			window.statusbar.SetVisible(true)
			window.statusbar.SetDockType(DockBottom)
		} else {
			window.statusbar.SetVisible(false)
			window.statusbar.SetDockType(DockNone)
		}
	})

	sdisplay := view.AddSection("Display")
	sdisplay.AddImageButton("Full\nscreen", ".\\icons\\full-screen.png", RibbonButtonSizeBig).SetEnabled(false)
	sdisplay.AddImageButton("Thumbnail", ".\\icons\\thumbnail.png", RibbonButtonSizeBig).SetEnabled(false)

	ribbon.SetCurrentTab(home)
	ribbon.ResumeRepaint()
}

func (window *MainWindow) SetCurrentTool(newTool Tool) {
	if window.tools.GetCurrentTool() == newTool {
		return
	}
	for button, tool := range window.buttonToolPairs {
		if tool == newTool {
			for _, btn := range window.btnTools {
				btn.SetToggled(false)
			}
			button.SetToggled(true)
			// it will get enabled by the tool if required
			window.bsize.SetEnabled(false)
			window.tools.SetCurrentTool(tool)
			break
		}
	}
	window.workspace.canvas.Repaint()
}

func (window *MainWindow) initRibbonColorSection(home *RibbonTab) {
	logInfo("init color section...")
	colors := home.AddSection("Colors")

	window.color1 = colors.AddButton("Color\n1", Rgb(160, 160, 160), Rgb(0, 0, 0), RibbonButtonSizeBig)
	window.color2 = colors.AddButton("Color\n2", Rgb(160, 160, 160), Rgb(255, 255, 255), RibbonButtonSizeBig)

	window.color1.SetToggled(true)

	window.color1.SetClickEvent(func(e *RibbonButtonEvent) {
		window.color2.SetToggled(false)
		window.color1.SetToggled(true)
	})

	window.color2.SetClickEvent(func(e *RibbonButtonEvent) {
		window.color1.SetToggled(false)
		window.color2.SetToggled(true)
	})

	fForegroundBackground := func() RibbonButton {
		if window.color1.IsToggled() {
			return window.color1
		}
		return window.color2
	}

	fcolors := func(e *RibbonButtonEvent) {
		button := e.Button
		color := button.GetColor()
		fForegroundBackground().SetColor(color)
		window.workspace.canvas.RepaintVisible()
	}

	colorbuttons := colors.AddColorButtons([]RibbonColorButton{
		{Name: "Black", Color: Rgb(0, 0, 0), OnClick: fcolors},
		{Name: "White", Color: Rgb(255, 255, 255), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Gray", Color: Rgb(127, 127, 127), OnClick: fcolors},
		{Name: "Light Gray", Color: Rgb(195, 195, 195), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Dark Red", Color: Rgb(136, 0, 21), OnClick: fcolors},
		{Name: "Brown", Color: Rgb(185, 122, 87), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Red", Color: Rgb(255, 0, 0), OnClick: fcolors},
		{Name: "Rose", Color: Rgb(255, 174, 201), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Orange", Color: Rgb(255, 127, 39), OnClick: fcolors},
		{Name: "Gold", Color: Rgb(255, 201, 14), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Yellow", Color: Rgb(255, 242, 0), OnClick: fcolors},
		{Name: "Light Yellow", Color: Rgb(239, 228, 176), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Green", Color: Rgb(34, 177, 76), OnClick: fcolors},
		{Name: "Lime", Color: Rgb(181, 230, 29), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Turquoise", Color: Rgb(0, 162, 232), OnClick: fcolors},
		{Name: "Light Turquoise", Color: Rgb(153, 217, 234), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Indigo", Color: Rgb(63, 72, 204), OnClick: fcolors},
		{Name: "Blue-Gray", Color: Rgb(112, 146, 190), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},

		{Name: "Purple", Color: Rgb(163, 73, 164), OnClick: fcolors},
		{Name: "Lavender", Color: Rgb(200, 191, 231), OnClick: fcolors},
		{Name: "Custom", Color: Rgb(245, 245, 245), OnClick: fcolors},
	})

	beditcolors := colors.AddImageButton("Edit\ncolors", ".\\icons\\edit-colors.png", RibbonButtonSizeBig)

	customColorButtons := make([]RibbonButton, 0)
	for _, button := range colorbuttons {
		if button.GetText() == "Custom" {
			button.SetEnabled(false)
			customColorButtons = append(customColorButtons, button)
		}
	}

	beditcolors.SetClickEvent(func(e *RibbonButtonEvent) {
		var prevCustomColors [16]Color
		// Obtain the 10 colors from our custom color buttons
		for i, button := range customColorButtons {
			prevCustomColors[i] = button.GetColor()
		}
		// Fill the rest 6 with white colors
		for i := 10; i < 16; i++ {
			prevCustomColors[i] = Rgb(255, 255, 255)
		}
		prevColor := fForegroundBackground().GetColor()
		newColor, _ := ChoseColorDialog(window, prevColor, prevCustomColors)
		fForegroundBackground().SetColor(newColor)
		// Check if it's a new color, then we add it into the custom colors
		if !newColor.IsEqualTo(prevColor) {
			alreadyExist := false
			var emptyButton RibbonButton = nil
			for _, button := range customColorButtons {
				color := button.GetColor()
				if newColor.IsEqualTo(color) {
					alreadyExist = true
					break
				}
				// Also grab a empty custom color button
				if !button.IsEnabled() {
					// If we haven't grab one already
					if emptyButton == nil {
						emptyButton = button
					}
				}
			}
			// We have a new color that doesn't exit in the custom color list
			if !alreadyExist {
				// Check if we don't have an empty button to put the new color in
				if emptyButton == nil {
					// Then we move all the colors to left and we put the new color in the last button
					colorsSliced := prevCustomColors[1:10] // Slice it and grab the last 9
					for i, color := range colorsSliced {
						customColorButtons[i].SetColor(color)
					}
					// last button
					customColorButtons[9].SetColor(newColor)
				} else {
					// We good
					emptyButton.SetColor(newColor)
					emptyButton.SetEnabled(true)
				}
			}
		}
	})
}
