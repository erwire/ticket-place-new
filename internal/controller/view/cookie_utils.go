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

func (f *FyneApp) setupCookieIntoEntry() {
	f.DriverSetting.DriverPathEntry.Text = f.info.AppConfig.Driver.Path
	f.DriverSetting.DriverAddressEntry.Text = f.info.AppConfig.Driver.Connection
	f.DriverSetting.DriverComPortEntry.Text = f.info.AppConfig.Driver.Com

	if f.info.AppConfig.Driver.PollingPeriod.String() == "0s" {
		f.DriverSetting.DriverPollingPeriodSelect.Selected = "2s"
	} else {
		f.DriverSetting.DriverTimeoutSelect.Selected = f.info.AppConfig.Driver.TimeoutPeriod.String()
	}

	if f.info.AppConfig.Driver.TimeoutPeriod.String() == "0s" {
		f.DriverSetting.DriverTimeoutSelect.Selected = "20s"
	} else {
		f.DriverSetting.DriverPollingPeriodSelect.Selected = f.info.AppConfig.Driver.PollingPeriod.String()
	}

	f.authForm.loginEntry.Text = f.info.AppConfig.User.Login
	f.authForm.passwordEntry.Text = f.info.AppConfig.User.Password
	f.DriverSetting.DriverPathEntry.Refresh()
	f.DriverSetting.DriverAddressEntry.Refresh()
	f.DriverSetting.DriverComPortEntry.Refresh()
	f.DriverSetting.DriverPollingPeriodSelect.Refresh()

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
