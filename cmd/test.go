package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Toolbar Widget")

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),

		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
	image := canvas.NewImageFromFile("./content/system/icon/logo.png")
	image.FillMode = canvas.ImageFillStretch
	boxImage := container.NewGridWrap(fyne.NewSize(35, 35), image)
	text := canvas.NewText("Ticket Place", theme.ForegroundColor())
	text.TextStyle = fyne.TextStyle{Bold: true}
	text.TextSize = 18
	toolbarAct := container.New(layout.NewFormLayout(), container.NewHBox(boxImage, text), toolbar)

	content := container.NewBorder(toolbarAct, nil, nil, nil, widget.NewLabel("Content"))
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

//func main() {
//	myApp := app.New()
//	myWindow := myApp.NewWindow("Toolbar Widget")
//
//	toolbar := widget.NewToolbar(
//		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
//			log.Println("New document")
//		}),
//		widget.NewToolbarSeparator(),
//		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
//		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
//		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
//		widget.NewToolbarSpacer(),
//		widget.NewToolbarAction(theme.HelpIcon(), func() {
//			log.Println("Display help")
//		}),
//	)
//	toolbarString := container.New
//	content := container.NewBorder(toolbar, nil, nil, nil, widget.NewLabel("Content"))
//	myWindow.SetContent(content)
//	myWindow.ShowAndRun()
//}
