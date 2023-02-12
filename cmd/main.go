package main

import (
	"fptr/internal/controller/view"
	"fptr/internal/entities"
	"fptr/internal/gateways"
	"fptr/internal/services"
	"fyne.io/fyne/v2/app"
	fptr10 "github.com/EuginKostomarov/ftpr10"
	"log"
	"net/http"
	"time"
)

func main() {
	info := &entities.Info{}

	fptrDriver, err := fptr10.NewSafe()

	if err != nil {
		log.Println(err.Error())
	}

	client := &http.Client{Timeout: 2 * time.Second}
	gateway := gateways.NewGateway(client, fptrDriver)
	service := services.NewServices(gateway)
	view := view.NewFyneApp(app.New(), service, info)
	view.StartApp()

}
