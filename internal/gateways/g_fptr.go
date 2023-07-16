package gateways

import (
	"fmt"
	"fptr/internal/entities"
	apperr "fptr/internal/error_list"
	errorlog "fptr/pkg/error_logs"
	"fptr/pkg/fptr10"
	"fptr/pkg/notes"
	"strconv"
)

var (
	CheckType  = "check"
	TicketType = "ticket"
)

type KKTGateway struct {
	IFptr *fptr10.IFptr
}

func NewKKTGateway(IFptr *fptr10.IFptr) *KKTGateway {
	return &KKTGateway{IFptr: IFptr}
}

func (g *KKTGateway) Configurate() error {
	g.IFptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_PORT, strconv.Itoa(fptr10.LIBFPTR_PORT_USB))
	g.IFptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_AUTO_RECONNECT, "false")
	if err := g.IFptr.ApplySingleSettings(); err != nil {
		return err
	}
	return nil
}

func (g *KKTGateway) Open() error {
	if err := g.IFptr.Open(); err != nil {
		return err
	}
	return nil
}

func (g *KKTGateway) Close() error {
	return nil
}

func (g *KKTGateway) Destroy() {
	g.IFptr.Destroy()
}

func (g *KKTGateway) OpenShift() error {
	return g.IFptr.OpenShift()
}

func (g *KKTGateway) CloseShift() error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_CLOSE_SHIFT)
	return g.IFptr.Report()
}

func (g *KKTGateway) PrintTicket() error {
	return nil
}

/*

	Ticket []: Amount,
	PaymentType



*/

func (g *KKTGateway) PrintSell(sell entities.Sell) error {
	zeroAmountFlag, acceptedFlag, checkStatus := g.ZeroAmountStatus(sell), g.AcceptedForPrint(sell), g.CheckStatus(sell, CheckType)
	if !zeroAmountFlag && acceptedFlag && checkStatus {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL)

		if err := g.IFptr.OpenReceipt(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		var sum int

		for _, ticket := range sell.Data.Tickets {
			if ticket.Amount != 0 && g.TicketStatus(ticket, CheckType) {
				if err := g.PositionRegister(ticket); err != nil {
					g.IFptr.CancelReceipt()
					return err
				}
				sum += ticket.Amount
			}
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, sum)
		if err := g.IFptr.ReceiptTotal(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		switch sell.Data.PaymentType {
		case "cash":
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_CASH)
		default:
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_ELECTRONICALLY)
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, sum)
		if err := g.IFptr.Payment(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		if err := g.IFptr.CloseReceipt(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		for {
			if err := g.IFptr.CheckDocumentClosed(); err != nil {
				continue
			} else {
				break
			}
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_CLOSED) {
			g.IFptr.CancelReceipt()
			return apperr.NewFPTRError(fptr10.LIBFPTR_ERROR_NEED_CANCEL_DOCUMENT, apperr.LibfptrErrorNeedCancelDocument)
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_PRINTED) {
			// Можно сразу вызвать метод допечатывания документа, он завершится с ошибкой, если это невозможно
			for {
				if err := g.IFptr.ContinuePrint(); err != nil {
					continue
				}
			}
		}

		// Запрос информации о закрытом чеке
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_LAST_DOCUMENT)
		if err := g.IFptr.FnQueryData(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		return nil
	} else {
		return g.NotPrintReason(zeroAmountFlag, acceptedFlag, checkStatus)
	}

}
func (g *KKTGateway) PrintRefoundFromCheck(sell entities.Sell) error {
	zeroAmountFlag, acceptedFlag, checkStatus := g.ZeroAmountStatus(sell), g.AcceptedForPrint(sell), g.CheckStatus(sell, CheckType)
	if !zeroAmountFlag && acceptedFlag && checkStatus {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL_RETURN)
		if err := g.IFptr.OpenReceipt(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}
		g.IFptr.SetParam(1212, 4)
		g.IFptr.SetParam(1214, 4)
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "Возврат")
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
		g.IFptr.PrintText()

		var sum int
		for _, ticket := range sell.Data.Tickets {
			if ticket.Amount != 0 && g.TicketStatus(ticket, CheckType) {
				g.PositionRegister(ticket)
				sum += ticket.Amount
			}
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, sum)
		if err := g.IFptr.ReceiptTotal(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		switch sell.Data.PaymentType {
		case "cash":
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_CASH)
		case "card":
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_ELECTRONICALLY)
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, sum)

		if err := g.IFptr.Payment(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		if err := g.IFptr.CloseReceipt(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		for {
			if err := g.IFptr.CheckDocumentClosed(); err != nil {
				continue
			} else {
				break
			}
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_CLOSED) {
			return apperr.NewFPTRError(fptr10.LIBFPTR_ERROR_NEED_CANCEL_DOCUMENT, apperr.LibfptrErrorNeedCancelDocument)
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_PRINTED) {
			for {
				if err := g.IFptr.ContinuePrint(); err != nil {
					continue
				}
			}
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_LAST_DOCUMENT)
		if err := g.IFptr.FnQueryData(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}
		return nil
	} else {
		return g.NotPrintReason(zeroAmountFlag, acceptedFlag, checkStatus)
	}
}
func (g *KKTGateway) PrintRefound(refound entities.Refound) error {
	zeroAmountFlag, acceptedFlag, checkStatus := g.ZeroAmountStatus(refound), g.AcceptedForPrint(refound), g.CheckStatus(refound, CheckType)
	if !zeroAmountFlag && acceptedFlag && checkStatus {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL_RETURN)
		if err := g.IFptr.OpenReceipt(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}
		g.IFptr.SetParam(1212, 4)
		g.IFptr.SetParam(1214, 4)
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "Возврат")
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
		g.IFptr.PrintText()

		var sum int
		for _, ticket := range refound.Data.Tickets {
			if ticket.Amount != 0 && g.TicketStatus(ticket, CheckType) {
				g.PositionRegister(ticket) //добавить обработку ошибки
				sum += ticket.Amount
			}
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, sum)
		if err := g.IFptr.ReceiptTotal(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		if refound.Data.PaymentType != "cash" && refound.Data.Order.PaymentType != "cash" {
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_ELECTRONICALLY)
		} else {
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_CASH)
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, sum)
		err := g.IFptr.Payment()
		if err != nil {
			g.IFptr.CancelReceipt()
			return err
		}
		if err := g.IFptr.CloseReceipt(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		for {
			if err := g.IFptr.CheckDocumentClosed(); err != nil {
				continue
			} else {
				break
			}
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_CLOSED) {
			g.IFptr.CancelReceipt()
			return apperr.NewFPTRError(fptr10.LIBFPTR_ERROR_NEED_CANCEL_DOCUMENT, apperr.LibfptrErrorNeedCancelDocument)
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_PRINTED) {
			for {
				if err := g.IFptr.ContinuePrint(); err != nil {
					continue
				}
			}
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_LAST_DOCUMENT)
		if err := g.IFptr.FnQueryData(); err != nil {
			g.IFptr.CancelReceipt()
			return err
		}

		return nil
	} else {
		return g.NotPrintReason(zeroAmountFlag, acceptedFlag, checkStatus)
	}
}

func (g *KKTGateway) NotPrintReason(zeroAmountFlag, acceptedFlag, checkStatus bool) error {
	reason := ""
	switch true {
	case zeroAmountFlag:
		reason += "заказ с нулевыми позициями"
	case !acceptedFlag:
		if reason != "" {
			reason += ", "
		}
		reason += "заказ с необрабатываемой формой оплаты"
	case !checkStatus:
		if reason != "" {
			reason += ", "
		}
		reason += "в заказе не присутствует позиций с обрабатываемым статусом"
	}

	return apperr.NewBusinessError(fmt.Sprintf("Чек не печатается по причине: %s", reason), errorlog.ValidateError) // ? можно обработать варнинг с нуль-возвратом
}

//# Регистраторы

func (g *KKTGateway) PositionRegister(data entities.TicketData) error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_COMMODITY_NAME, fmt.Sprint(data.Number, ",", data.Event.Show.Name, ",", data.Event.Show.AgeLimit, ",", data.Event.DateTime, ",", data.Zona, ", Ряд:", data.RowSector, ", Место:", data.SeatNumber))
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PRICE, data.Amount)
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_QUANTITY, 1)
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_TAX_TYPE, fptr10.LIBFPTR_TAX_NO)
	g.IFptr.SetParam(1212, 4)
	g.IFptr.SetParam(1214, 4)

	return g.IFptr.Registration()
}
func (g *KKTGateway) NewCashierRegister(fullName string) error {
	g.IFptr.SetParam(1021, fullName)
	g.IFptr.SetParam(1203, "500100732259")
	if err := g.IFptr.OperatorLogin(); err != nil {
		return err
	}
	return nil
}

func (g *KKTGateway) IsOpened() bool {
	return g.IFptr.IsOpened()
}
func (g *KKTGateway) ShiftIsExpired() bool {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_DATA_TYPE, fptr10.LIBFPTR_DT_STATUS)
	g.IFptr.QueryData()
	return g.IFptr.GetParamInt(fptr10.LIBFPTR_PARAM_SHIFT_STATE) == fptr10.LIBFPTR_SS_EXPIRED
}
func (g *KKTGateway) ShiftIsOpened() bool {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_DATA_TYPE, fptr10.LIBFPTR_DT_STATUS)
	g.IFptr.QueryData()
	return g.IFptr.GetParamInt(fptr10.LIBFPTR_PARAM_SHIFT_STATE) == fptr10.LIBFPTR_SS_OPENED
}
func (g *KKTGateway) ShiftIsClosed() bool {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_DATA_TYPE, fptr10.LIBFPTR_DT_STATUS)
	g.IFptr.QueryData()
	return g.IFptr.GetParamInt(fptr10.LIBFPTR_PARAM_SHIFT_STATE) == fptr10.LIBFPTR_SS_CLOSED
}
func (g *KKTGateway) CurrentShiftStatus() uint {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_DATA_TYPE, fptr10.LIBFPTR_DT_STATUS)
	g.IFptr.QueryData()
	return g.IFptr.GetParamInt(fptr10.LIBFPTR_PARAM_SHIFT_STATE)
}

func (g *KKTGateway) PrintXReport() error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_X)
	return g.IFptr.Report()
}
func (g *KKTGateway) CashIncome(income float64) error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, income)
	return g.IFptr.CashIncome()
}

//# Доступ к печати чеков и билетов

func (g *KKTGateway) ZeroAmountStatus(data interface{}) bool {
	ZeroAmountStatus := true
	switch data.(type) {
	case entities.Sell:
		sell := data.(entities.Sell)
		for _, ticket := range sell.Data.Tickets {
			if ticket.Amount != 0 {
				ZeroAmountStatus = false
				break
			}
		}
		return ZeroAmountStatus
	case *entities.Sell:
		sell := data.(*entities.Sell)
		for _, ticket := range sell.Data.Tickets {
			if ticket.Amount != 0 {
				ZeroAmountStatus = false
				break
			}
		}
		return ZeroAmountStatus
	case entities.Refound:
		refound := data.(entities.Refound)
		for _, ticket := range refound.Data.Tickets {
			if ticket.Amount != 0 {
				ZeroAmountStatus = false
				break
			}
		}
		return ZeroAmountStatus
	case *entities.Refound:
		refound := data.(*entities.Refound)
		for _, ticket := range refound.Data.Tickets {
			if ticket.Amount != 0 {
				ZeroAmountStatus = false
				break
			}
		}
		return ZeroAmountStatus
	default:
		return true
	}

	return ZeroAmountStatus
} //# Проверка на нулевые заказы
func (g *KKTGateway) AcceptedForPrint(data interface{}) bool {
	switch data.(type) {
	case entities.Sell:
		sell := data.(entities.Sell)
		return sell.Data.PaymentType == "cash" || sell.Data.PaymentType == "card"
	case *entities.Sell:
		sell := data.(*entities.Sell)
		return sell.Data.PaymentType == "cash" || sell.Data.PaymentType == "card"
	case entities.Refound:
		refound := data.(entities.Refound)
		return refound.Data.Order.PaymentType == "cash" || refound.Data.Order.PaymentType == "card" || refound.Data.PaymentType == "cash" || refound.Data.PaymentType == "card"
	case *entities.Refound:
		refound := data.(*entities.Refound)
		return refound.Data.Order.PaymentType == "cash" || refound.Data.Order.PaymentType == "card" || refound.Data.PaymentType == "cash" || refound.Data.PaymentType == "card"
	default:
		return false
	}
} //# Доступ к печати по оплате
func (g *KKTGateway) TicketStatus(ticket entities.TicketData, ticketType string) bool {
	switch ticketType {
	case CheckType:
		return ticket.Status == "returned" || ticket.Status == "payed"
	case TicketType:
		return ticket.Status == "created" || ticket.Status == "payed"
	default:
		return false
	}
}
func (g *KKTGateway) CheckStatus(data interface{}, typeValue string) bool {
	status := false
	switch data.(type) {
	case entities.Sell:
		sell := data.(entities.Sell)
		for _, ticket := range sell.Data.Tickets {
			if g.TicketStatus(ticket, typeValue) {
				status = true
				break
			}
		}
		return status
	case *entities.Sell:
		sell := data.(*entities.Sell)
		for _, ticket := range sell.Data.Tickets {
			if g.TicketStatus(ticket, typeValue) {
				status = true
				break
			}
		}
		return status
	case entities.Refound:
		refound := data.(entities.Refound)
		for _, ticket := range refound.Data.Tickets {
			if g.TicketStatus(ticket, typeValue) {
				status = true
				break
			}
		}
		return status
	case *entities.Refound:
		refound := data.(*entities.Refound)
		for _, ticket := range refound.Data.Tickets {
			if g.TicketStatus(ticket, typeValue) {
				status = true
				break
			}
		}
		return status
	default:
		return false
	}
}

func (g *KKTGateway) CurrentErrorStatusCode() error {
	status := g.IFptr.ErrorDescription()
	code := g.IFptr.ErrorCode()
	if status != "" || status != apperr.LibfptrErrLibfptrStatusOk {
		return apperr.NewFPTRError(code, status)
	}
	return nil
}

func (g *KKTGateway) PrintLastCheckPressedFromKKT() error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_LAST_DOCUMENT)
	return g.IFptr.Report()
}

func (g *KKTGateway) WarningBeep() {

	var notes = [...]int{
		notes.NoteC4, notes.NoteG4,
	}

	var times = [...]int{200, 200}

	for key := range notes {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FREQUENCY, notes[key])
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_DURATION, times[key])
		g.IFptr.Beep()
	}

}

func (g *KKTGateway) ErrorBeep() {

	var notes = [...]int{
		311, 466, 392,
	}

	var times = [...]int{
		400, 100, 750,
	}

	for key := range notes {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FREQUENCY, notes[key])
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_DURATION, times[key])
		g.IFptr.Beep()
	}

}
