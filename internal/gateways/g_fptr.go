package gateways

import (
	"errors"
	"fmt"
	"fptr/internal/entities"
	errorlog "fptr/pkg/error_logs"
	"fptr/pkg/fptr10"
	"log"
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

func (g *KKTGateway) OpenShift() error {
	return g.IFptr.OpenShift()
}

func (g *KKTGateway) CloseShift() error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_CLOSE_SHIFT)
	return g.IFptr.Report()
}

func (g *KKTGateway) PrintTicket() error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}
	return nil
}

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

func (g *KKTGateway) PrintSell(sell entities.Sell) error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}

	if g.ShiftIsExpired() {
		return errorlog.ShiftIsExpired // ! не обязательна обработка, если смена истекла - выйдет ошибка FPTR.ERROR
	}
	if !g.ZeroAmountStatus(sell) && g.AcceptedForPrint(sell) {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL)

		if err := g.IFptr.OpenReceipt(); err != nil {
			if g.IFptr.ErrorCode() == fptr10.LIBFPTR_ERROR_DENIED_IN_OPENED_RECEIPT {
				err = g.IFptr.CancelReceipt()
				if err != nil {
					return fmt.Errorf("%w: %s", errorlog.CantCancelReceipt, err.Error())
				}
				err = g.IFptr.OpenReceipt()
				if err != nil {
					return fmt.Errorf("%w: %s", errorlog.OpenReceiptError, err.Error())
				}
			}

			return fmt.Errorf("%w: %s", errorlog.OpenReceiptError, err.Error())
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
			return err
		}

		if err := g.IFptr.CloseReceipt(); err != nil {
			if g.ShiftIsExpired() {
				g.CloseShift()
				g.OpenShift()
			}
		}

		for {
			if err := g.IFptr.CheckDocumentClosed(); err != nil {
				log.Println(g.IFptr.ErrorDescription()) //позднее добавить механизм по таймеру для избежания циклов
				continue
			} else {
				break
			}
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_CLOSED) {
			if err := g.IFptr.CancelReceipt(); err != nil {
				return errorlog.CantCancelReceipt
			}
			return errorlog.DocumentNotClosed
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_PRINTED) {
			// Можно сразу вызвать метод допечатывания документа, он завершится с ошибкой, если это невозможно
			for {
				if err := g.IFptr.ContinuePrint(); err != nil {
					log.Printf("Не удалось напечатать документ (Ошибка \"%v\"). Устраните неполадку и повторите.", g.IFptr.ErrorDescription()) //исправить использование бесконечного цикла
					continue
				}
			}
		}

		// Запрос информации о закрытом чеке
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_LAST_DOCUMENT)
		if err := g.IFptr.FnQueryData(); err != nil {
			return err
		}

		return nil
	}
	return nil // обработка нулевых заказов
}

func (g *KKTGateway) PrintXReport() error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_X)
	return g.IFptr.Report()
}

func (g *KKTGateway) CashIncome(income float64) error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, income)
	return g.IFptr.CashIncome()
}

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
}

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
}

func (g *KKTGateway) PrintRefoundFromCheck(sell entities.Sell) error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}

	if g.ShiftIsExpired() {
		return errorlog.ShiftIsExpired
	}

	if !g.ZeroAmountStatus(sell) && g.AcceptedForPrint(sell) {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL_RETURN)
		if err := g.IFptr.OpenReceipt(); err != nil {
			if g.IFptr.ErrorCode() == fptr10.LIBFPTR_ERROR_DENIED_IN_OPENED_RECEIPT {
				err = g.IFptr.CancelReceipt()
				if err != nil {
					return fmt.Errorf("%w: %s", errorlog.CantCancelReceipt, err.Error())
				}
				err = g.IFptr.OpenReceipt()
				if err != nil {
					return fmt.Errorf("%w: %s", errorlog.OpenReceiptError, err.Error())
				}
			}

			return fmt.Errorf("%w: %s", errorlog.OpenReceiptError, err.Error())
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
		if err := g.IFptr.ReceiptTotal(); err != nil { //!
			err = g.IFptr.CancelReceipt()
			if err != nil {
				return fmt.Errorf("%w: %s", errorlog.CantCancelReceipt, err.Error())
			}
			return fmt.Errorf("чек закрыт из-за проблем с расчетом итога: %w", err)
		}

		switch sell.Data.PaymentType {
		case "cash":
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_CASH)
		case "card":
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_ELECTRONICALLY)
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, sum)
		err := g.IFptr.Payment() //! err
		if err != nil {
			err = g.IFptr.CancelReceipt()
			if err != nil {
				return fmt.Errorf("чек закрыт из-за проблем с расчетом полученной суммы: %w", err)
			}
		}
		if err := g.IFptr.CloseReceipt(); err != nil {
			if g.ShiftIsExpired() {
				g.CloseShift()
				g.OpenShift()
			}
		}

		for {
			if err := g.IFptr.CheckDocumentClosed(); err != nil {
				log.Println(g.IFptr.ErrorDescription()) //позднее добавить механизм по таймеру для избежания циклов
				continue
			} else {
				break
			}
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_CLOSED) {
			if err := g.IFptr.CancelReceipt(); err != nil {
				return errorlog.CantCancelReceipt
			}
			return errorlog.DocumentNotClosed
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_PRINTED) {
			// Можно сразу вызвать метод допечатывания документа, он завершится с ошибкой, если это невозможно
			for {
				if err := g.IFptr.ContinuePrint(); err != nil {
					log.Printf("Не удалось напечатать документ (Ошибка \"%v\"). Устраните неполадку и повторите.", g.IFptr.ErrorDescription()) //исправить использование бесконечного цикла
					continue
				}
			}
		}

		// Запрос информации о закрытом чеке
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_LAST_DOCUMENT)
		if err := g.IFptr.FnQueryData(); err != nil {
			return err
		}

		return nil
	} else {
		return nil // ? можно обработать варнинг с нуль-возвратом
	}
}

func (g *KKTGateway) PrintRefound(refound entities.Refound) error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}

	if g.ShiftIsExpired() {
		return errorlog.ShiftIsExpired
	}

	if !g.ZeroAmountStatus(refound) && g.AcceptedForPrint(refound) {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL_RETURN)
		if err := g.IFptr.OpenReceipt(); err != nil {
			if g.IFptr.ErrorCode() == fptr10.LIBFPTR_ERROR_DENIED_IN_OPENED_RECEIPT {
				err = g.IFptr.CancelReceipt()
				if err != nil {
					return fmt.Errorf("%w: %s", errorlog.CantCancelReceipt, err.Error())
				}
				err = g.IFptr.OpenReceipt()
				if err != nil {
					return fmt.Errorf("%w: %s", errorlog.OpenReceiptError, err.Error())
				}
			}

			return fmt.Errorf("%w: %s", errorlog.OpenReceiptError, err.Error())
		}
		g.IFptr.SetParam(1212, 4)
		g.IFptr.SetParam(1214, 4)
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "Возврат")
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
		g.IFptr.PrintText()

		var sum int
		for _, ticket := range refound.Data.Tickets {
			if ticket.Amount != 0 && g.TicketStatus(ticket, CheckType) {
				g.PositionRegister(ticket)
				sum += ticket.Amount
			}
		}
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, sum)
		if err := g.IFptr.ReceiptTotal(); err != nil { //!
			err = g.IFptr.CancelReceipt()
			if err != nil {
				return fmt.Errorf("%w: %s", errorlog.CantCancelReceipt, err.Error())
			}
			return fmt.Errorf("чек закрыт из-за проблем с расчетом итога: %w", err)
		}
		if refound.Data.PaymentType != "cash" && refound.Data.Order.PaymentType != "cash" {
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_ELECTRONICALLY)
		} else {
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_CASH)
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, sum)
		err := g.IFptr.Payment() //! err
		if err != nil {
			err = g.IFptr.CancelReceipt()
			if err != nil {
				return fmt.Errorf("чек закрыт из-за проблем с расчетом полученной суммы: %w", err)
			}
		}
		if err := g.IFptr.CloseReceipt(); err != nil {
			if g.ShiftIsExpired() {
				g.CloseShift()
				g.OpenShift()
			}
		}

		for {
			if err := g.IFptr.CheckDocumentClosed(); err != nil {
				log.Println(g.IFptr.ErrorDescription()) //позднее добавить механизм по таймеру для избежания циклов
				continue
			} else {
				break
			}
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_CLOSED) {
			if err := g.IFptr.CancelReceipt(); err != nil {
				return errorlog.CantCancelReceipt
			}
			return errorlog.DocumentNotClosed
		}

		if !g.IFptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_PRINTED) {
			// Можно сразу вызвать метод допечатывания документа, он завершится с ошибкой, если это невозможно
			for {
				if err := g.IFptr.ContinuePrint(); err != nil {
					log.Printf("Не удалось напечатать документ (Ошибка \"%v\"). Устраните неполадку и повторите.", g.IFptr.ErrorDescription()) //исправить использование бесконечного цикла
					continue
				}
			}
		}

		// Запрос информации о закрытом чеке
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_LAST_DOCUMENT)
		if err := g.IFptr.FnQueryData(); err != nil {
			return err
		}

		return nil
	} else {
		return nil // ? можно обработать варнинг с нуль-возвратом
	}
}

func (g *KKTGateway) PositionRegister(data entities.TicketData) error {
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_COMMODITY_NAME, fmt.Sprint(data.Number, ",", data.Event.Show.Name, ",", data.Event.Show.AgeLimit, ",", data.Event.DateTime, ",", data.Zona, ", Ряд:", data.RowSector, ", Место:", data.SeatNumber))
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PRICE, data.Amount)
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_QUANTITY, 1)
	g.IFptr.SetParam(1212, 4)
	g.IFptr.SetParam(1214, 4)
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_TAX_TYPE, fptr10.LIBFPTR_TAX_NO)
	return g.IFptr.Registration()
}

func (g *KKTGateway) IsOpened() bool {
	return g.IFptr.IsOpened()
}

func (g *KKTGateway) NewCashierRegister(info entities.SessionInfo) error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}

	g.IFptr.SetParam(1021, info.UserData.FullName)
	g.IFptr.SetParam(1203, "500100732259")
	if err := g.IFptr.OperatorLogin(); err != nil {
		return err
	}
	return nil
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

func (g *KKTGateway) CurrentErrorStatusCode() error {
	status := g.IFptr.ErrorDescription()
	if status != "" {
		return errors.New(status)
	}
	return nil
}
