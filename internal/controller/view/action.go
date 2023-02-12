package view

import (
	"context"
	"fptr/internal/entities"
	"fptr/pkg/toml"
	"log"
	"time"
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
	f.UpdateSession(entities.SessionInfo{})
	f.header.usernameLabel.Text = ""
	f.authForm.form.Show()
	f.context.cancel()
}

func (f *FyneApp) printXReportPressed() {
	//Механизм напечатания X-отчета
}

func (f *FyneApp) WarningPressed() {

}

func (f *FyneApp) AuthorizationPressed(choice bool) { //! обработчик действий
	if choice {
		appConfig := f.formAppConfig()
		session, message := f.service.Login(appConfig)
		if len(message) != 0 {
			f.authForm.form.Show()
			f.ShowWarning(message)
		} else {
			session.CreatedAt = time.Now()
			err := f.UpdateUserInfo(appConfig.User)
			if err != nil {
				log.Println(err.Error())
			}
			err = f.UpdateDriverInfo(appConfig.Driver)
			if err != nil {
				log.Println(err.Error())
			}
			err = f.UpdateSession(*session)
			if err != nil {
				log.Println(err.Error())
			}
			f.header.usernameLabel.Text = f.info.Session.UserData.Username
			f.header.usernameLabel.Refresh()
			f.context.ctx, f.context.cancel = context.WithCancel(context.Background())
			go f.Listen(f.context.ctx)
		}

	} else {
		f.mainWindow.Close()
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

func (f *FyneApp) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Context Closed")
			return
		default:

		}
	}

}
