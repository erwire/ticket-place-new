package view

import (
	"context"
	"fptr/internal/entities"
	"fptr/pkg/toml"
	"fyne.io/fyne/v2/theme"
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

			click, message := f.service.GetLastReceipt(f.info.AppConfig.Driver.Connection, f.info.Session)
			if message != "" {
				f.ShowWarning("Внимание, возможно задвоение чека!")
				log.Println(err.Error()) //обработку сделать нормальную
			}
			err = toml.WriteToml(toml.ClickPath, click)
			if err != nil {
				log.Println(err.Error())
			}

			f.context.ctx, f.context.cancel = context.WithCancel(context.Background())
			go f.Listen(f.context.ctx, *f.info)
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

func (f *FyneApp) Listen(ctx context.Context, info entities.Info) {

	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			log.Println("Context Closed")
			return
		default:
			if f.flag.StopListen {
				continue
			}
			time.Sleep(info.AppConfig.Driver.PollingPeriod * time.Nanosecond)
			clickCache := &entities.Click{}
			err := toml.ReadToml(toml.ClickPath, clickCache)

			if err != nil {
				log.Println(err.Error())
				continue
			}

			click, message := f.service.GetLastReceipt(info.AppConfig.Driver.Connection, info.Session)

			if message != "" {
				//обработка ошибки (аля потеряно соединение с сервером)
			}

			if clickCache.Data.Id == click.Data.Id {
				log.Println("Билеты одинаковы!")
				//continue
			}

			if f.flag.DebugOn {
				logMes := ""

				if f.flag.PrintOnPrinterTicketBox {
					logMes += "Печать билета на принтере/"
				}

				if f.flag.PrintOnKKTTicketCheckBox {
					logMes += "Печать билета на ККТ/"
				}

				if f.flag.PrintCheckBox {
					logMes += "Печать чека/"
				}

				log.Println(logMes)

				toml.WriteToml(toml.ClickPath, click)
			}

		}
	}

}
