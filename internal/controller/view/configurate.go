package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const iconPath = "./content/system/icon/main.png"
const adminPassword = "ticket_admin"

func (f *FyneApp) ConfigureFlags() {
	f.flag.SoundError = true
	f.flag.ProgressWorking = false
	f.flag.AuthJustHide = false
}

func (f *FyneApp) ConfigureSettingWindow() {
	f.NewSettingWindow()
	f.ConfigureFlags()
	f.SettingWindow.Resize(fyne.NewSize(500, 500))
}

func (f *FyneApp) ConfigureMainWindows() {
	f.NewMainWindow()
	f.ConfigureMainWindowAccordion()
	f.NewProgressAction()
	f.MainWindow.Resize(fyne.NewSize(600, 600))
	icoResource, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
		f.service.Errorf("Ошибка установки иконки: %v", err)
	} else {
		f.MainWindow.SetIcon(icoResource)
	}
	//f.application.Settings().SetTheme(theme.DarkTheme())
	f.ProgressAction.ProgressBox = container.NewVBox(f.ProgressAction.StatusText, f.ProgressAction.Progress)
	content := container.NewBorder(f.Toolbar.Box, nil, nil, nil, container.NewVBox(
		f.ConfigureMainWindowHeader(),
		f.PrintSettingsItem.PrintSettingsContainer,
		f.MainWindowAccordion,
		f.ProgressAction.ProgressBox,
	))
	f.MainWindow.SetContent(content)
}

func (f *FyneApp) ConfigureAuthDialogForm() {
	f.NewAuthForm()
	f.authForm.form.Resize(fyne.NewSize(500, 250))
	f.authForm.form.Show()
}

func (f *FyneApp) ConfigureProgresser() {
	f.NewProgresser()
	f.Reconnector.Progresser.Resize(fyne.NewSize(350, 100))
	f.Reconnector.Progresser.SetFixedSize(true)

}

func (f *FyneApp) ConfigureMainWindowHeader() *fyne.Container {
	f.NewMainWindowHeader()

	f.header.usernameLabel.TextSize = 18
	f.header.localTimeLabel.TextSize = 18

	f.header.localTimeLabel.Refresh()
	f.header.usernameLabel.Refresh()

	box := container.NewHBox(
		f.header.usernameLabel,
		f.header.localTimeLabel,
		layout.NewSpacer(),
		//f.header.listenerStatus.listenerToolbar,
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

	f.PrintSettingsItem.printXReportButton.Importance = widget.MediumImportance

	f.PrintSettingsItem.CashIncomeForm.SubmitText = "Внести"
	f.PrintSettingsItem.CashIncomeForm.OnSubmit = f.CashIncomeOnSubmit
	_ = container.NewVBox(
		container.NewGridWithColumns(2, f.PrintSettingsItem.printLastСheckButton,
			f.PrintSettingsItem.printXReportButton), f.PrintSettingsItem.CashIncomeForm,
	)

	textCont := container.NewMax(f.PrintSettingsItem.AdditionalText)
	textCont.Resize(fyne.NewSize(600, 600))
	f.PrintSettingsItem.PrintCheck.SetChecked(true)
	f.PrintSettingsItem.PrintSettingsContainer =
		container.NewVBox(
			container.NewGridWithColumns(2, f.PrintSettingsItem.exitButton,
				f.PrintSettingsItem.exitAndCloseShiftButton),
			container.NewCenter(
				container.NewHBox(f.PrintSettingsItem.PrintCheck,
					f.PrintSettingsItem.PrintOnKKT,
					f.PrintSettingsItem.PrintOnPrinter),
			),
			widget.NewSeparator(),
			container.NewGridWithColumns(2, f.PrintSettingsItem.printLastСheckButton, f.PrintSettingsItem.printXReportButton),
			//container.NewGridWithColumns(1, f.PrintSettingsItem.reconnectButton),
			f.PrintSettingsItem.CashIncomeForm,
		)

}

func (f *FyneApp) ConfigurePrintsRefoundAndDepositsAccordionItem() {
	f.NewPrintsRefoundAndDepositsAccordionItem()
	f.PrintsRefoundAndDeposits.RefoundForm.SubmitText = "Вернуть"
	f.PrintsRefoundAndDeposits.RefoundForm.OnSubmit = f.RefoundOnSubmit
	f.PrintsRefoundAndDeposits.RefoundForm.Hide()
	f.PrintsRefoundAndDeposits.PrintCheckForm.SubmitText = "Печатать"
	f.PrintsRefoundAndDeposits.PrintCheckForm.OnSubmit = f.PrintCheckOnSubmit
	f.PrintsRefoundAndDeposits.PrintCheckForm.Hide()
	f.PrintsRefoundAndDeposits.AdminForm.SubmitText = "Войти"
	f.PrintsRefoundAndDeposits.AdminForm.OnSubmit = func() {
		if f.PrintsRefoundAndDeposits.AdminEntry.Text == adminPassword {
			go f.BlockItemControl()
			f.PrintsRefoundAndDeposits.AdminEntry.Text = ""
			f.PrintsRefoundAndDeposits.AdminForm.Hide()
			f.PrintsRefoundAndDeposits.RefoundForm.Show()
			f.PrintsRefoundAndDeposits.PrintCheckForm.Show()
		} else {
			dialog.ShowInformation("Неправильный пароль", "Вы ввели неправильный пароль. Попробуйте снова!", f.MainWindow)
		}
	}
	box := container.NewVBox(
		f.PrintsRefoundAndDeposits.AdminForm,
		f.PrintsRefoundAndDeposits.RefoundForm, widget.NewLabel(""),
		f.PrintsRefoundAndDeposits.PrintCheckForm, widget.NewLabel(""),
	)
	f.PrintsRefoundAndDeposits.RefoundAndDepositsAccordionItem = widget.NewAccordionItem("Печать чеков продаж и возвратов", box)
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

func (f *FyneApp) ConfigurateErrorAlert() {
	f.NewErrorAlert()
	f.Error.ErrorWindow.Resize(fyne.NewSize(400, 100))

	f.Error.ErrorWindow.Hide()
	title := canvas.NewText("Во время исполнения произошла ошибка", theme.ForegroundColor())
	title.TextSize = 20
	title.Alignment = fyne.TextAlignCenter
	box := container.NewBorder(title, container.NewCenter(f.Error.ErrorConfirmButton), container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50))), container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50))), f.Error.ErrorText)
	f.Error.ErrorWindow.SetContent(box)
}

func (f *FyneApp) ConfigurateCriticalErrorAlert() {
	f.NewCriticalAlert()
	f.CriticalError.ErrorWindow.Resize(fyne.NewSize(400, 100))
	f.CriticalError.ErrorWindow.Hide()
	title := canvas.NewText("Во время исполнения произошла критическая ошибка", theme.ForegroundColor())
	title.TextSize = 20
	title.Alignment = fyne.TextAlignCenter
	title.Hide()
	f.CriticalError.ErrorLinkButton.Alignment = fyne.TextAlignCenter
	f.CriticalError.AdditionalText.Alignment = fyne.TextAlignCenter
	f.CriticalError.AdditionalText.TextSize = 18
	image := canvas.NewImageFromResource(theme.ErrorIcon())
	image.FillMode = canvas.ImageFillContain
	boxImage := container.NewGridWrap(fyne.NewSize(130, 130), image)
	boxCenter := container.NewCenter(boxImage)
	box := container.NewBorder(title, container.NewCenter(f.CriticalError.ErrorConfirmButton), container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50))), container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50))), container.NewVBox(boxCenter, f.CriticalError.AdditionalText, f.CriticalError.ErrorLinkButton))
	f.CriticalError.ErrorWindow.SetContent(box)
}

func (f *FyneApp) ConfigurateToolbar() {
	f.NewToolbar()

	//boxImage := container.NewGridWrap(fyne.NewSize(165, 40), f.Toolbar.Logo)
	f.Toolbar.Toolbar = widget.NewToolbar(widget.NewToolbarSpacer(), f.Toolbar.Info)
	f.Toolbar.Box = container.New(layout.NewMaxLayout(), f.Toolbar.Toolbar)

}

func (f *FyneApp) ConfigurateAboutDialogWindow() {
	f.NewAboutDialog()
	f.AboutDialog.Version.Text = fmt.Sprintf("Версия текущего ПО: %s", f.AppInfo.version)
	f.AboutDialog.Version.Alignment = fyne.TextAlignCenter
	f.AboutDialog.Information.Text = fmt.Sprintf("Дополнительная информация о ПО: ")
	f.AboutDialog.Information.Text = ""
	//f.AboutDialog.Information.Hide()
	f.AboutDialog.Version.Alignment = fyne.TextAlignCenter
	box := container.NewVBox(f.AboutDialog.Version, container.NewCenter(f.AboutDialog.CheckUpdateButton), widget.NewSeparator())
	f.AboutDialog.Dialog = dialog.NewCustom("О программе", "Закрыть", box, f.MainWindow)

	f.AboutDialog.Dialog.Hide()
}

func (f *FyneApp) ConfigurateDoubleConfirm() {
	f.NewPrintDoubleConfirm()
	box := container.NewVBox(
		widget.NewLabel("Вы пытаетесь распечатать чек, который был напечатан в течении 15 секунд. Повторить печать?"),
		f.PrintDoubleConfirm.PDConfirm,
		f.PrintDoubleConfirm.PDDismiss,
	)
	f.PrintDoubleConfirm.Window = f.application.NewWindow("Повторная печать")
	f.PrintDoubleConfirm.Window.SetContent(box)
	f.PrintDoubleConfirm.Window.Hide()
}
