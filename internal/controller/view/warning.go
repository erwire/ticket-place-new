package view

func (f *FyneApp) ShowWarning(err string) {
	f.Warning.WarningText.Text = err
	f.Warning.WarningText.Refresh()
	f.Warning.WarningWindow.Show()
}
