package update

import (
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/blang/semver"
	"github.com/google/logger"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func (w *MainWindow) CheckUpdate() {
	logger.Infof("Запущен процесс поиска обновлений")

	logger.Infof("Исполняемый файл: " + w.Config.execPath + "\n")
	logger.Infof("Файл обновления: " + w.Config.updatePath + "\n")
	logger.Infof("Версия: " + w.Config.version + "\n")
	logger.Infof("Источник: " + w.Config.owner + "/" + w.Config.repo)

	latest, found, err := selfupdate.DetectLatest(w.Config.owner + "/" + w.Config.repo)
	if err != nil {
		logger.Errorf("Ошибка в получении данных по обновлению")
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("Ошибка при получении данных по обновлению")), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)
		})
		d.Show()
		return
		//ошибка подключения к источнику обновления
	}

	v := semver.MustParse(w.Config.version)

	if !found {
		logger.Infof("Обновления не найдены в данном репозитории")
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("В указанном источнике не найдены обновления")), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)
		})
		d.Show()

		return
	}
	logger.Infof("Последняя версия в источнике: " + latest.Version.String() + "\n")
	if latest.Version.LTE(v) {
		logger.Infof("Версия текущего релиза является последней")
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("У вас последняя версия ПО")), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)
		})
		d.Show()
		logger.Info("Центр обновлений завершил работу")
		return

	}

	if err := selfupdate.UpdateTo(latest.AssetURL, w.Config.execPath); err != nil {
		logger.Errorf("Ошибка при обновлении ПО: %v", err)
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("Ошибка при обновлении ПО")), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)
		})
		d.Show()
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, w.Config.updatePath); err != nil {
		logger.Errorf("Ошибка при обновлении центра обновлений: %v", err)
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("Ошибка при обновлении центра обновления")), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)
		})

		d.Show()

		return
	}

	pid, err := strconv.Atoi(w.Config.PID)
	if err != nil {
		logger.Errorf("Ошибка при конвертации в число строки, содержащей PID: %v", err)
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("Внутренняя ошибка центра обновления")), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)
		})

		d.Show()

		return
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		logger.Errorf("Ошибка при обнаружении процесса: %v", err)
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("Ошибка в обнаружении процесса приложения")), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)
		})

		d.Show()
		return
	}
	proc.Kill()
	time.Sleep(2 * time.Second)
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", w.Config.execPath)

	if err := cmd.Start(); err != nil {
		logger.Errorf("Ошибка при перезапуске процесса: %v", err)
		d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel("Ошибка в перезапуске процесса приложения"+err.Error())), w.Window)
		d.SetOnClosed(func() {
			logger.Info("Центр обновлений завершил работу")
			os.Exit(1)

		})

		d.Show()
		return
	}
	cmd.Process.Release()
	logger.Info("Центр обновлений успешно обновил ПО до версии ", latest.Version.String())
	d := dialog.NewCustom("Обновление", "Закрыть", container.NewCenter(widget.NewLabel(fmt.Sprintf("Успешное обновление ПО до версии %s", latest.Version))), w.Window)
	d.SetOnClosed(func() {
		logger.Info("Центр обновлений завершил работу")
		os.Exit(0)

	})

	d.Show()

}

func (w *MainWindow) Close() {
	logger.Info("Центр обновлений завершил работу")
	os.Exit(0)
}
