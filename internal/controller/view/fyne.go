package view

import (
	"context"
	"fptr/internal/entities"
	"fptr/internal/services"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Flags struct {
	PrintOnKKTTicketCheckBox, PrintCheckBox, PrintOnPrinterTicketBox bool
}

type Selected struct {
}

type FyneApp struct {
	info    *entities.Info
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
		usernameLabel        *canvas.Text
		localTimeLabel       *canvas.Text
		printLastСheckButton *widget.Button
		exitButton           *widget.Button
		printXReportButton   *widget.Button
	}
	PrintSettingsItem struct {
		PrintSettingsAccordionItem *widget.AccordionItem
		PrintCheck                 *widget.Check
		PrintOnKKT                 *widget.Check
		PrintOnPrinter             *widget.Check
		AdditionalText             *widget.Entry
		SetAdditionalText          *widget.Button
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

	MainWindowAccordion *widget.Accordion

	//! Окно ошибок
	Warning struct {
		WarningWindow dialog.Dialog
		WarningText   *canvas.Text
	}

	AlertWindow   dialog.Dialog
	SettingWindow dialog.Dialog
	//Флаги
	flag Flags
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
	}

	if f.info.Session.IsDead() {
		f.UpdateSession(entities.SessionInfo{})
	} else {
		f.header.usernameLabel.Text = f.info.Session.UserData.Username
		f.authForm.form.Hide()
		f.context.ctx, f.context.cancel = context.WithCancel(context.Background())
		go f.Listen(f.context.ctx)
	}

	f.setupCookieIntoEntry()
	f.mainWindow.ShowAndRun()
}

// + Главное окно
// ! Окно
