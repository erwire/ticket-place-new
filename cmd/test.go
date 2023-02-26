package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	LongText := "dafajkfas kas fas fas afs afsjk fajks faf asf jashf asjkh fafjhajsf jks fa fas fjas fa fa fasa fsjk faf aka fa fas fasf naknas f"

	ap := app.New()
	alert := ap.NewWindow("Ошибка")
	text := canvas.NewText(LongText, theme.ForegroundColor())
	text.TextSize = 20
	alert.Resize(fyne.NewSize(400, 200))
	alert.SetPadded(true)
	box := container.NewBorder(container.NewCenter(widget.NewLabel("Во время исполнения произошла ошибка")), container.NewCenter(widget.NewButtonWithIcon("Подвердить", theme.ConfirmIcon(), func() {
		alert.Hide()
	})), text, nil)
	alert.SetContent(container.NewMax(canvas.NewRectangle(theme.BackgroundColor()), box))

	w := ap.NewWindow("Hello")
	w.Resize(fyne.NewSize(200, 200))
	selectEntryOne := widget.NewSelectEntry([]string{"odin", "dva", "tri"})
	but := widget.NewButtonWithIcon("Press Me@", theme.AccountIcon(), func() {
		alert.Hide()
		alert.Show()
	})

	w.SetContent(container.NewVBox(selectEntryOne, but))
	w.ShowAndRun()

}
