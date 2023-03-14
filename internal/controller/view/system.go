package view

func (f *FyneApp) SetAppInfo(version, path, updType string) {
	f.AppInfo = &AppInfo{
		version:    version,
		updatePath: path,
		updateType: updType,
	}
}
