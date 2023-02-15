package main

import (
	"fptr/pkg/fptr10"
	"log"
	"strconv"
)

func main() {
	fptr := fptr10.New()
	defer fptr.Destroy()

	fptr.SetSingleSetting(fptr10.LIBFPTR_SETTING_PORT, strconv.Itoa(fptr10.LIBFPTR_PORT_USB))
	if err := fptr.ApplySingleSettings(); err != nil {
		log.Println(err)
		return
	}

	// Соединение с ККТ
	if err := fptr.Open(); err != nil {
		log.Println(err)
		return
	}

	// Регистрация кассира
	fptr.SetParam(1021, "Иванов И.И.")
	fptr.SetParam(1203, "500100732259")
	if err := fptr.OperatorLogin(); err != nil {
		log.Println(err)
		return
	}

	// Открытие чека (с передачей телефона получателя)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_RECEIPT_TYPE, fptr10.LIBFPTR_RT_SELL)
	fptr.SetParam(1008, "+79161234567")
	if err := fptr.OpenReceipt(); err != nil {
		log.Println(err)
		return
	}

	fptr.SetParam(1262, "020")
	fptr.SetParam(1263, "14.12.2018")
	fptr.SetParam(1264, "1556")
	fptr.SetParam(1265, "tm=mdlp&sid=00000000105200")
	fptr.UtilFormTlv()

	industryInfo := fptr.GetParamByteArray(fptr10.LIBFPTR_PARAM_TAG_VALUE)
	validationResult := fptr.GetParamInt(fptr10.LIBFPTR_PARAM_MARKING_CODE_ONLINE_VALIDATION_RESULT)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_COMMODITY_NAME, "Афобазол")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_PRICE, 450)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_QUANTITY, 1.000)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_MEASUREMENT_UNIT, fptr10.LIBFPTR_IU_PIECE)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_TAX_TYPE, fptr10.LIBFPTR_TAX_VAT10)
	fptr.SetParam(1212, 33)
	fptr.SetParam(1214, 4)
	fptr.SetParam(1260, industryInfo)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_MARKING_CODE, "mark")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_MARKING_CODE_STATUS, "status")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_MARKING_CODE_TYPE, fptr10.LIBFPTR_MCT12_AUTO)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_MARKING_CODE_ONLINE_VALIDATION_RESULT, validationResult)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_MARKING_PROCESSING_MODE, 0)

	if err := fptr.Registration(); err != nil {
		log.Println(err)
		return
	}

	// Регистрация итога (отбрасываем копейки)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_SUM, 369.0)
	if err := fptr.ReceiptTotal(); err != nil {
		log.Println(err)
		return
	}

	// Оплата наличными
	fptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_TYPE, fptr10.LIBFPTR_PT_CASH)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_PAYMENT_SUM, 1000)
	if err := fptr.Payment(); err != nil {
		log.Println(err)
		return
	}

	// Закрытие чека
	_ = fptr.CloseReceipt()
	for {
		err := fptr.CheckDocumentClosed()
		if err != nil {
			// Не удалось проверить состояние документа. Вывести пользователю текст ошибки, попросить устранить неполадку и повторить запрос
			log.Println(fptr.ErrorDescription())
			continue
		} else {
			break
		}
	}

	if !fptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_CLOSED) {
		// Документ не закрылся. Требуется его отменить (если это чек) и сформировать заново
		_ = fptr.CancelReceipt()
		return
	}

	if !fptr.GetParamBool(fptr10.LIBFPTR_PARAM_DOCUMENT_PRINTED) {
		// Можно сразу вызвать метод допечатывания документа, он завершится с ошибкой, если это невозможно
		for {
			if err := fptr.ContinuePrint(); err != nil {
				// Если не удалось допечатать документ - показать пользователю ошибку и попробовать еще раз.
				log.Printf("Не удалось напечатать документ (Ошибка \"%v\"). Устраните неполадку и повторите.", fptr.ErrorDescription())
				continue
			}
		}
	}

	// Запрос информации о закрытом чеке
	fptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_LAST_DOCUMENT)
	if err := fptr.FnQueryData(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("Fiscal Sign = %v", fptr.GetParamString(fptr10.LIBFPTR_PARAM_FISCAL_SIGN))
	log.Printf("Fiscal Document Number = %v", fptr.GetParamInt(fptr10.LIBFPTR_PARAM_DOCUMENT_NUMBER))

	// Формирование слипа ЕГАИС
	if err := fptr.BeginNonfiscalDocument(); err != nil {
		log.Println(err)
		return
	}

	fptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "ИНН: 111111111111 КПП: 222222222")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
	if err := fptr.PrintText(); err != nil {
		log.Println(err)
		return
	}

	fptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "КАССА: 1               СМЕНА: 11")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
	if err := fptr.PrintText(); err != nil {
		log.Println(err)
		return
	}

	fptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "ЧЕК: 314  ДАТА: 20.11.2017 15:39")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
	if err := fptr.PrintText(); err != nil {
		log.Println(err)
		return
	}

	fptr.SetParam(fptr10.LIBFPTR_PARAM_BARCODE, "https://check.egais.ru?id=cf1b1096-3cbc-11e7-b3c1-9b018b2ba3f7")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_BARCODE_TYPE, fptr10.LIBFPTR_BT_QR)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_SCALE, 5)
	if err := fptr.PrintBarcode(); err != nil {
		log.Println(err)
		return
	}

	if err := fptr.PrintText(); err != nil {
		log.Println(err)
		return
	}

	fptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT, "https://check.egais.ru?id=cf1b1096-3cbc-11e7-b3c1-9b018b2ba3f7")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT_WRAP, fptr10.LIBFPTR_TW_CHARS)
	if err := fptr.PrintText(); err != nil {
		log.Println(err)
		return
	}

	if err := fptr.PrintText(); err != nil {
		log.Println(err)
		return
	}

	fptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT,
		"10 58 1c 85 bb 80 99 84 40 b1 4f 35 8a 35 3f 7c "+
			"78 b0 0a ff cd 37 c1 8e ca 04 1c 7e e7 5d b4 85 "+
			"ff d2 d6 b2 8d 7f df 48 d2 5d 81 10 de 6a 05 c9 "+
			"81 74")
	fptr.SetParam(fptr10.LIBFPTR_PARAM_ALIGNMENT, fptr10.LIBFPTR_ALIGNMENT_CENTER)
	fptr.SetParam(fptr10.LIBFPTR_PARAM_TEXT_WRAP, fptr10.LIBFPTR_TW_WORDS)
	if err := fptr.PrintText(); err != nil {
		log.Println(err)
		return
	}

	if err := fptr.EndNonfiscalDocument(); err != nil {
		log.Println(err)
		return
	}

	// Отчет о закрытии смены
	fptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_CLOSE_SHIFT)
	if err := fptr.Report(); err != nil {
		log.Println(err)
		return
	}

	// Получение информации о неотправленных документах
	fptr.SetParam(fptr10.LIBFPTR_PARAM_FN_DATA_TYPE, fptr10.LIBFPTR_FNDT_OFD_EXCHANGE_STATUS)
	if err := fptr.FnQueryData(); err != nil {
		log.Println(err)
		return
	}
	log.Printf("Unsent documents count = %v", fptr.GetParamInt(fptr10.LIBFPTR_PARAM_DOCUMENTS_COUNT))
	log.Printf("First unsent document number = %v", fptr.GetParamInt(fptr10.LIBFPTR_PARAM_DOCUMENT_NUMBER))
	log.Printf("First unsent document date = %v", fptr.GetParamDateTime(fptr10.LIBFPTR_PARAM_DATE_TIME))

	// Завершение работы
	if err := fptr.Close(); err != nil {
		log.Println(err)
		return
	}
}
