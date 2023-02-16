package view

import (
	"context"
	"fmt"
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
			userPrev := &entities.UserInfo{}
			if err := toml.ReadToml(toml.UserInfoPath, userPrev); err != nil {
				log.Println(err)
				return
			}
			if userPrev.ValidateUser() {
				if userPrev.Login != appConfig.User.Login {
					if f.service.ShiftIsOpened() || f.service.ShiftIsExpired() {
						if message := f.service.CloseShift(); message != "" {
							f.service.LoggerService.Errorf(message)
							f.ShowWarning("Проблемы соединения с ККТ")
							f.authForm.form.Show()
							return
						}
					}
				}
			}

			session.CreatedAt = time.Now()
			err := f.UpdateUserInfo(appConfig.User)
			if err != nil {
				f.service.Errorf("%v\n", err.Error())
			}
			err = f.UpdateDriverInfo(appConfig.Driver)
			if err != nil {
				f.service.Errorf("%v\n", err.Error())
			}
			err = f.UpdateSession(*session)
			if err != nil {
				f.service.Errorf("%v\n", err.Error())
				return
			}
			f.header.usernameLabel.Text = f.info.Session.UserData.Username
			f.header.usernameLabel.Refresh()

			click, message := f.service.GetLastReceipt(f.info.AppConfig.Driver.Connection, f.info.Session)
			if message != "" {
				f.ShowWarning("Внимание, возможно задвоение чека!")
			}
			err = toml.WriteToml(toml.ClickPath, click)
			if err != nil {
				f.authForm.form.Show()
				f.service.Errorf("%v\n", err.Error())
				return
			}

			message = f.service.MakeSession(*f.info)
			if message != "" {
				f.authForm.form.Show()
				f.ShowWarning(message)
				return
			}

			f.context.ctx, f.context.cancel = context.WithCancel(context.Background())
			f.service.Infof("Успешная авторизация под пользователем %s\n", f.info.Session.UserData.FullName)
			f.service.Infoln("Горутина прослушивания запущена")
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

func (f *FyneApp) exitAndCloseShiftButtonPressed() {
	message := f.service.KKT.CloseShift()
	if message != "" {
		f.ShowWarning(message)
	}
	f.exitPressed()
}

func (f *FyneApp) Listen(ctx context.Context, info entities.Info) {
	for {
		select {
		case <-ctx.Done():
			f.service.Info("Горутина прослушивания закрыта")
			return
		default:
			if f.info.Session.IsDead() {
				f.UpdateSession(entities.SessionInfo{})
				f.ShowWarning("Ваша сессия устарела. Пожалуйста, пройдите повторную авторизацию.")
				f.context.cancel()
			}
			if f.flag.StopListen {
				continue
			}
			time.Sleep(info.AppConfig.Driver.PollingPeriod * time.Nanosecond)
			clickCache := &entities.Click{}

			if err := toml.ReadToml(toml.ClickPath, clickCache); err != nil {
				f.service.Infoln(err)
				continue
			}

			click, message := f.service.GetLastReceipt(info.AppConfig.Driver.Connection, info.Session)

			if message != "" {
				//обработка ошибки (аля потеряно соединение с сервером)
			}

			if clickCache.Data.Id == click.Data.Id {
				continue
			}

			if err := toml.WriteToml(toml.ClickPath, click); err != nil {
				f.service.Infoln(err)
				continue
			}
			if f.flag.DebugOn {
				// отладка
			} else {
				if f.flag.PrintCheckBox {
					switch click.Data.Type {
					case "order":
						message = f.service.PrintSell(*f.info, fmt.Sprint(click.Data.OrderId))
						if message != "" {
							f.ShowWarning(message)
							continue
						}
					default:
						message = f.service.PrintRefound(*f.info, fmt.Sprint(click.Data.OrderId))
						if message != "" {
							f.ShowWarning(message)
							continue
						}
					}

				}

				if f.flag.PrintOnPrinterTicketBox {

				}

				if f.flag.PrintOnKKTTicketCheckBox {

				}
			}

		}
	}

}
