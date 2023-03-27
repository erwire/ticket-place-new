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

func (f *FyneApp) ShowPrintConfirm() {
	f.PrintDoubleConfirm.Window.Hide()
	f.PrintDoubleConfirm.Window.CenterOnScreen()
	f.PrintDoubleConfirm.Window.RequestFocus()
	f.PrintDoubleConfirm.Window.Show()
}
