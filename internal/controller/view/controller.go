package view

import (
	"context"
	"fmt"
	"fptr/internal/entities"
	"fptr/pkg/fptr10"
	"fptr/pkg/toml"
	"time"
)

func (f *FyneApp) GetClickAndWriteIntoToml() {
	click, err := f.service.GetLastReceipt(f.info.AppConfig.Driver.Connection, f.info.Session)
	if err != nil {
		f.ErrorHandler(err, LoginResponsibility)
		return
	}
	if err := toml.WriteToml(toml.ClickPath, click); err != nil {
		f.ErrorHandler(err, LoginResponsibility)
		return
	}
}

func (f *FyneApp) IsNewUser(newUser entities.UserInfo) bool {
	oldUser := &entities.UserInfo{}
	err := toml.ReadToml(toml.UserInfoPath, oldUser)
	if err != nil {
		f.ErrorHandler(err, LoginResponsibility)
		return true
	}
	if !oldUser.ValidateUser() {
		return true
	}

	return oldUser.Login != newUser.Login
}

func (f *FyneApp) Login(conf entities.AppConfig) {
	session, err := f.service.Login(conf)
	if err != nil {
		f.authForm.form.Show()
		f.ErrorHandler(err, LoginResponsibility)
		return
	}
	if f.IsNewUser(conf.User) {
		switch f.service.CurrentShiftStatus() {
		case fptr10.LIBFPTR_SS_CLOSED:
			break
		case fptr10.LIBFPTR_SS_OPENED, fptr10.LIBFPTR_SS_EXPIRED:
			if err := f.service.CloseShift(); err != nil {
				f.authForm.form.Show()
				f.ErrorHandler(err, LoginResponsibility)
			}
		default:
			break
		}
	}

	session.CreatedAt = time.Now()

	if err := f.UpdateUserInfo(conf.User); err != nil {
		f.authForm.form.Show()
		f.ErrorHandler(err, LoginResponsibility)
		return
	}
	if err := f.UpdateDriverInfo(conf.Driver); err != nil {
		f.authForm.form.Show()
		f.ErrorHandler(err, LoginResponsibility)
		return
	}
	if err := f.UpdateSession(*session); err != nil {
		f.authForm.form.Show()
		f.ErrorHandler(err, LoginResponsibility)
		return
	}

	if err := f.service.MakeSession(*f.info); err != nil {
		f.authForm.form.Show()
		f.ErrorHandler(err, LoginResponsibility)
		return
	}

	f.header.usernameLabel.Text = f.info.Session.UserData.Username
	f.header.usernameLabel.Refresh()
	f.GetClickAndWriteIntoToml()

	f.context.ctx, f.context.cancel = context.WithCancel(context.Background())
	f.StartListen()

} //# вход в смену

func (f *FyneApp) Logout() {
	f.UpdateSession(entities.SessionInfo{})
	f.header.usernameLabel.Text = ""
	f.authForm.form.Show()
	f.StopListen()
} //# выход из сессии

func (f *FyneApp) LogoutWS() {
	err := f.service.CloseShift()
	if err != nil {
		f.service.Errorf("Ошибка закрытия смены при выходе из сессии: %v", err)
	}

	f.StopListen()

	f.UpdateSession(entities.SessionInfo{})
	f.header.usernameLabel.Text = ""
	f.authForm.form.Show()
	f.service.Infof("Произошел успешный выход из сессии с закрытием смены")
} //# выход из сессии с закрытием смены

func (f *FyneApp) StartListen() {
	f.service.Infof("Запущен поток прослушивания по адресу: %s", f.info.AppConfig.Driver.Connection)
	go f.Listen(f.context.ctx, *f.info)
} //# Начать прослушку

func (f *FyneApp) StopListen() {
	if f.context.ctx != nil {
		select {
		case <-f.context.ctx.Done():
			f.service.Warningf("Попытка остановить закрытый поток")
			return
		default:
			f.context.cancel()
			f.service.Infof("Остановлен поток прослушивания по адресу: %s", f.info.AppConfig.Driver.Connection)
			return
		}
	}

	f.service.Warningf("Попытка закрыть поток, который еще не создан!\n")
} //# Закрыть прослушку

func (f *FyneApp) WarningWSShow() {

} //# Критическая ошибка смены

func (f *FyneApp) WarningShow() {

} //# Ошибка

func (f *FyneApp) Listen(ctx context.Context, info entities.Info) {
	for {

		select {
		case <-ctx.Done():
			return
		default:
			if f.flag.StopListen {
				continue
			}

			if f.info.Session.IsDead() {
				f.Logout()
				f.ShowWarning("Ваша сессия устарела. Пожалуйста, авторизуйтесь снова!")
			}

			time.Sleep(info.AppConfig.Driver.PollingPeriod * time.Nanosecond)

			clickCache := &entities.Click{}

			if err := toml.ReadToml(toml.ClickPath, clickCache); err != nil {
				f.service.Infoln(err)
				continue
			}

			click, err := f.service.GetLastReceipt(info.AppConfig.Driver.Connection, info.Session)

			if err != nil {
				f.ErrorHandler(err, ClickResponsibility)
				continue
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
						if err = f.service.PrintSell(*f.info, fmt.Sprint(click.Data.OrderId)); err != nil {
							f.ErrorHandler(err, SellResponsibility)
							continue
						}
					default:
						if err = f.service.PrintRefound(*f.info, fmt.Sprint(click.Data.OrderId)); err != nil {
							f.ErrorHandler(err, RefoundResponsibility)
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
