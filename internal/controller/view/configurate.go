package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const iconPath = "./content/system/icon/main.png"

func (f *FyneApp) ConfigureSettingWindow() {
	f.NewSettingWindow()
	f.SettingWindow.Resize(fyne.NewSize(500, 500))
}

func (f *FyneApp) ConfigureMainWindows() {
	f.NewMainWindow()
	f.ConfigureMainWindowAccordion()
	f.mainWindow.Resize(fyne.NewSize(600, 600))
	icoResource, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
		f.service.Errorf("Ошибка установки иконки: %v", err)
	} else {
		f.mainWindow.SetIcon(icoResource)
	}
	//f.application.Settings().SetTheme(theme.DarkTheme())
	f.mainWindow.SetContent(
		container.NewVBox(
			f.ConfigureMainWindowHeader(),
			f.PrintSettingsItem.PrintSettingsContainer,
			f.MainWindowAccordion,
		),
	)
}

func (f *FyneApp) ConfigureAuthDialogForm() {
	f.NewAuthForm()
	f.authForm.form.Resize(fyne.NewSize(500, 250))
	f.authForm.form.Show()
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
		f.header.exitAndCloseShiftButton,
		f.header.printXReportButton,
		f.header.listenerStatus.listenerToolbar,
	)
	return box
}

func (f *FyneApp) ConfigurePrintSettingsContainer() {
	f.NewPrintSettingsContainer()
	//f.PrintSettingsItem.AdditionalText.Wrapping = fyne.TextWrapBreak
	//f.PrintSettingsItem.AdditionalText.MultiLine = true
	//f.PrintSettingsItem.AdditionalText.Resize(fyne.NewSize(300, 300))
	//f.PrintSettingsItem.AdditionalText.SetPlaceHolder("Сообщение")
	//f.PrintSettingsItem.AdditionalText.Refresh()
	textCont := container.NewMax(f.PrintSettingsItem.AdditionalText)
	textCont.Resize(fyne.NewSize(600, 600))
	f.PrintSettingsItem.PrintCheck.SetChecked(true)
	f.PrintSettingsItem.PrintSettingsContainer = container.NewVBox(
		//textCont,
		container.NewHBox(
			f.PrintSettingsItem.PrintCheck,
			f.PrintSettingsItem.PrintOnKKT,
			f.PrintSettingsItem.PrintOnPrinter,
			//f.PrintSettingsItem.SetAdditionalText,
		),
	)
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

func (f *FyneApp) ConfigureDriverSettingAccordionItem() {
	f.NewDriverSettingAccordionItem()
	f.DriverSetting.DriverSettingForm.SubmitText = "Подтвердить"
	f.DriverSetting.DriverSettingForm.OnSubmit = f.DriverSettingFormOnSubmit

	box := container.NewVBox(
		widget.NewLabel("Настройки принтера"),
		container.NewHBox(f.DriverSetting.DriverSettingButton, f.DriverSetting.DriverPrintHistoryButton),
		//f.DriverSetting.DriverSettingLabel,
		//f.DriverSetting.DriverSettingForm,
	)
	f.DriverSetting.DriverSettingAccordion = widget.NewAccordionItem("Настройки драйвера", box)
}

func (f *FyneApp) ConfigureMainWindowAccordion() {
	f.NewMainWindowAccordion()
}

func (f *FyneApp) ConfigureWarningAlert() {
	f.NewWarningAlert()
	f.Warning.WarningWindow.Resize(fyne.NewSize(300, 100))
	f.Warning.WarningWindow.Hide()
}
