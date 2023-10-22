package view

import (
	"fptr/internal/entities"
	"fptr/pkg/toml"
	"time"
)

// Собирает данные из формы авторизации в структуру
func (f *FyneApp) formAuthData() entities.UserInfo {
	return entities.UserInfo{
		Login:    f.authForm.loginEntry.Text,
		Password: f.authForm.passwordEntry.Text,
	}
}

// Собирает данные из формы настроек драйвера в структуру
func (f *FyneApp) formDriverData() entities.DriverInfo {
	duration, _ := time.ParseDuration(f.DriverSetting.DriverPollingPeriodSelect.Selected)
	timeoutDuration, err := time.ParseDuration(f.DriverSetting.DriverTimeoutSelect.Selected)
	if err != nil {
		f.service.Logger.Warningf("Ошибка в процессе сбора данных из формы в структуру данных", err)
	}
	return entities.DriverInfo{
		Path:                  f.DriverSetting.DriverPathEntry.Text,
		Com:                   f.DriverSetting.DriverComPortEntry.Text,
		Connection:            f.DriverSetting.DriverAddressEntry.Text,
		PollingPeriod:         duration,
		TimeoutPeriod:         timeoutDuration,
		UpdatePath:            f.DriverSetting.DriverUpdatePath.Text,
		PrinterServiceAddress: f.PrinterSettings.PrinterServiceAddress.Text,
		PrinterName:           f.PrinterSettings.SelectPrinter.Selected,
	}
}

// Собирает данные из форм авторизации и настроек в структуру
func (f *FyneApp) formAppConfig() entities.AppConfig {
	return entities.AppConfig{
		User:   f.formAuthData(),
		Driver: f.formDriverData(),
	}
}

// Считывает данные из файлов хранения в кэш
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

// Устанавливает значения по умолчанию в структуру кэша
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

	if f.info.AppConfig.Driver.PrinterServiceAddress == "" {
		f.info.AppConfig.Driver.PrinterServiceAddress = "1000"
		f.service.Logger.Infof("Установлено значение по умолчанию %s", f.info.AppConfig.Driver.PrinterServiceAddress)
	}

}

// Устанавливает значения настроек в поля формы из структуры кэша
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
	f.setupPrinterSettingsIntoEntry()
}

func (f *FyneApp) setupPrinterSettingsIntoEntry() {
	f.PrinterSettings.SelectPrinter.Selected = f.info.AppConfig.Driver.PrinterName
	f.PrinterSettings.PrinterServiceAddress.Text = f.info.AppConfig.Driver.PrinterServiceAddress

	f.PrinterSettings.SelectPrinter.Refresh()
	f.PrinterSettings.PrinterServiceAddress.Refresh()
}

// Обновляет значения сессии в структуре кэша и в файлах хранения
func (f *FyneApp) UpdateSession(session entities.SessionInfo) error {
	err := toml.WriteToml(toml.SessionPath, session)
	if err != nil {
		return err
	}
	f.info.Session = session
	return nil
}

// Обновляет значения данных пользователя в структуре кэша и в файлах хранения
func (f *FyneApp) UpdateUserInfo(info entities.UserInfo) error {
	err := toml.WriteToml(toml.UserInfoPath, info)
	if err != nil {
		return err
	}
	f.info.AppConfig.User = info
	return nil
}

// Обновляет значения данных настроек в структуре кэша и в файлах хранения
func (f *FyneApp) UpdateDriverInfo(info entities.DriverInfo) error {
	err := toml.WriteToml(toml.DriverInfoPath, info)
	if err != nil {
		return err
	}
	f.info.AppConfig.Driver = info
	return nil
}
