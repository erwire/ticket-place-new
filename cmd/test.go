package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

var data = []string{"a", "string", "list"}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")
	var data = [][]string{[]string{"top left", "top right", "ыфв", "ыфв", "выфв"},
		[]string{"top left", "top right", "ыфв", "ыфв", "выфв"},
		[]string{"top left", "top right", "ыфв", "ыфв", "выфв"},
		[]string{"top left", "top right", "ыфв", "ыфв", "выфв"},
		[]string{"top left", "top right", "ыфв", "ыфв", "выфв"}}

	table := widget.NewTableWithHeaders(
		func() (rows int, cols int) {
			return 5, 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Hi")
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(data[id.Row][id.Col])
		},
	)
	table.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		if id.Row == -1 {
			switch id.Col {
			case 0:
				template.(*widget.Label).SetText("Событие")
			case 1:
				template.(*widget.Label).SetText("Сумма")
			case 2:
				template.(*widget.Label).SetText("Ошибка")
			case 3:
				template.(*widget.Label).SetText("Статус")
			case 4:
				template.(*widget.Label).SetText("Сумма")
			}

		}

		if id.Row >= 0 {
			template.(*widget.Label).SetText(strconv.Itoa(id.Row))
		}
	}

	for i := 0; i < 5; i++ {
		table.SetColumnWidth(i, 100)
	}

	y := container.NewGridWithRows(2,
		widget.NewButton("Привет", func() {

		}), widget.NewButton("Пока", func() {

		}))
	table.Resize(fyne.NewSize(500, 500))
	myWindow.SetContent(container.NewBorder(nil, y, nil, nil, container.NewVScroll(table)))
	myWindow.ShowAndRun()
}
