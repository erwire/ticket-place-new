package view

import (
	"context"
	"fptr/internal/entities"
	"fptr/internal/services"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"time"
)

type Flags struct {
	PrintOnKKTTicketCheckBox, PrintCheckBox, PrintOnPrinterTicketBox bool
	StopListen                                                       bool
	DebugOn                                                          bool
	SoundError                                                       bool
	ProgressWorking                                                  bool
	AuthJustHide                                                     bool
	FirstStart                                                       bool
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
	MainWindow fyne.Window
	//элементы окна
	authForm struct {
		form                      dialog.Dialog
		loginEntry, passwordEntry *widget.Entry
		settingButton             *widget.Button
	}

	header struct {
		usernameLabel  *canvas.Text
		localTimeLabel *canvas.Text

		listenerStatus struct {
			listenerToolbar     *widget.Toolbar
			listenerToolbarItem *widget.ToolbarAction
		}
	}
	PrintSettingsItem struct {
		exitButton              *widget.Button
		exitAndCloseShiftButton *widget.Button
		//PrintSettingsAccordionItem *widget.AccordionItem
		PrintSettingsContainer *fyne.Container
		PrintCheck             *widget.Check
		PrintOnKKT             *widget.Check
		PrintOnPrinter         *widget.Check
		AdditionalText         *widget.Entry
		SetAdditionalText      *widget.Button
		CashIncomeForm         *widget.Form
		CashIncomeFormItem     *widget.FormItem
		CashIncomeEntry        *widget.Entry
		printLastСheckButton   *widget.Button
		printXReportButton     *widget.Button
	}

	PrintsRefoundAndDeposits struct {
		RefoundAndDepositsAccordionItem *widget.AccordionItem
		RefoundForm                     *widget.Form
		RefoundFormItem                 *widget.FormItem
		RefoundEntry                    *widget.Entry
		AdminEntry                      *widget.Entry
		AdminFormItem                   *widget.FormItem
		AdminForm                       *widget.Form
		PrintCheckForm                  *widget.Form
		PrintCheckFormItem              *widget.FormItem
		PrintCheckEntry                 *widget.Entry
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
		CloseShiftButton                                        *widget.Button
		ErrorSoundButton                                        *widget.Button
		PrintLastButton                                         *widget.Button
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
	InfoDialogWindow dialog.Dialog
	InfoDialogText   canvas.Text
	AlertWindow      dialog.Dialog
	SettingWindow    dialog.Dialog
	//Флаги
	flag     Flags
	selected Selected

	Error struct {
		ErrorWindow        fyne.Window
		ErrorConfirmButton *widget.Button
		ErrorText          *canvas.Text
	}

	CriticalError struct {
		ErrorWindow        fyne.Window
		ErrorText          *canvas.Text
		AdditionalText     *canvas.Text
		ErrorConfirmButton *widget.Button
		ErrorLinkButton    *widget.Hyperlink
	}

	Reconnector struct {
		Progresser      fyne.Window
		ProgresserText  *canvas.Text
		ProgressBar     *widget.ProgressBarInfinite
		ProgressConfirm *widget.Button
		ProgressDismiss *widget.Button
	}
}

func NewFyneApp(a fyne.App, view *services.Services, inf *entities.Info) *FyneApp { //, service *services.Service
	return &FyneApp{
		application: a,
		service:     view,
		info:        inf,
	}
}

func (f *FyneApp) StartApp() {
	if err := f.service.Open(); err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
	}
	f.ConfigureMainWindows()
	f.ConfigureAuthDialogForm()
	f.ConfigureWarningAlert()
	f.ConfigurateErrorAlert()
	f.ConfigurateCriticalErrorAlert()
	f.ConfigureSettingWindow()
	f.ConfigureProgresser()
	err := f.InitializeCookie()
	if err != nil {
		f.ShowWarning("Данные по прошлой сессии повреждены или отсутствуют")
		f.UpdateSession(entities.SessionInfo{})
	}

	if f.info.Session.IsDead() {
		f.Logout()
	} else {

		f.header.usernameLabel.Text = f.info.Session.UserData.Username
		f.HideAuthForm()
		f.context.ctx, f.context.cancel = context.WithCancel(context.Background())

		err = f.service.MakeSession(*f.info)
		if err != nil {
			f.ErrorHandler(err, LoginResponsibility)
		}
		f.service.Infof("Данные из сессии подгрузились. Успешная авторизация под пользователем %s", f.info.Session.UserData.FullName)
		f.flag.StopListen = true
		f.flag.FirstStart = true
		go f.Listen(f.context.ctx, *f.info)
	}

	f.setupCookieIntoEntry()
	go f.ClockUpdater()
	f.flag.DebugOn = false
	go func() {
		time.Sleep(2 * time.Second)
		f.flag.StopListen = false
	}()
	f.MainWindow.ShowAndRun()
	//f.PopUp.Text = widget.NewLabel("")
	//f.PopUp.PopUp = widget.NewPopUp(f.PopUp.Text, f.MainWindow.Canvas())

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
