package main

import (
	"fptr/internal/controller/view"
	"fyne.io/fyne/v2/app"
)

func main() {

	v := view.NewFyneApp(app.New())
	v.StartApp()
}
