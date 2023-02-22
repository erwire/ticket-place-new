package view

import "fptr/internal/entities"

func (f *FyneApp) Login() {

} //# вход в смену

func (f *FyneApp) Logout() {
	f.UpdateSession(entities.SessionInfo{})
	f.header.usernameLabel.Text = ""
	f.authForm.form.Show()
	f.context.cancel()
} //# выход из сессии

func (f *FyneApp) LoginWS() {

} //# авторизация со сменой

func (f *FyneApp) LogoutWS() {
	err := f.service.CloseShift()
	if err != nil {
		f.service.Errorf("Ошибка закрытия смены при выходе из сессии: %v", err)
	}
	f.UpdateSession(entities.SessionInfo{})
	f.header.usernameLabel.Text = ""
	f.authForm.form.Show()
	f.context.cancel()
} //# выход из сессии с закрытием смены

func (f *FyneApp) StartListen() {

} //# Начать прослушку

func (f *FyneApp) StopListen() {
	f.context.cancel()
} //# Закрыть прослушку

func (f *FyneApp) WarningWSShow() {

} //# Критическая ошибка смены

func (f *FyneApp) WarningShow() {

} //# Ошибка
