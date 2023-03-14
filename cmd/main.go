package main

import (
	"flag"
	"fmt"
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

var version = "1.0.7"
var updatePath = "jahngeor/ticket-place-new"
var updType = "github"

func main() {
	info := &entities.Info{}
	createAppDirectories()
	//log.Printf("Количество горутин в начале запуска: %d", runtime.NumGoroutine())

	logVerbose := flag.Bool("log", false, "используется для логирования в консоль всех уровней")
	flag.Parse()
	mainLogger := services.NewLogger(*logVerbose)

	if err := mainLogger.InitLog(); err != nil {
		log.Fatal(err.Error())
	}

	fptrDriver, err := fptr10.NewSafe()
	fptrDriver.SetSingleSetting(fptr10.LIBFPTR_SETTING_AUTO_RECONNECT, "false")
	fptrDriver.ApplySingleSettings()
	fmt.Println(fptrDriver.GetSingleSetting(fptr10.LIBFPTR_SETTING_AUTO_RECONNECT))
	defer fptrDriver.Destroy()
	mainLogger.Infoln("Запуск драйвера KKT")

	client := &http.Client{Timeout: 20 * time.Second}
	gateway := gateways.NewGateway(client, fptrDriver)
	service := services.NewServices(gateway, mainLogger)
	view := view.NewFyneApp(app.New(), service, info)

	view.SetAppInfo(version, updatePath, updType)

	service.LoggerService.Infoln("Начало работы приложения")
	//service.LoggerService.ReinitDebugger(time.Hour * 24)

	defer service.Logger.Close()
	defer service.Logger.Infoln("Завершение работы приложения")

	view.StartApp()

	if err != nil {
		service.Errorf("Запуск драйвера ККТ завершился с ошибкой: %v", err)
		view.ShowCriticalError(err, "Пожалуйста, скачайте драйвер и перезапустите приложение", "https://atoldriver.ru/")
		return
	}

}

func createAppDirectories() {
	paths := directoriesList()
	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			log.Printf("Создаем папку %s\n", path)
			_ = os.Mkdir(path, 0660)
		}
	}
}

func directoriesList() []string {
	return []string{
		"./log",
		"./debug_info",
		"./debug_info/sell",
		"./debug_info/refound",
		"./debug_info/click",
		"./debug_info/login",
		"./content",
		"./content/system",
		"./content/system/icon",
		"./cookie",
		"./cookie/appconfig",
		"./cookie/click",
		"./cookie/session",
		"./cookie/userdata",
		"./db",
	}
}
