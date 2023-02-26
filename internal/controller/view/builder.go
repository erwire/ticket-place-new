package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"time"
)

func (f *FyneApp) NewMainWindow() {
	f.mainWindow = f.application.NewWindow("Ticket-Place")
	f.mainWindow.SetMaster()
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
	f.authForm.form = dialog.NewForm("Авторизация", "Войти", "Выйти", authFormItems, f.AuthorizationPressed, f.mainWindow)
}

func (f *FyneApp) NewMainWindowHeader() {
	f.header.usernameLabel = canvas.NewText("", theme.ForegroundColor())
	f.header.localTimeLabel = canvas.NewText(time.Now().Format("2.01.2006 15:04:05"), theme.ForegroundColor())
	f.header.printLastСheckButton = widget.NewButton("Напечатать последний чек", f.printLastCheckPressed)
	f.header.printXReportButton = widget.NewButton("Напечатать X-отчет", f.printXReportPressed)
	f.header.exitButton = widget.NewButton("Выйти", f.exitPressed)
	f.header.exitAndCloseShiftButton = widget.NewButton("Выйти и закрыть смену", f.exitAndCloseShiftButtonPressed)
	f.header.listenerStatus.listenerToolbarItem = widget.NewToolbarAction(theme.CancelIcon(), f.listenerStatusAction)
	f.header.listenerStatus.listenerToolbar = widget.NewToolbar(
		f.header.listenerStatus.listenerToolbarItem,
	)

}

func (f *FyneApp) NewPrintSettingsContainer() {
	f.PrintSettingsItem.PrintCheck = widget.NewCheckWithData("Печатать чек", binding.BindBool(&f.flag.PrintCheckBox))
	f.PrintSettingsItem.PrintOnKKT = widget.NewCheckWithData("Печатать билет на кассе", binding.BindBool(&f.flag.PrintOnKKTTicketCheckBox))
	f.PrintSettingsItem.PrintOnPrinter = widget.NewCheckWithData("Печатать билет на принтере", binding.BindBool(&f.flag.PrintOnPrinterTicketBox))
	f.PrintSettingsItem.AdditionalText = widget.NewEntry() //widget.NewEntry()
	f.PrintSettingsItem.SetAdditionalText = widget.NewButton("Записать", f.SetAdditionalTextPressed)
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

func (f *FyneApp) NewMainWindowAccordion() {
	f.ConfigurePrintSettingsContainer()
	f.ConfigurePrintsRefoundAndDepositsAccordionItem()
	f.ConfigureDriverSettingAccordionItem()
	f.MainWindowAccordion = widget.NewAccordion(
		f.PrintsRefoundAndDeposits.RefoundAndDepositsAccordionItem, f.DriverSetting.DriverSettingAccordion,
	)
}

func (f *FyneApp) NewWarningAlert() {
	f.Warning.WarningText = canvas.NewText("", theme.ForegroundColor())
	textError := canvas.NewText("Возникла ошибка во время выполнения: ", theme.ForegroundColor())

	box := container.NewVBox(container.NewHBox(textError), container.NewHBox(f.Warning.WarningText), container.NewHBox(widget.NewLabel("")))
	f.Warning.WarningWindow = dialog.NewCustom("Ошибка", "Понятно", box, f.mainWindow)
	f.Warning.WarningWindow.SetOnClosed(f.WarningPressed)
}

func (f *FyneApp) NewErrorAlert() {
	f.Error.ErrorWindow = f.application.NewWindow("Ошибка")
	f.Error.ErrorText = canvas.NewText("", theme.ForegroundColor())
	f.Error.ErrorText.Alignment = fyne.TextAlignCenter
	f.Error.ErrorText.TextSize = 18
	f.Error.ErrorConfirmButton = widget.NewButtonWithIcon("Хорошо", theme.ConfirmIcon(), func() {
		f.Error.ErrorWindow.Hide()
	})
	f.Error.ErrorWindow.SetIcon(theme.WarningIcon())
	f.Error.ErrorWindow.SetCloseIntercept(func() {
		f.Error.ErrorWindow.Hide()
	})
}
