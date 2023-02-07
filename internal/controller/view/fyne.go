package view

import (
	"fmt"
	"fptr/internal/services"
	"fptr/pkg/toml"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

type Flags struct {
	PrintOnKKTTicketCheckBox, PrintCheckBox, PrintOnPrinterTicketBox bool
}

type Selected struct {
}

type FyneApp struct {
	//приложение
	service     *services.Service
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

func NewFyneApp(a fyne.App) *FyneApp { //, service *services.Service
	return &FyneApp{
		application: a,
		//service:     service,
	}
}

func (f *FyneApp) StartApp() {
	f.ConfigureMainWindows()
	f.ConfigureAuthDialogForm()
	f.ConfigureWarningAlert()
	f.ConfigureSettingWindow()
	f.SetupCookie()
	f.mainWindow.ShowAndRun()
}

// + Главное окно
// ! Окно
func (f *FyneApp) NewMainWindow() {
	f.mainWindow = f.application.NewWindow("Ticket-Place")
}

func (f *FyneApp) NewSettingWindow() {
	content := container.NewVBox(
		container.New(layout.NewFormLayout(), widget.NewLabel("Путь к драйверу"), f.DriverSetting.DriverPathEntry),
		container.New(layout.NewFormLayout(), widget.NewLabel("Адрес сервера"), f.DriverSetting.DriverAddressEntry),
		container.New(layout.NewFormLayout(), widget.NewLabel("COM-порт ККТ"), f.DriverSetting.DriverComPortEntry),
		container.New(layout.NewFormLayout(), widget.NewLabel("Период опроса сервера"), f.DriverSetting.DriverPollingPeriodSelect),
		widget.NewLabel(""),
	)

	f.SettingWindow = dialog.NewCustomConfirm("Настройки приложения", "Сохранить", "Отменить", content, f.SettingWindowPressed, f.mainWindow)
}

func (f *FyneApp) ConfigureSettingWindow() {
	f.NewSettingWindow()
	f.SettingWindow.Resize(fyne.NewSize(500, 500))
}

func (f *FyneApp) SettingWindowPressed(choice bool) {
	settings := f.formDriverData()
	if choice {
		fmt.Println(toml.WriteToml(toml.DriverInfoPath, settings))
	}
}

func (f *FyneApp) ConfigureMainWindows() {
	f.NewMainWindow()
	f.ConfigureMainWindowAccordion()
	f.mainWindow.Resize(fyne.NewSize(500, 500))
	f.SetupMainWindowsContent()
}

func (f *FyneApp) SetupMainWindowsContent() {
	f.mainWindow.SetContent(
		container.NewVBox(
			f.ConfigureMainWindowHeader(),
			f.MainWindowAccordion,
		),
	)
}

// ! Контент

//? Контент Авторизации

func (f *FyneApp) ConfigureAuthDialogForm() {
	f.NewAuthForm()
	f.authForm.form.Resize(fyne.NewSize(500, 250))
	f.authForm.form.Show()
}

func (f *FyneApp) NewAuthForm() {
	f.authForm.loginEntry, f.authForm.passwordEntry = widget.NewEntry(), widget.NewPasswordEntry()
	f.authForm.settingButton = widget.NewButton("Настройки", func() {
		f.SettingWindow.Show()
	})
	var authFormItems []*widget.FormItem
	authFormItems = append(authFormItems,
		widget.NewFormItem("Логин", f.authForm.loginEntry),
		widget.NewFormItem("Пароль", f.authForm.passwordEntry),
		widget.NewFormItem("Настройки", f.authForm.settingButton),
	)
	f.authForm.form = dialog.NewForm("Авторизация", "Войти", "Выйти", authFormItems, f.Authorization, f.mainWindow)
}

func (f *FyneApp) Authorization(choice bool) { //! обработчик действий
	if choice {
		//f.header.usernameLabel.Text = f.authForm.loginEntry.Text
		//f.header.usernameLabel.Refresh()
		//f.ShowWarning("Ошибка доступа")
	} else {
		f.mainWindow.Close()
	}
}

//? Контент Header

func (f *FyneApp) NewMainWindowHeader() {
	f.header.usernameLabel = canvas.NewText("", color.White)
	f.header.localTimeLabel = canvas.NewText(time.Now().Format("2.01.2006 15:04:05"), color.White)
	f.header.printLastСheckButton = widget.NewButton("Напечатать последний чек", f.printLastCheckPressed)
	f.header.printXReportButton = widget.NewButton("Напечатать X-отчет", f.exitPressed)
	f.header.exitButton = widget.NewButton("Выйти", f.printXReportPressed)
}

func (f *FyneApp) ConfigureMainWindowHeader() *fyne.Container {
	f.NewMainWindowHeader()
	f.header.usernameLabel.TextSize = 18
	f.header.localTimeLabel.TextSize = 18
	f.header.printXReportButton.Importance = widget.MediumImportance
	f.header.localTimeLabel.Refresh()
	f.header.usernameLabel.Refresh()

	box := container.NewHBox(
		f.header.usernameLabel,
		f.header.localTimeLabel,
		f.header.printLastСheckButton,
		f.header.exitButton,
		f.header.printXReportButton,
	)
	return box
}

func (f *FyneApp) printLastCheckPressed() {
	//Напечатать последний чек
}

func (f *FyneApp) exitPressed() {
	//Механизм разлогинивания
}

func (f *FyneApp) printXReportPressed() {
	//Механизм напечатания X-отчета
}

func (f *FyneApp) NewPrintSettingsAccordionItem() {
	f.PrintSettingsItem.PrintCheck = widget.NewCheckWithData("Печатать чек", binding.BindBool(&f.flag.PrintCheckBox))
	f.PrintSettingsItem.PrintOnKKT = widget.NewCheckWithData("Печатать билет на кассе", binding.BindBool(&f.flag.PrintOnKKTTicketCheckBox))
	f.PrintSettingsItem.PrintOnPrinter = widget.NewCheckWithData("Печатать билет на принтере", binding.BindBool(&f.flag.PrintOnPrinterTicketBox))
	f.PrintSettingsItem.AdditionalText = widget.NewEntry()
	f.PrintSettingsItem.SetAdditionalText = widget.NewButton("Записать", f.SetAdditionalTextPressed)
}

func (f *FyneApp) ConfigurePrintSettingsAccordionItem() {
	f.NewPrintSettingsAccordionItem()
	f.PrintSettingsItem.AdditionalText.MultiLine = true
	f.PrintSettingsItem.AdditionalText.Resize(fyne.NewSize(200, 600))
	f.PrintSettingsItem.AdditionalText.Refresh()

	f.PrintSettingsItem.PrintSettingsAccordionItem = widget.NewAccordionItem("Параметры печати", container.NewVBox(
		container.NewHBox(
			f.PrintSettingsItem.PrintCheck,
			f.PrintSettingsItem.PrintOnKKT,
			f.PrintSettingsItem.PrintOnPrinter,
		), container.NewVBox(
			f.PrintSettingsItem.AdditionalText,
			f.PrintSettingsItem.SetAdditionalText,
		),
	))
}

func (f *FyneApp) SetAdditionalTextPressed() {

}

func (f *FyneApp) NewPrintsRefoundAndDepositsAccordionItem() {
	f.PrintsRefoundAndDeposits.RefoundEntry = widget.NewEntry()
	f.PrintsRefoundAndDeposits.RefoundFormItem = widget.NewFormItem("Возврат заказа", f.PrintsRefoundAndDeposits.RefoundEntry)
	f.PrintsRefoundAndDeposits.CashIncomeEntry = widget.NewEntry()
	f.PrintsRefoundAndDeposits.CashIncomeFormItem = widget.NewFormItem("Внесение наличных", f.PrintsRefoundAndDeposits.CashIncomeEntry)
	f.PrintsRefoundAndDeposits.PrintCheckEntry = widget.NewEntry()
	f.PrintsRefoundAndDeposits.PrintCheckFormItem = widget.NewFormItem("Печать заказа", f.PrintsRefoundAndDeposits.PrintCheckEntry)
	f.PrintsRefoundAndDeposits.RefoundForm = widget.NewForm(f.PrintsRefoundAndDeposits.RefoundFormItem)
	f.PrintsRefoundAndDeposits.CashIncomeForm = widget.NewForm(f.PrintsRefoundAndDeposits.CashIncomeFormItem)
	f.PrintsRefoundAndDeposits.PrintCheckForm = widget.NewForm(f.PrintsRefoundAndDeposits.PrintCheckFormItem)
}

func (f *FyneApp) ConfigurePrintsRefoundAndDepositsAccordionItem() {
	f.NewPrintsRefoundAndDepositsAccordionItem()
	f.PrintsRefoundAndDeposits.RefoundForm.SubmitText = "Вернуть"
	f.PrintsRefoundAndDeposits.RefoundForm.OnSubmit = f.RefoundOnSubmit
	f.PrintsRefoundAndDeposits.CashIncomeForm.SubmitText = "Внести"
	f.PrintsRefoundAndDeposits.CashIncomeForm.OnSubmit = f.CashIncomeOnSubmit
	f.PrintsRefoundAndDeposits.PrintCheckForm.SubmitText = "Печатать"
	f.PrintsRefoundAndDeposits.PrintCheckForm.OnSubmit = f.PrintCheckOnSubmit
	box := container.NewVBox(
		f.PrintsRefoundAndDeposits.RefoundForm, widget.NewLabel(""),
		f.PrintsRefoundAndDeposits.CashIncomeForm, widget.NewLabel(""),
		f.PrintsRefoundAndDeposits.PrintCheckForm, widget.NewLabel(""),
	)
	f.PrintsRefoundAndDeposits.RefoundAndDepositsAccordionItem = widget.NewAccordionItem("Печать заказа, возвраты и внесения", box)
}

func (f *FyneApp) CashIncomeOnSubmit() {

}

func (f *FyneApp) PrintCheckOnSubmit() {

}

func (f *FyneApp) RefoundOnSubmit() {

}

func (f *FyneApp) NewDriverSettingAccordionItem() {
	f.DriverSetting.DriverSettingButton = widget.NewButton("Открыть настройки принтера", f.DriverSettingButtonPressed)
	f.DriverSetting.DriverPrintHistoryButton = widget.NewButton("Открыть историю печати", f.DriverPrintHistoryButtonPressed)
	f.DriverSetting.DriverSettingLabel = widget.NewLabel("Настройки локального драйвера")

	f.DriverSetting.DriverComPortEntry = widget.NewEntry()
	f.DriverSetting.DriverPathEntry = widget.NewEntry()
	f.DriverSetting.DriverAddressEntry = widget.NewEntry()
	f.DriverSetting.DriverPollingPeriodSelect = widget.NewSelect(
		[]string{"1s", "2s", "3s", "4s", "5s", "10s", "15s"},
		f.DriverPollingPeriodSelected,
	)
	f.DriverSetting.DriverPollingPeriodSelect.Resize(fyne.NewSize(300, 300))
	f.DriverSetting.DriverPollingPeriodSelect.Refresh()
	f.DriverSetting.DriverKKTComFormItem = widget.NewFormItem("COM-порт кассы", f.DriverSetting.DriverComPortEntry)
	f.DriverSetting.DriverKKTPathFormItem = widget.NewFormItem("Путь к драйверу", f.DriverSetting.DriverPathEntry)
	f.DriverSetting.DriverApiAddressFormItem = widget.NewFormItem("Адрес сервера", f.DriverSetting.DriverAddressEntry)
	f.DriverSetting.DriverPollingPeriodFormItem = widget.NewFormItem("Период опроса сервера", f.DriverSetting.DriverPollingPeriodSelect)
	f.DriverSetting.DriverSettingForm = widget.NewForm(
		f.DriverSetting.DriverKKTComFormItem,
		f.DriverSetting.DriverKKTPathFormItem,
		f.DriverSetting.DriverApiAddressFormItem,
		f.DriverSetting.DriverPollingPeriodFormItem,
	)

}

func (f *FyneApp) ConfigureDriverSettingAccordionItem() {
	f.NewDriverSettingAccordionItem()
	f.DriverSetting.DriverSettingForm.SubmitText = "Подтвердить"
	f.DriverSetting.DriverSettingForm.OnSubmit = f.DriverSettingFormOnSubmit

	box := container.NewVBox(
		widget.NewLabel("Настройки принтера"),
		container.NewHBox(f.DriverSetting.DriverSettingButton, f.DriverSetting.DriverPrintHistoryButton),
		f.DriverSetting.DriverSettingLabel,
		f.DriverSetting.DriverSettingForm,
	)
	f.DriverSetting.DriverSettingAccordion = widget.NewAccordionItem("Настройки драйвера", box)
}

func (f *FyneApp) DriverSettingButtonPressed() {

}

func (f *FyneApp) DriverPrintHistoryButtonPressed() {

}

func (f *FyneApp) DriverSettingFormOnSubmit() {

}

func (f *FyneApp) DriverApiFormOnSubmit() {

}

func (f *FyneApp) DriverPollingPeriodSelected(selected string) {

}

func (f *FyneApp) NewMainWindowAccordion() {
	f.ConfigurePrintSettingsAccordionItem()
	f.ConfigurePrintsRefoundAndDepositsAccordionItem()
	f.ConfigureDriverSettingAccordionItem()
	f.MainWindowAccordion = widget.NewAccordion(
		f.PrintSettingsItem.PrintSettingsAccordionItem, f.PrintsRefoundAndDeposits.RefoundAndDepositsAccordionItem, f.DriverSetting.DriverSettingAccordion,
	)
}

func (f *FyneApp) ConfigureMainWindowAccordion() {
	f.NewMainWindowAccordion()
}

func (f *FyneApp) NewWarningAlert() {
	f.Warning.WarningText = canvas.NewText("", color.White)
	textError := canvas.NewText("Возникла ошибка во время выполнения: ", color.White)

	box := container.NewVBox(container.NewHBox(textError), container.NewHBox(f.Warning.WarningText), container.NewHBox(widget.NewLabel("")))
	f.Warning.WarningWindow = dialog.NewCustom("Ошибка", "Понятно", box, f.mainWindow)
	f.Warning.WarningWindow.SetOnClosed(f.WarningPressed)
}

func (f *FyneApp) ConfigureWarningAlert() {
	f.NewWarningAlert()
	f.Warning.WarningWindow.Resize(fyne.NewSize(300, 100))
	f.Warning.WarningWindow.Hide()
}

func (f *FyneApp) WarningPressed() {

}

func (f *FyneApp) ShowWarning(err string) {
	f.Warning.WarningText.Text = err
	f.Warning.WarningText.Refresh()
	f.Warning.WarningWindow.Show()
}
