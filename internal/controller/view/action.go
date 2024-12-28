package view

import (
	"fmt"
	"fptr/internal/entities"
	"fptr/pkg/toml"
	"fyne.io/fyne/v2/theme"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

func (f *FyneApp) CashIncomeOnSubmit() {
	incomeStr := f.PrintSettingsItem.CashIncomeEntry.Text
	income, err := strconv.ParseFloat(incomeStr, 32)
	if err != nil {
		f.ShowWarning("Некорректные данные в поле ввода суммы")
		return
	}
	f.service.Infof("Запрос на внесение: %f руб.", income)

	if err := f.service.CashIncome(income); err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) PrintCheckOnSubmit() {
	id := f.PrintsRefoundAndDeposits.PrintCheckEntry.Text
	if id == "" {
		f.ShowWarning("Пожалуйста, вставьте значение в поле номера чека")
		return
	}
	f.service.Infof("Запрос на печать заказа с номером %s", id)

	if err := f.service.PrintSell(*f.info, id, nil); err != nil {
		f.ErrorHandler(err, SellResponsibility)
		return
	}
}

func (f *FyneApp) RefoundOnSubmit() {
	id := f.PrintsRefoundAndDeposits.RefoundEntry.Text
	if id == "" {
		f.ShowWarning("Пожалуйста, вставьте значение в поле возврата")
		return
	}
	f.service.Infof("Запрос на печать заказа с номером %s", id)

	if err := f.service.PrintRefoundFromSell(*f.info, id); err != nil {
		f.ErrorHandler(err, RefoundResponsibility)
		return
	}
}

func (f *FyneApp) SetAdditionalTextPressed() {

}

func (f *FyneApp) printLastCheckPressedFromCRM() {
	click := &entities.Click{}
	err := toml.ReadToml(toml.ClickPath, click)
	if err != nil {
		message := "Ошибка при прочтении истории печати"
		f.service.Errorf("%s: %v", message, err)
		f.ShowWarning(message)
		return
	}

	id := fmt.Sprint(click.Data.OrderId)
	err = nil
	switch click.Data.Type {
	case "order":
		err = f.service.PrintSell(*f.info, id, nil)
	default:
		err = f.service.PrintRefound(*f.info, id, nil)
	}

	if err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) printLastCheckPressedFromKKT() {
	if err := f.service.KKT.PrintLastCheckPressedFromKKT(); err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) exitPressed() {
	f.Logout()
}

func (f *FyneApp) printXReportPressed() {
	if err := f.service.PrintXReport(); err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) WarningPressed() {

}

func (f *FyneApp) AuthorizationPressed(choice bool) { //! обработчик действий
	if choice {
		conf := f.formAppConfig()

		if entities.NewCalculationType(f.authForm.taxesCalculationTypeComboBox.Text) == entities.UndefinedTaxes {
			f.authForm.form.Show()
			f.ShowWarning("Пожалуйста, установите параметр расчета НДС")
			return
		}

		f.Login(conf)
	} else {
		if f.flag.AuthJustHide {
			f.flag.AuthJustHide = false
			return
		}
		f.MainWindow.Close()
	}
}

func (f *FyneApp) SettingWindowPressed(choice bool) {
	settings := f.formDriverData()
	if choice {
		err := toml.WriteToml(toml.DriverInfoPath, settings)
		if err != nil {
			//заполнить
		}

	}
}

func (f *FyneApp) listenerStatusAction() {
	switch f.header.listenerStatus.listenerToolbarItem.Icon {
	case theme.CancelIcon():
		f.header.listenerStatus.listenerToolbarItem.Icon = theme.ConfirmIcon()
		f.flag.StopListen = true
	case theme.ConfirmIcon():
		f.header.listenerStatus.listenerToolbarItem.Icon = theme.CancelIcon()
		f.flag.StopListen = false
	}

	f.header.listenerStatus.listenerToolbar.Refresh()
}

func (f *FyneApp) exitAndCloseShiftButtonPressed() {
	f.LogoutWS()
}

func (f *FyneApp) CloseShift() {
	err := f.service.CloseShift()
	if err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) OpenConnection() {
	err := f.service.KKT.MakeSession(f.info.Session.UserData.FullName, f.info.Session.UserData.Inn)
	if err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) ToolbarInfoPressed() {
	f.AboutDialog.Dialog.Show()

}

func (f *FyneApp) CheckUpdateAction() {

	if _, err := os.Stat("./updater_windows_amd64.exe"); err != nil && os.IsNotExist(err) {
		f.service.Logger.Errorf("Ошибка при запуске центра обновления: %v", err)
		f.ShowWarning("Отсутствует исполняемый файл центра обновлений")
		return
	}

	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", ".\\updater_windows_amd64.exe")
	cmdname, err := os.Executable()

	if err != nil {
		f.service.Logger.Errorf("Ошибка при запуске центра обновления: %v", err)
	}

	updAdrr := strings.Split(f.info.AppConfig.Driver.UpdatePath, "/")
	if len(updAdrr) != 2 {
		f.ShowWarning("В настройках указан неверный источник обновления ПО")
		return
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("version=%s", f.AppInfo.version), fmt.Sprintf("exec_name=%s", cmdname))
	cmd.Env = append(cmd.Env, fmt.Sprintf("pid=%d", os.Getpid()))
	cmd.Env = append(cmd.Env, fmt.Sprintf("repo=%s", updAdrr[1]))
	cmd.Env = append(cmd.Env, fmt.Sprintf("owner=%s", updAdrr[0]))
	f.service.Logger.Infof("Запуск центра обновлений")

	if err := cmd.Start(); err != nil {
		f.service.Logger.Errorf("Ошибка при запуске центра обновления: %v", err)
	}

	if err := cmd.Process.Release(); err != nil {
		f.service.Logger.Errorf("Ошибка при создании свободного процесса: %v", err)
	}

}
