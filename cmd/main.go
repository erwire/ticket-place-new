package main

import (
	"fptr/internal/controller/view"
	"fptr/internal/entities"
	"fptr/internal/gateways"
	"fptr/internal/services"
	"fptr/pkg/fptr10"
	"fyne.io/fyne/v2/app"
	"log"
	"net/http"
	"time"
)

func main() {
	info := &entities.Info{}
	//log.Printf("Количество горутин в начале запуска: %d", runtime.NumGoroutine())

	fptrDriver, err := fptr10.NewSafe()
	defer fptrDriver.Destroy()

	if err != nil {
		log.Println(err.Error())
	}

	client := &http.Client{Timeout: 2 * time.Second}
	gateway := gateways.NewGateway(client, fptrDriver)
	service := services.NewServices(gateway)
	view := view.NewFyneApp(app.New(), service, info)

	view.StartApp()

}
