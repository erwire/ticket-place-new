package view

import (
	"context"
	"fptr/internal/entities"
	"fptr/internal/services"
	"fptr/pkg/toml"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type Flags struct {
	PrintOnKKTTicketCheckBox, PrintCheckBox, PrintOnPrinterTicketBox bool
	StopListen                                                       bool
	DebugOn                                                          bool
}

type Selected struct {
	Additional *string
}

type FyneApp struct {
	info *entities.Info

	context struct {
		ctx    context.Context
		cancel context.CancelFunc
	}
	//приложение
	service     *services.Services
	application fyne.App
	//! главное окно
	mainWindow fyne.Window
	//элементы окна
	authForm struct {
		form                      dialog.Dialog
		loginEntry, passwordEntry *widget.Entry
		settingButton             *widget.Button
	}

	header struct {
		usernameLabel           *canvas.Text
		localTimeLabel          *canvas.Text
		printLastСheckButton    *widget.Button
		exitButton              *widget.Button
		exitAndCloseShiftButton *widget.Button
		printXReportButton      *widget.Button

		listenerStatus struct {
			listenerToolbar     *widget.Toolbar
			listenerToolbarItem *widget.ToolbarAction
		}
	}
	PrintSettingsItem struct {
		//PrintSettingsAccordionItem *widget.AccordionItem
		PrintSettingsContainer *fyne.Container
		PrintCheck             *widget.Check
		PrintOnKKT             *widget.Check
		PrintOnPrinter         *widget.Check
		AdditionalText         *widget.Entry
		SetAdditionalText      *widget.Button
	}

	PrintsRefoundAndDeposits struct {
		RefoundAndDepositsAccordionItem *widget.AccordionItem
		RefoundForm                     *widget.Form
		RefoundFormItem                 *widget.FormItem
		RefoundEntry                    *widget.Entry

		CashIncomeForm     *widget.Form
		CashIncomeFormItem *widget.FormItem
		CashIncomeEntry    *widget.Entry

		PrintCheckForm     *widget.Form
		PrintCheckFormItem *widget.FormItem
		PrintCheckEntry    *widget.Entry
	}

	DriverSetting struct {
		DriverSettingAccordion                                  *widget.AccordionItem
		DriverSettingButton, DriverPrintHistoryButton           *widget.Button
		DriverSettingLabel                                      *widget.Label
		DriverComPortEntry, DriverPathEntry, DriverAddressEntry *widget.Entry
		DriverPollingPeriodSelect                               *widget.Select
		DriverSettingForm                                       *widget.Form
		DriverKKTComFormItem                                    *widget.FormItem
		DriverKKTPathFormItem                                   *widget.FormItem
		DriverApiAddressFormItem                                *widget.FormItem
		DriverPollingPeriodFormItem                             *widget.FormItem
	}

	PopUp struct {
		PopUp *widget.PopUp
		Text  *widget.Label
	}

	MainWindowAccordion *widget.Accordion

	//! Окно ошибок
	Warning struct {
		WarningWindow dialog.Dialog
		WarningText   *canvas.Text
	}

	AlertWindow   dialog.Dialog
	SettingWindow dialog.Dialog
	//Флаги
	flag     Flags
	selected Selected
}

func NewFyneApp(a fyne.App, view *services.Services, inf *entities.Info) *FyneApp { //, service *services.Service
	return &FyneApp{
		application: a,
		service:     view,
		info:        inf,
	}
}

func (f *FyneApp) StartApp() {

	f.ConfigureMainWindows()
	f.ConfigureAuthDialogForm()
	f.ConfigureWarningAlert()
	f.ConfigureSettingWindow()
	err := f.InitializeCookie()
	if err != nil {
		f.ShowWarning("Данные по прошлой сессии повреждены или отсутствуют")
		f.UpdateSession(entities.SessionInfo{})
	}

	if f.info.Session.IsDead() {
		f.ShowWarning("Ваша сессия устарела. Пожалуйста, авторизуйтесь снова!")
		f.UpdateSession(entities.SessionInfo{})
	} else {
		f.header.usernameLabel.Text = f.info.Session.UserData.Username
		f.authForm.form.Hide()
		f.context.ctx, f.context.cancel = context.WithCancel(context.Background())
		click, message := f.service.GetLastReceipt(f.info.AppConfig.Driver.Connection, f.info.Session)
		if message != "" {
			f.ShowWarning("Внимание, возможно задвоение чека!")
			log.Println(err.Error()) //обработку сделать нормальную
		}
		err = toml.WriteToml(toml.ClickPath, click)
		if err != nil {
			log.Println(err.Error())
		}
		message = f.service.MakeSession(*f.info)
		if message != "" {
			f.UpdateSession(entities.SessionInfo{})
			f.authForm.form.Show()
		}
		f.service.Infof("Данные из сессии подгрузились. Успешная авторизация под пользователем %s", f.info.Session.UserData.FullName)
		go f.Listen(f.context.ctx, *f.info)
	}

	f.setupCookieIntoEntry()
	go f.ClockUpdater()
	f.flag.DebugOn = false

	//f.PopUp.Text = widget.NewLabel("")
	//f.PopUp.PopUp = widget.NewPopUp(f.PopUp.Text, f.mainWindow.Canvas())

	f.mainWindow.ShowAndRun()
}

func (f *FyneApp) ClockUpdater() {
	for {
		time.Sleep(time.Second)
		timeLog, err := f.service.LoggerService.CurrentTime()
		if err != nil {
			f.service.Warning(err)
		}
		if time.Now().Format("02-01-2006") != timeLog.Format("02-01-2006") {
			if err := f.service.LoggerService.Reinit(); err != nil {
				f.service.Warning(err)
			}
		}
		f.header.localTimeLabel.Text = time.Now().Format("02.01.2006 15:04:05")
		f.header.localTimeLabel.Refresh()
	}
}

// + Главное окно
// ! Окно
