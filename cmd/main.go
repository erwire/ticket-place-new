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
	"os"
	"time"
)

var (
	version    = "1.0.11"
	updatePath = "jahngeor"
	updateRepo = "test"
	updType    = "github"
)

func main() {
	info := &entities.Info{}

	logVerbose := flag.Bool("log", false, "используется для логирования в консоль всех уровней")
	flag.Parse()
	mainLogger := services.NewLogger(*logVerbose)

	if err := mainLogger.InitLog(); err != nil {
		log.Fatal(err.Error())
	}

	middleware.NewMiddleware(mainLogger.Logger).BasicMiddleware()

	fptrDriver, err := fptr10.NewSafe()
	if err != nil {
		mainLogger.Fatal(err.Error())
		os.Exit(1)
	}

	fptrDriver.SetSingleSetting(fptr10.LIBFPTR_SETTING_AUTO_RECONNECT, "false")
	if err := fptrDriver.ApplySingleSettings(); err != nil {
		mainLogger.Fatalf("Ошибка активации настроек ККТ: %s", err.Error())
	}

	defer fptrDriver.Destroy()
	mainLogger.Infoln("Запуск драйвера KKT")

	client := &http.Client{Timeout: 20 * time.Second}
	gateway := gateways.NewGateway(client, fptrDriver)
	service := services.NewServices(gateway, mainLogger)
	viewController := view.NewFyneApp(app.New(), service, info)

	viewController.SetAppInfo(version, updatePath, updType, updateRepo)

	service.LoggerService.Infoln("Начало работы приложения")
	//service.LoggerService.ReinitDebugger(time.Hour * 24)

	defer service.LoggerService.Close()
	defer service.Logger.Infoln("Завершение работы приложения")

	viewController.StartApp()

	//if err != nil {
	//	service.Errorf("Запуск драйвера ККТ завершился с ошибкой: %v", err)
	//	view.ShowCriticalError(err, "Пожалуйста, скачайте драйвер и перезапустите приложение", "https://atoldriver.ru/")
	//	return
	//}

}
