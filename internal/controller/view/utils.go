package view

import (
	"fptr/internal/entities"
	"fptr/pkg/toml"
	"time"
)

func (f *FyneApp) formAuthData() entities.UserInfo {
	return entities.UserInfo{
		Login:    f.authForm.loginEntry.Text,
		Password: f.authForm.passwordEntry.Text,
	}
}

func (f *FyneApp) formDriverData() entities.DriverInfo {
	duration, _ := time.ParseDuration(f.DriverSetting.DriverPollingPeriodSelect.Selected)
	return entities.DriverInfo{
		Path:          f.DriverSetting.DriverPathEntry.Text,
		Com:           f.DriverSetting.DriverComPortEntry.Text,
		Connection:    f.DriverSetting.DriverAddressEntry.Text,
		PollingPeriod: duration,
	}
}

func (f *FyneApp) SetupCookie() {
	config := &entities.DriverInfo{}
	err := toml.ReadToml(toml.DriverInfoPath, config)

	if err != nil {
		f.ShowWarning(err.Error())
	}

	f.DriverSetting.DriverPathEntry.Text = config.Path
	f.DriverSetting.DriverAddressEntry.Text = config.Connection
	f.DriverSetting.DriverComPortEntry.Text = config.Com
	f.DriverSetting.DriverPollingPeriodSelect.Selected = config.PollingPeriod.String()
}
