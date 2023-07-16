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
	"net/url"
	"time"
)

const (
	PageA4        = "A4"
	PageA5        = "A5"
	PagePortrait  = "Portrait"
	PageLandscape = "Landscape"
)

func (f *FyneApp) NewMainWindow() {
	f.MainWindow = f.application.NewWindow("Ticket Place")
	f.MainWindow.SetMaster()
}

func (f *FyneApp) NewPrinterSettings() {
	f.PrinterSettings.CheckPrinterStatus = widget.NewButton("Проверить статус службы", f.CheckStatusAction)
	f.PrinterSettings.RunPrinterService = widget.NewButton("Запустить службу печати", f.RunPrintServiceAction)
	f.PrinterSettings.SelectPrinter = widget.NewSelect([]string{}, f.SelectPrinterAction)
	f.PrinterSettings.StatusImage = canvas.NewImageFromResource(theme.QuestionIcon())
	f.PrinterSettings.StopPrinterService = widget.NewButton("", f.StopPrinterServiceAction)
	f.PrinterSettings.GetListOfPrinters = widget.NewButton("Обновить список принтеров", f.GetListOfPrintersAction)
	f.PrinterSettings.PrinterServiceAddress = widget.NewEntry()
	f.PrinterSettings.PrinterServiceAddress.PlaceHolder = "Адрес службы печати"
}

func (f *FyneApp) StopPrinterServiceAction() {

}

func (f *FyneApp) GetListOfPrintersAction() {
	printers, err := f.service.GetListOfPrinters(f.info.AppConfig.Driver)
	if err != nil {
		return
	}
	f.PrinterSettings.SelectPrinter.Options = printers
	f.PrinterSettings.SelectPrinter.Refresh()
}

func (f *FyneApp) RunPrintServiceAction() {

}

func (f *FyneApp) SelectPrinterAction(change string) {

}

func (f *FyneApp) CheckStatusAction() {
	if err := f.service.Ping(f.info.AppConfig.Driver); err != nil {
		f.PrinterSettings.StatusImage.Resource = theme.ErrorIcon()
		f.PrinterSettings.StatusImage.Refresh()
		return
	}

	f.PrinterSettings.StatusImage.Resource = theme.ConfirmIcon()
	f.PrinterSettings.StatusImage.Refresh()

}

func (f *FyneApp) NewSettingWindow() {
	ButtonSeparatorText := canvas.NewText("Основные функции", theme.ForegroundColor())
	f.DriverSetting.PrintLastButton = widget.NewButtonWithIcon("Напечатать последний чек", theme.DocumentIcon(), f.printLastCheckPressedFromKKT)
	f.DriverSetting.ErrorSoundButton = widget.NewButtonWithIcon("Выключить звук оповещений", theme.MediaPauseIcon(), func() {
		if f.flag.SoundError {
			f.DriverSetting.ErrorSoundButton.Icon = theme.MediaPlayIcon()
			f.DriverSetting.ErrorSoundButton.Text = "Включить звук оповещений"
			f.flag.SoundError = false
		} else {
			f.DriverSetting.ErrorSoundButton.Icon = theme.MediaPauseIcon()
			f.DriverSetting.ErrorSoundButton.Text = "Выключить звук оповещений"
			f.flag.SoundError = true
		}
		f.DriverSetting.ErrorSoundButton.Refresh()
	})
	ButtonSeparatorText.TextSize = 18
	ButtonSeparatorText.Alignment = fyne.TextAlignCenter
	f.DriverSetting.DriverUpdatePath = widget.NewEntry()
	f.DriverSetting.CloseShiftButton = widget.NewButtonWithIcon("Закрыть смену", theme.CancelIcon(), f.CloseShift)
	f.DriverSetting.PrinterSettings = widget.NewButtonWithIcon("Настройки принтера", theme.DocumentPrintIcon(), f.OpenPrinterSettings)
	content := container.NewVBox(
		container.New(layout.NewFormLayout(), widget.NewLabel("Путь к драйверу"), f.DriverSetting.DriverPathEntry),
		container.New(layout.NewFormLayout(), widget.NewLabel("Адрес сервера"), f.DriverSetting.DriverAddressEntry),
		container.New(layout.NewFormLayout(), widget.NewLabel("COM-порт ККТ"), f.DriverSetting.DriverComPortEntry),
		container.New(layout.NewFormLayout(), widget.NewLabel("Период опроса сервера"), f.DriverSetting.DriverPollingPeriodSelect),
		container.New(layout.NewFormLayout(), widget.NewLabel("Длительность опроса"), f.DriverSetting.DriverTimeoutSelect),
		container.New(layout.NewFormLayout(), widget.NewLabel("Источник обновления"), f.DriverSetting.DriverUpdatePath),

		widget.NewSeparator(),
		ButtonSeparatorText,
		container.NewGridWithColumns(2, f.DriverSetting.CloseShiftButton, f.DriverSetting.ErrorSoundButton, f.DriverSetting.PrintLastButton, f.DriverSetting.PrinterSettings),
	)
	f.SettingWindow = dialog.NewCustomConfirm("Настройки приложения", "Сохранить", "Отменить", content, f.SettingWindowPressed, f.MainWindow)
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
	f.authForm.form = dialog.NewForm("Авторизация", "Войти", "Выйти", authFormItems, f.AuthorizationPressed, f.MainWindow)
}

func (f *FyneApp) NewMainWindowHeader() {
	f.header.usernameLabel = canvas.NewText("", theme.ForegroundColor())
	f.header.localTimeLabel = canvas.NewText(time.Now().Format("2.01.2006 15:04:05"), theme.ForegroundColor())

	f.header.listenerStatus.listenerToolbarItem = widget.NewToolbarAction(theme.CancelIcon(), f.listenerStatusAction)
	f.header.listenerStatus.listenerToolbar = widget.NewToolbar(
		f.header.listenerStatus.listenerToolbarItem,
	)

}

func (f *FyneApp) NewPrintSettingsContainer() {
	f.PrintSettingsItem.exitButton = widget.NewButton("Выйти", f.exitPressed)
	f.PrintSettingsItem.exitAndCloseShiftButton = widget.NewButton("Выйти и закрыть смену", f.exitAndCloseShiftButtonPressed)
	f.PrintSettingsItem.PrintCheck = widget.NewCheckWithData("Печатать чек", binding.BindBool(&f.flag.printCheckBox.PrintCheckBox))
	f.PrintSettingsItem.PrintOnKKT = widget.NewCheckWithData("Печатать билет на кассе", binding.BindBool(&f.flag.printCheckBox.PrintOnKKTTicketCheckBox))
	f.PrintSettingsItem.PrintOnPrinter = widget.NewCheckWithData("Печатать билет на принтере", binding.BindBool(&f.flag.printCheckBox.PrintOnPrinterTicketBox))
	f.PrintSettingsItem.AdditionalText = widget.NewEntry() //widget.NewEntry()
	f.PrintSettingsItem.SetAdditionalText = widget.NewButton("Записать", f.SetAdditionalTextPressed)
	f.PrintSettingsItem.printLastСheckButton = widget.NewButton("Напечатать последний чек", f.printLastCheckPressedFromKKT)
	f.PrintSettingsItem.printXReportButton = widget.NewButton("Напечатать X-отчет", f.printXReportPressed)
	f.PrintSettingsItem.CashIncomeEntry = widget.NewEntry()
	f.PrintSettingsItem.CashIncomeFormItem = widget.NewFormItem("Внесение наличных", f.PrintSettingsItem.CashIncomeEntry)
	f.PrintSettingsItem.CashIncomeForm = widget.NewForm(f.PrintSettingsItem.CashIncomeFormItem)
	f.PrintSettingsItem.reconnectButton = widget.NewButton("Восстановить соединение с кассой", f.OpenConnection)
	f.PrintSettingsItem.PageSizeRadioGroup = widget.NewRadioGroup([]string{PageA4, PageA5}, func(s string) {
		f.flag.pageParams.PageSize = s
	})
	f.PrintSettingsItem.PageSizeRadioGroup.Horizontal = true
	f.PrintSettingsItem.PageOrientationRadioGroup = widget.NewRadioGroup([]string{PageLandscape, PagePortrait}, func(s string) {
		f.flag.pageParams.PageOrientation = s

	})
	f.PrintSettingsItem.PageOrientationRadioGroup.Horizontal = true
}

func (f *FyneApp) NewPrintsRefoundAndDepositsAccordionItem() {
	f.PrintsRefoundAndDeposits.RefoundEntry = widget.NewEntry()
	f.PrintsRefoundAndDeposits.RefoundFormItem = widget.NewFormItem("Возврат заказа", f.PrintsRefoundAndDeposits.RefoundEntry)

	f.PrintsRefoundAndDeposits.PrintCheckEntry = widget.NewEntry()
	f.PrintsRefoundAndDeposits.PrintCheckFormItem = widget.NewFormItem("Печать заказа", f.PrintsRefoundAndDeposits.PrintCheckEntry)
	f.PrintsRefoundAndDeposits.RefoundForm = widget.NewForm(f.PrintsRefoundAndDeposits.RefoundFormItem)

	f.PrintsRefoundAndDeposits.AdminEntry = widget.NewPasswordEntry()
	f.PrintsRefoundAndDeposits.AdminFormItem = widget.NewFormItem("Введите пароль", f.PrintsRefoundAndDeposits.AdminEntry)
	f.PrintsRefoundAndDeposits.AdminForm = widget.NewForm(f.PrintsRefoundAndDeposits.AdminFormItem)

	f.PrintsRefoundAndDeposits.PrintCheckForm = widget.NewForm(f.PrintsRefoundAndDeposits.PrintCheckFormItem)
}

func (f *FyneApp) NewDriverSettingAccordionItem() {
	f.DriverSetting.DriverSettingButton = widget.NewButton("Открыть настройки принтера", f.DriverSettingButtonPressed)
	f.DriverSetting.DriverPrintHistoryButton = widget.NewButton("Открыть историю печати", f.DriverPrintHistoryButtonPressed)
	f.DriverSetting.DriverSettingLabel = widget.NewLabel("Настройки локального драйвера")

	f.DriverSetting.DriverComPortEntry = widget.NewEntry()
	f.DriverSetting.DriverPathEntry = widget.NewEntry()
	f.DriverSetting.DriverAddressEntry = widget.NewEntry()
	f.DriverSetting.DriverTimeoutSelect = widget.NewSelect([]string{"2s", "5s", "10s", "20s", "40s", "60s", "70s", "80s", "90s", "100s", "120s"}, nil)
	f.DriverSetting.DriverPollingPeriodSelect = widget.NewSelect(
		[]string{"1s", "2s", "3s", "4s", "5s", "10s", "15s"},
		f.DriverPollingPeriodSelected,
	)
	f.DriverSetting.DriverTimeoutSelect.Selected = "20s"
	f.DriverSetting.DriverPollingPeriodSelect.Selected = "3s"

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
	f.Warning.WarningWindow = dialog.NewCustom("Ошибка", "Понятно", box, f.MainWindow)
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

func (f *FyneApp) NewCriticalAlert() {
	f.CriticalError.ErrorWindow = f.application.NewWindow("Критическая ошибка")
	f.CriticalError.ErrorText = canvas.NewText("", theme.ForegroundColor())
	f.CriticalError.ErrorText.Alignment = fyne.TextAlignCenter
	f.CriticalError.ErrorText.TextSize = 18
	f.CriticalError.ErrorConfirmButton = widget.NewButtonWithIcon("Выйти", theme.ConfirmIcon(), func() {
		f.application.Quit()
	})
	f.CriticalError.AdditionalText = canvas.NewText("", theme.ForegroundColor())
	f.CriticalError.ErrorWindow.SetIcon(theme.WarningIcon())
	f.CriticalError.ErrorLinkButton = widget.NewHyperlink("", &url.URL{})
	f.CriticalError.ErrorWindow.SetCloseIntercept(func() {
		f.application.Quit()
	})
}

func (f *FyneApp) NewProgresser() {
	f.Reconnector.ProgressBar = widget.NewProgressBarInfinite()
	f.Reconnector.ProgressConfirm = widget.NewButton("Попробовать снова", f.ProgresserPressedConfirm)
	f.Reconnector.ProgressDismiss = widget.NewButton("Выйти", f.ProgresserPressedDismiss)
	f.Reconnector.ProgresserText = canvas.NewText("Попытка подключиться к серверу номер 1", theme.ForegroundColor())
	box := container.NewVBox(container.NewCenter(f.Reconnector.ProgresserText), f.Reconnector.ProgressBar, container.NewCenter(container.NewHBox(f.Reconnector.ProgressConfirm, f.Reconnector.ProgressDismiss)))
	f.Reconnector.Progresser = f.application.NewWindow("Подключаемся к серверу")
	f.Reconnector.Progresser.SetIcon(theme.AccountIcon())
	f.Reconnector.Progresser.SetContent(box)
	f.Reconnector.Progresser.Hide()
}

func (f *FyneApp) NewToolbar() {
	//f.Toolbar.Logo = canvas.NewImageFromFile("")
	//f.Toolbar.Logo.FillMode = canvas.ImageFillStretch
	f.Toolbar.Info = widget.NewToolbarAction(theme.InfoIcon(), f.ToolbarInfoPressed)
}

func (f *FyneApp) ProgresserPressedConfirm() {
	f.ShowProgresser()
}

func (f *FyneApp) ProgresserPressedDismiss() {
	f.Reconnector.Progresser.Hide()
}

func (f *FyneApp) NewAboutDialog() {
	f.AboutDialog.Img = canvas.NewImageFromFile("./content/system/icon/logo.png")
	f.AboutDialog.Version = canvas.NewText("", theme.ForegroundColor())
	f.AboutDialog.Information = canvas.NewText("", theme.ForegroundColor())
	f.AboutDialog.CheckUpdateButton = widget.NewButtonWithIcon("Проверить обновления", theme.DownloadIcon(), f.CheckUpdateAction)
}

func (f *FyneApp) NewPrintDoubleConfirm() {
	f.PrintDoubleConfirm.PDConfirm = widget.NewButton("Да", func() {
		f.flag.Waiter <- true
		f.PrintDoubleConfirm.Window.Hide()
	})
	f.PrintDoubleConfirm.PDDismiss = widget.NewButton("Нет", func() {
		f.flag.Waiter <- false
		f.PrintDoubleConfirm.Window.Hide()
	})
	f.PrintDoubleConfirm.Text = canvas.NewText("", theme.ForegroundColor())
}

func (f *FyneApp) NewProgressAction() { //+ Внедрение Progress Bar
	f.ProgressAction.ProgressValue = binding.NewFloat()
	f.ProgressAction.ProgressStatus = binding.NewString()
	f.ProgressAction.ProgressStatus.Set("Текущий статус: нет задач")
	f.service.SetProgressData(f.ProgressAction.ProgressValue, f.ProgressAction.ProgressStatus)
	f.ProgressAction.StatusText = widget.NewLabelWithData(f.ProgressAction.ProgressStatus)
	f.ProgressAction.Progress = widget.NewProgressBarWithData(f.ProgressAction.ProgressValue)
}
