package update

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/logger"
	"os"
)

const logPath = "./log/"
const icoPath = "./content/system/icon/main.png"

func Build() {

	var verbose = flag.Bool("log", false, "печать информ-логи в консоль")
	flag.Parse()

	file, err := os.OpenFile(logPath+"upd.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Errorf("Ошибка при открытии файла логов")
	}

	defer logger.Init("Logger", *verbose, true, file).Close()
	defer file.Close()
	logger.Info("Запущен центр обновлений")

	version := os.Getenv("version")
	execName := os.Getenv("exec_name")
	repo := os.Getenv("repo")
	owner := os.Getenv("owner")
	pid := os.Getenv("pid")

	fmt.Println(version)

	app := app.New()
	mainWindow := app.NewWindow("Центр обновления")
	resource, err := fyne.LoadResourceFromPath(icoPath)
	if err != nil {
		logger.Infof("Во время привязки иконки произошла ошибка: %v", err)
	}
	mw := NewMainWindow(mainWindow)
	mw.Window.SetIcon(resource)
	updatePath, err := os.Executable()
	if err != nil {
		mw.EnvError("Невозможно определить имя исполняемого файла центра обновления")
	}

	if execName == "" {
		mw.EnvError("В центр обновления не передано имя исполняемого файла")
		return
	}
	if version == "" {
		mw.EnvError("В центр обновления не переданы данные о версии ПО")
		return
	}

	if repo == "" || owner == "" {
		mw.EnvError("В центр обновления не переданы данные о источнике обновления")
		return
	}

	if pid == "" {
		mw.EnvError("Не переданы данные о номере процесса")
		return
	}

	mw.Config.updatePath = updatePath
	mw.Config.version = version
	mw.Config.execPath = execName
	mw.Config.repo = repo
	mw.Config.owner = owner
	mw.Config.PID = pid

	labelsNames := []string{labelQuestion, labelVersion}
	labelVers := canvas.NewText("Текущая версия: "+version, theme.ForegroundColor())
	labelVers.Alignment = fyne.TextAlignCenter
	labelQuest := canvas.NewText("Проверить обновление?", theme.ForegroundColor())
	labelQuest.Alignment = fyne.TextAlignCenter
	logoImage := canvas.NewImageFromFile("./content/system/icon/logo.png")
	logoImage.FillMode = canvas.ImageFillContain
	imageBox := container.NewCenter(container.NewGridWrap(fyne.NewSize(200, 200), logoImage))

	buttonsNames := []string{buttonConfirm, buttonCancel}
	buttonConfirm := widget.NewButtonWithIcon("Да", theme.ConfirmIcon(), mw.CheckUpdate)
	buttonCancel := widget.NewButtonWithIcon("Нет", theme.CancelIcon(), mw.Close)

	box := container.NewVBox(imageBox, labelVers, labelQuest, buttonConfirm, buttonCancel)

	//инициализация компонентов главного окна
	mw.InitButtons(buttonsNames, buttonConfirm, buttonCancel)
	mw.InitLabels(labelsNames, labelQuest, labelVers)
	mw.InitBox(box)
	mw.ConfigurateWindow()

	windows := NewWindows(mw)
	updateApp := NewUpdater(app, windows)
	updateApp.StartApp()
}
