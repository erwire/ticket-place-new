package view

func (f *FyneApp) ShowWarning(err string) {

	f.Error.ErrorText.Text = err
	f.Error.ErrorText.Refresh()
	f.Error.ErrorWindow.Hide()
	f.Error.ErrorWindow.CenterOnScreen()
	f.Error.ErrorWindow.RequestFocus()
	f.Error.ErrorWindow.Show()

	//f.Warning.WarningText.Text = err
	//f.Warning.WarningText.Refresh()
	//f.Warning.WarningWindow.Show()

}
