package main

import (
	"flag"
	"fptr/cmd/middleware"
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

var version = "1.1.0"
var updatePath = "jahngeor"
var updateRepo = "test"
var updType = "github"

func main() {
	info := &entities.Info{}

	//log.Printf("Количество горутин в начале запуска: %d", runtime.NumGoroutine())

	logVerbose := flag.Bool("log", false, "используется для логирования в консоль всех уровней")
	flag.Parse()
	mainLogger := services.NewLogger(*logVerbose)

	if err := mainLogger.InitLog(); err != nil {
		log.Fatal(err.Error())
	}

	middleware.NewMiddleware(mainLogger.Logger).BasicMiddleware()

	fptrDriver, err := fptr10.NewSafe()

	mainLogger.Infoln("Запуск драйвера KKT")

	client := &http.Client{Timeout: 20 * time.Second}
	gateway := gateways.NewGateway(client, fptrDriver)
	service := services.NewServices(gateway, mainLogger)
	fyneApp := view.NewFyneApp(app.New(), service, info)

	fyneApp.SetAppInfo(version, updatePath, updType, updateRepo)

	if err != nil {
		service.Errorf("Запуск драйвера ККТ завершился с ошибкой: %v", err)
		fyneApp.ShowCriticalError(err, "Пожалуйста, скачайте драйвер и перезапустите приложение", "https://atoldriver.ru/")
		return
	}

	service.LoggerService.Infoln("Начало работы приложения")
	//service.LoggerService.ReinitDebugger(time.Hour * 24)

	defer service.LoggerService.Close()
	defer service.Logger.Infoln("Завершение работы приложения")

	fyneApp.StartApp()

}
