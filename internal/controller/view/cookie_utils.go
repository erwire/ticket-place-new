package view

import (
	"fptr/internal/entities"
	"fptr/pkg/toml"
	"log"
	"time"
)

func (f *FyneApp) formAuthData() entities.UserInfo {
	return entities.UserInfo{
		Login:    f.authForm.loginEntry.Text,
		Password: f.authForm.passwordEntry.Text,
	}
} //собирает данные из программы в структуру

func (f *FyneApp) formDriverData() entities.DriverInfo {
	duration, _ := time.ParseDuration(f.DriverSetting.DriverPollingPeriodSelect.Selected)
	timeoutDuration, err := time.ParseDuration(f.DriverSetting.DriverTimeoutSelect.Selected)
	log.Println(timeoutDuration, err)
	return entities.DriverInfo{
		Path:          f.DriverSetting.DriverPathEntry.Text,
		Com:           f.DriverSetting.DriverComPortEntry.Text,
		Connection:    f.DriverSetting.DriverAddressEntry.Text,
		PollingPeriod: duration,
		TimeoutPeriod: timeoutDuration,
		UpdatePath:    f.DriverSetting.DriverUpdatePath.Text,
	}
} //собирает данные из программы в структуру

func (f *FyneApp) formAppConfig() entities.AppConfig {
	return entities.AppConfig{
		User:   f.formAuthData(),
		Driver: f.formDriverData(),
	}
}

func (f *FyneApp) InitializeCookie() error {
	userInfo, driverInfo := &entities.UserInfo{}, &entities.DriverInfo{}
	session := &entities.SessionInfo{}

	err := toml.ReadToml(toml.DriverInfoPath, driverInfo)
	if err != nil {
		return err
	}

	err = toml.ReadToml(toml.UserInfoPath, userInfo)
	if err != nil {
		return err
	}

	err = toml.ReadToml(toml.SessionPath, session)

	if err != nil {
		return err
	}

	f.info.AppConfig.Driver = *driverInfo
	f.info.AppConfig.User = *userInfo
	f.info.Session = *session

	f.setupCookieIntoEntry()

	return nil
}

func (f *FyneApp) setupDefaultIntoCookie() {
	if f.info.AppConfig.Driver.Connection == "" {
		f.info.AppConfig.Driver.Connection = "https://ticket-place.ru"
		f.service.Logger.Infof("Установлено значение по умолчанию для адреса: %s", f.info.AppConfig.Driver.Connection)
	}

	if f.info.AppConfig.Driver.TimeoutPeriod.Seconds() == 0 {
		f.info.AppConfig.Driver.TimeoutPeriod = time.Second * 20
		f.service.Logger.Infof("Установлено значение по умолчанию для времени жизни запроса: %s", f.info.AppConfig.Driver.TimeoutPeriod)

	}

	if f.info.AppConfig.Driver.PollingPeriod.Seconds() == 0 {
		f.info.AppConfig.Driver.PollingPeriod = time.Second * 5
		f.service.Logger.Infof("Установлено значение по умолчанию %s", f.info.AppConfig.Driver.PollingPeriod)
	}

	if f.info.AppConfig.Driver.UpdatePath == "" {
		f.info.AppConfig.Driver.UpdatePath = "jahngeor/ticket-place"
		f.service.Logger.Infof("Установлено значение по умолчанию %s", f.info.AppConfig.Driver.UpdatePath)
	}

}

func (f *FyneApp) setupCookieIntoEntry() {
	f.setupDefaultIntoCookie()

	f.DriverSetting.DriverPathEntry.Text = f.info.AppConfig.Driver.Path
	f.DriverSetting.DriverAddressEntry.Text = f.info.AppConfig.Driver.Connection
	f.DriverSetting.DriverComPortEntry.Text = f.info.AppConfig.Driver.Com
	f.DriverSetting.DriverTimeoutSelect.Selected = f.info.AppConfig.Driver.TimeoutPeriod.String()
	f.DriverSetting.DriverPollingPeriodSelect.Selected = f.info.AppConfig.Driver.PollingPeriod.String()
	f.DriverSetting.DriverUpdatePath.Text = f.info.AppConfig.Driver.UpdatePath

	f.authForm.loginEntry.Text = f.info.AppConfig.User.Login
	f.authForm.passwordEntry.Text = f.info.AppConfig.User.Password

	f.DriverSetting.DriverPathEntry.Refresh()
	f.DriverSetting.DriverAddressEntry.Refresh()
	f.DriverSetting.DriverComPortEntry.Refresh()
	f.DriverSetting.DriverPollingPeriodSelect.Refresh()
	f.DriverSetting.DriverTimeoutSelect.Refresh()
	f.DriverSetting.DriverUpdatePath.Refresh()
}

func (f *FyneApp) UpdateSession(session entities.SessionInfo) error {
	err := toml.WriteToml(toml.SessionPath, session)
	if err != nil {
		return err
	}
	f.info.Session = session
	return nil
} //заносит данные внутрь структуры, также заносит данные внутрь TOML

func (f *FyneApp) UpdateUserInfo(info entities.UserInfo) error {
	err := toml.WriteToml(toml.UserInfoPath, info)
	if err != nil {
		return err
	}
	f.info.AppConfig.User = info
	return nil
}

func (f *FyneApp) UpdateDriverInfo(info entities.DriverInfo) error {
	err := toml.WriteToml(toml.DriverInfoPath, info)
	if err != nil {
		return err
	}
	f.info.AppConfig.Driver = info
	return nil
}
