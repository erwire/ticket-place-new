package view

import (
	"fmt"
	"fptr/pkg/toml"
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

}

func (f *FyneApp) PrintCheckOnSubmit() {

}

func (f *FyneApp) RefoundOnSubmit() {

}

func (f *FyneApp) SetAdditionalTextPressed() {

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

func (f *FyneApp) WarningPressed() {

}

func (f *FyneApp) AuthorizationPressed(choice bool) { //! обработчик действий
	if choice {
		//f.header.usernameLabel.Text = f.authForm.loginEntry.Text
		//f.header.usernameLabel.Refresh()
		//f.ShowWarning("Ошибка доступа")
	} else {
		f.mainWindow.Close()
	}
}

func (f *FyneApp) SettingWindowPressed(choice bool) {
	settings := f.formDriverData()
	if choice {
		fmt.Println(toml.WriteToml(toml.DriverInfoPath, settings))
	}
}
