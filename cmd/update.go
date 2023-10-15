package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var data = []string{"a", "string", "list"}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(4, widget.NewCheck("", func(b bool) {

			}), widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[1].(*widget.Label).SetText("Привет")
			o.(*fyne.Container).Objects[2].(*widget.Label).SetText("Мир")
		})

	myWindow.SetContent(list)
	myWindow.ShowAndRun()
}
