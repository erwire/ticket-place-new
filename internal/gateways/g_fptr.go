package gateways

import (
	"fmt"
	"fptr/internal/entities"
	errorlog "fptr/pkg/error_logs"
	"fptr/pkg/fptr10"
	"log"
	"strconv"
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
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_CLOSE_SHIFT)
	return g.IFptr.Report()
}

func (g *KKTGateway) PrintTicket() error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}
	return nil
}

func (g *KKTGateway) PrintSell(sell entities.Sell) error {
	if !g.IsOpened() {
		return errorlog.BoxOfficeIsNotOpenError
	}
	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL)

	if err := g.IFptr.OpenReceipt(); err != nil {
		if g.IFptr.ErrorCode() == fptr10.LIBFPTR_ERROR_DENIED_IN_OPENED_RECEIPT {
			g.IFptr.CancelReceipt()
		}

		return fmt.Errorf("%w: %s", errorlog.OpenReceiptError, err.Error())
	}

	var sum int

	for _, ticket := range sell.Data.Tickets {
		if ticket.Amount != 0 && ticket.Status == "payed" {
			if err := g.PositionRegister(ticket); err != nil {
				return err
			}
			sum += ticket.Amount
		}
	}

	g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, sum)
	if err := g.IFptr.ReceiptTotal(); err != nil {
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
			g.CloseShift() //-
			g.OpenShift()  //-
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

func (g *KKTGateway) PrintRefound(refound entities.Refound) error {
	ZeroAmountStatus := true
	for _, ticket := range refound.Data.Tickets {
		if ticket.Amount != 0 {
			ZeroAmountStatus = false
			break
		}
	}

	if !ZeroAmountStatus {
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL_RETURN)
		if err := g.IFptr.OpenReceipt(); err != nil {
			if g.IFptr.ErrorCode() == fptr10.LIBFPTR_ERROR_DENIED_IN_OPENED_RECEIPT {
				g.IFptr.CancelReceipt()
				g.IFptr.OpenReceipt()
			} else {
				return errorlog.OpenReceiptError
			}
		}
		g.IFptr.SetParam(1212, 4)
		g.IFptr.SetParam(1214, 4)
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "Возврат")
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
		g.IFptr.PrintText()

		var sum int
		for _, value := range refound.Data.Tickets {
			if value.Amount != 0 {
				g.PositionRegister(value)
				sum += value.Amount
			}
		}
		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, sum)
		if err := g.IFptr.ReceiptTotal(); err != nil { //!

		}
		if refound.Data.PaymentType != "cash" && refound.Data.Order.PaymentType != "cash" {
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_ELECTRONICALLY)
		} else {
			g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_CASH)
		}

		g.IFptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, sum)
		g.IFptr.Payment() //! err

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
