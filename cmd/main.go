package main

import (
	"flag"
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

const commonLogPath = "./log/common.log"

func main() {
	info := &entities.Info{}

	//log.Printf("Количество горутин в начале запуска: %d", runtime.NumGoroutine())

	logVerbose := flag.Bool("log", false, "используется для логирования в консоль всех уровней")
	flag.Parse()
	mainLogger := services.NewLogger(*logVerbose)

	if err := mainLogger.InitLog(); err != nil {
		log.Fatal(err.Error())
	}

	fptrDriver, err := fptr10.NewSafe()
	defer fptrDriver.Destroy()
	mainLogger.Infoln("Запуск драйвера KKT")
	if err != nil {
		mainLogger.Errorf("Запуск драйвера ККТ завершился с ошибкой: %v", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	gateway := gateways.NewGateway(client, fptrDriver)
	service := services.NewServices(gateway, mainLogger)
	view := view.NewFyneApp(app.New(), service, info)

	service.LoggerService.Infoln("Начало работы приложения")
	//service.LoggerService.ReinitDebugger(time.Hour * 24)

	defer service.Logger.Close()
	defer service.Logger.Infoln("Завершение работы приложения")

	view.StartApp()

}
