package view

func (f *FyneApp) SetAppInfo(version, path, updType, updateRepo string) {
	f.AppInfo = &AppInfo{
		version:    version,
		updatePath: path,
		updateType: updType,
		updateRepo: updateRepo,
	}
}
