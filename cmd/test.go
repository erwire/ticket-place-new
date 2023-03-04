package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	ap := app.New()
	w := ap.NewWindow("Test")
	label1 := widget.NewLabel("Текст 1")
	label2 := widget.NewLabel("Текст 2")
	button1 := widget.NewButton("Кнопка 1", nil)
	button2 := widget.NewButton("Кнопка 2", nil)
	w.SetContent(container.NewHBox(label1, label2, layout.NewSpacer(), button1, button2))
	w.ShowAndRun()
}
