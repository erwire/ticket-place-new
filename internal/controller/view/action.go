package view

import (
	"context"
	"fmt"
	"fptr/internal/entities"
	"fptr/pkg/toml"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (f *FyneApp) DriverSettingButtonPressed() {

}

func (f *FyneApp) ShowUnprintedWindow() {
	f.FullfilUnprintedWindow()
	f.UnprintedWindow.Show()
}

func (f *FyneApp) PrintUnprinted() {
	var ecount int = 0
	for _, value := range f.flag.CheckedUnprinted {
		status, err := value.Checked.Get()
		if err != nil {
			ecount++
			continue
		}
		if status {
			err := f.service.PrintSell(*f.info, value.Data, nil, f.flag.pageParams, entities.PrintCheckBox{PrintCheckBox: f.Unprinted.CheckBoxCheck.Checked, PrintOnPrinterTicketBox: f.Unprinted.CheckBoxTicket.Checked})
			if err != nil {
				ecount++
			}
		}

	}

	f.UnprintedWindow.Hide()
	if f.flag.UnprintedContext != nil {
		f.flag.UnprintedCancel()
	}
}

func (f *FyneApp) FullfilUnprintedWindow() {
	f.flag.CheckedUnprinted = make(map[int]*UnprintedValue)

	notes, err := f.service.DS.GetUnfinishedSellsNote(entities.ErrorStatus)
	if err != nil {
		f.ErrorHandler(err, "")
		return
	}

	f.Unprinted.SellsTable.Length = func() (rows int) {
		return len(notes)
	}
	f.Unprinted.SellsTable.CreateItem = func() fyne.CanvasObject {
		var temp bool
		return container.NewGridWithColumns(4, widget.NewCheckWithData("", binding.BindBool(&temp)), widget.NewLabel(""), widget.NewLabel(""), widget.NewLabel(""))
	}
	f.Unprinted.SellsTable.Refresh()
	f.Unprinted.SellsTable.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		if len(notes) > 0 {
			log.Println(len(f.flag.CheckedUnprinted))
			f.flag.CheckedUnprinted[id] = &UnprintedValue{Checked: binding.NewBool(), Data: notes[id].SellID}
			f.flag.CheckedUnprinted[id].Checked.Set(false)
			item.(*fyne.Container).Objects[0].(*widget.Check).Bind(f.flag.CheckedUnprinted[id].Checked)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(notes[id].SellID)
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(notes[id].Date)
			item.(*fyne.Container).Objects[3].(*widget.Label).SetText(notes[id].Event)

		}
	}

}

func (f *FyneApp) DriverPrintHistoryButtonPressed() {
	f.HistoryWindow.CenterOnScreen()
	f.FullfilHistoryWindow()
	f.HistoryWindow.Show()

}

func (f *FyneApp) FullfilHistoryWindow() {
	f.flag.CheckedHistory = make(map[int]*HistoryValue)
	notes, err := f.service.DS.GetAllSellsNote()
	if err != nil {
		f.ErrorHandler(err, "")
		return
	}
	f.History.SellsTable.Length = func() (rows int, cols int) {
		return len(notes), 6
	}

	f.History.SellsTable.UpdateCell = func(id widget.TableCellID, template fyne.CanvasObject) {
		f.flag.CheckedHistory[id.Row] = &HistoryValue{Checked: false, Data: notes[id.Row].SellID}
		switch id.Col {
		case 0:
			//template.(*fyne.Container).Objects[0].Hide()
			template.(*widget.Label).SetText(notes[id.Row].SellID)
		case 1:
			//template.(*fyne.Container).Objects[0].Hide()
			template.(*widget.Label).SetText(notes[id.Row].Date)
		case 2:
			//template.(*fyne.Container).Objects[0].Hide()
			template.(*widget.Label).SetText(notes[id.Row].Event)
		case 3:
			//template.(*fyne.Container).Objects[0].Hide()
			template.(*widget.Label).SetText(strconv.Itoa(int(notes[id.Row].Amount)))
		case 4:
			//template.(*fyne.Container).Objects[0].Hide()
			template.(*widget.Label).SetText(notes[id.Row].Status)
		case 5:
			//template.(*fyne.Container).Objects[0].Hide()
			template.(*widget.Label).SetText(notes[id.Row].Error)
		}
	}

	// запрос к базе

}

func (f *FyneApp) DriverSettingFormOnSubmit() {

}

func (f *FyneApp) DriverApiFormOnSubmit() {

}

func (f *FyneApp) DriverPollingPeriodSelected(selected string) {

}

func (f *FyneApp) CashIncomeOnSubmit() {
	incomeStr := f.PrintSettingsItem.CashIncomeEntry.Text
	income, err := strconv.ParseFloat(incomeStr, 32)
	if err != nil {
		f.ShowWarning("Некорректные данные в поле ввода суммы")
		return
	}
	f.service.Infof("Запрос на внесение: %f руб.", income)

	if err := f.service.CashIncome(income); err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) PrintCheckOnSubmit() {
	id := f.PrintsRefoundAndDeposits.PrintCheckEntry.Text
	if id == "" {
		f.ShowWarning("Пожалуйста, вставьте значение в поле номера чека")
		return
	}
	f.service.Infof("Запрос на печать заказа с номером %s", id)

	if err := f.service.PrintSell(*f.info, id, nil, f.flag.pageParams, f.flag.printCheckBox); err != nil {

		f.ErrorHandler(err, SellResponsibility)
		return
	}
}

func (f *FyneApp) RefoundOnSubmit() {
	id := f.PrintsRefoundAndDeposits.RefoundEntry.Text
	if id == "" {
		f.ShowWarning("Пожалуйста, вставьте значение в поле возврата")
		return
	}
	f.service.Infof("Запрос на печать заказа с номером %s", id)

	if err := f.service.PrintRefoundFromSell(*f.info, id); err != nil {
		f.ErrorHandler(err, RefoundResponsibility)
		return
	}
}

func (f *FyneApp) SetAdditionalTextPressed() {

}

func (f *FyneApp) printLastCheckPressedFromCRM() {
	click := &entities.Click{}
	err := toml.ReadToml(toml.ClickPath, click)
	if err != nil {
		message := "Ошибка при прочтении истории печати"
		f.service.Errorf("%s: %v", message, err)
		f.ShowWarning(message)
		return
	}

	id := fmt.Sprint(click.Data.OrderId)
	err = nil
	switch click.Data.Type {
	case "order":
		err = f.service.PrintSell(*f.info, id, nil, f.flag.pageParams, f.flag.printCheckBox)
	default:
		err = f.service.PrintRefound(*f.info, id, nil)
	}

	if err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) printLastCheckPressedFromKKT() {
	if err := f.service.KKT.PrintLastCheckPressedFromKKT(); err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) exitPressed() {
	f.Logout()
}

func (f *FyneApp) printXReportPressed() {
	if err := f.service.PrintXReport(); err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) WarningPressed() {

}

func (f *FyneApp) AuthorizationPressed(choice bool) { //! обработчик действий
	if choice {
		conf := f.formAppConfig()
		f.Login(conf)

	} else {
		if f.flag.AuthJustHide {
			f.flag.AuthJustHide = false
			return
		}
		f.MainWindow.Close()
	}
}

func (f *FyneApp) SettingWindowPressed(choice bool) {
	settings := f.formDriverData()
	if choice {
		//переписать
		err := toml.WriteToml(toml.DriverInfoPath, settings)
		if err != nil {
			//заполнить
		}
		f.info.AppConfig.Driver.PrinterServiceAddress = f.PrinterSettings.PrinterServiceAddress.Text

	}
}

func (f *FyneApp) listenerStatusAction() {
	switch f.header.listenerStatus.listenerToolbarItem.Icon {
	case theme.CancelIcon():
		f.header.listenerStatus.listenerToolbarItem.Icon = theme.ConfirmIcon()
		f.flag.StopListen = true
	case theme.ConfirmIcon():
		f.header.listenerStatus.listenerToolbarItem.Icon = theme.CancelIcon()
		f.flag.StopListen = false
	}

	f.header.listenerStatus.listenerToolbar.Refresh()
}

func (f *FyneApp) exitAndCloseShiftButtonPressed() {
	data, err := f.service.DS.GetUnfinishedSellsNote(entities.ErrorStatus)
	if err == nil {
		if len(data) > 0 {
			f.flag.UnprintedContext, f.flag.UnprintedCancel = context.WithCancel(context.TODO())
			f.ShowUnprintedWindow()
			select {
			case <-f.flag.UnprintedContext.Done():
				break
			}
		}
	}
	f.LogoutWS()
}

func (f *FyneApp) CloseShift() {
	err := f.service.CloseShift()
	if err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) OpenConnection() {
	err := f.service.KKT.MakeSession(f.info.Session.UserData.FullName)
	if err != nil {
		f.ErrorHandler(err, FunctionResponsibility)
		return
	}
}

func (f *FyneApp) ToolbarInfoPressed() {
	f.AboutDialog.Dialog.Show()

}

func (f *FyneApp) CheckUpdateAction() {

	if _, err := os.Stat("./updater_windows_amd64.exe"); err != nil && os.IsNotExist(err) {
		f.service.Logger.Errorf("Ошибка при запуске центра обновления: %v", err)
		f.ShowWarning("Отсутствует исполняемый файл центра обновлений")
		return
	}

	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", ".\\updater_windows_amd64.exe")
	cmdname, err := os.Executable()

	if err != nil {
		f.service.Logger.Errorf("Ошибка при запуске центра обновления: %v", err)
	}

	updAdrr := strings.Split(f.info.AppConfig.Driver.UpdatePath, "/")
	if len(updAdrr) != 2 {
		f.ShowWarning("В настройках указан неверный источник обновления ПО")
		return
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("version=%s", f.AppInfo.version), fmt.Sprintf("exec_name=%s", cmdname))
	cmd.Env = append(cmd.Env, fmt.Sprintf("pid=%d", os.Getpid()))
	cmd.Env = append(cmd.Env, fmt.Sprintf("repo=%s", updAdrr[1]))
	cmd.Env = append(cmd.Env, fmt.Sprintf("owner=%s", updAdrr[0]))
	f.service.Logger.Infof("Запуск центра обновлений")

	if err := cmd.Start(); err != nil {
		f.service.Logger.Errorf("Ошибка при запуске центра обновления: %v", err)
	}

	if err := cmd.Process.Release(); err != nil {
		f.service.Logger.Errorf("Ошибка при создании свободного процесса: %v", err)
	}

}

func (f *FyneApp) CloseCritical() {
	f.application.Quit()
}

func (f *FyneApp) OpenPrinterSettings() {
	f.PrinterSettingsWindow.Show()
}
