package entities

import (
	"encoding/base64"
	"github.com/skip2/go-qrcode"
	"html/template"
	"strings"
	"time"
)

type OrderDTO struct {
	Id             int         `json:"id"`           //r.Data.ID
	Type           string      `json:"type"`         //
	CashierName    string      `json:"cashier_name"` //r.Data.KassirName
	PaymentType    string      `json:"payment_type"` //r.Data.PaymentType
	Ticket         []TicketDTO `json:"ticket"`
	AdditionalText string      `json:"additional_text"`
	AdditionalData []template.HTML
}

type TicketDTO struct {
	Id            int       `json:"id"`
	Status        string    `json:"status"`
	Number        string    `json:"number"`
	OrderDate     time.Time `json:"order_date"`
	OrderTimeSep  string    `json:"order_time_sep"`
	OrderDateSep  string    `json:"order_date_sep"`
	EventManager  string    `json:"event_manager"`
	EventAddress  string    `json:"event_address"`
	EventName     string    `json:"event_name"`
	EventAgeLimit string    `json:"event_age_limit"`
	EventDate     time.Time `json:"event_date"`
	EventTimeSep  string    `json:"event_time_sep"`
	EventDateSep  string    `json:"event_date_sep"`
	SeatRow       string    `json:"seat_row"`
	SeatPlace     string    `json:"seat_place"`
	SeatZona      string    `json:"seat_zona"`
	Amount        int       `json:"amount"`
	QRBase64      string    `json:"qrcode_base64"`
}

func NewTicketDTO() *TicketDTO {
	return &TicketDTO{}
}

func NewOrderDTO() *OrderDTO {
	return &OrderDTO{}
}

type RequestToPrintDataDTO struct {
	PageSize        string `json:"page_size"`
	PageOrientation string `json:"page_orientation"`
	Ticket          string `json:"ticket"`
	PrinterName     string `json:"printer_name"`
}

func NewRequestToPrintDataDTO() *RequestToPrintDataDTO {
	return &RequestToPrintDataDTO{}
}

type PageParamsDTO struct {
	PageSize        string `json:"page_size"`
	PageOrientation string `json:"page_orientation"`
}

func NewPageParamsDTO(pageSize string, pageOrientation string) *PageParamsDTO {
	return &PageParamsDTO{PageSize: pageSize, PageOrientation: pageOrientation}
}

func NewPrinterParamsDTO() *RequestToPrintDataDTO {
	return &RequestToPrintDataDTO{}
}

func (s *OrderDTO) InitAdditionalSeparator() {
	data := strings.Split(s.AdditionalText, "\n\n")
	s.AdditionalData = make([]template.HTML, 3)

	for key := range data {
		s.AdditionalData[key] = template.HTML(strings.Replace(data[key], "\n", "<br>", -1))
	}

}

func (s *OrderDTO) InitQR() error {
	for key := range s.Ticket {

		qr, err := qrcode.New(s.Ticket[key].Number, qrcode.Highest)
		qr.DisableBorder = true
		qrbyytes, err := qr.PNG(90)
		//qr, err := qrcode.Encode(, , 128)
		if err != nil {
			return err
		}

		s.Ticket[key].QRBase64 = base64.StdEncoding.EncodeToString(qrbyytes)
	}
	return nil
}

func (s *OrderDTO) InitDateSeparate() {
	for key := range s.Ticket {
		s.Ticket[key].EventDateSep = s.Ticket[key].EventDate.Format("02.01.2006")
		s.Ticket[key].EventTimeSep = s.Ticket[key].EventDate.Format("15:04")
		s.Ticket[key].OrderDateSep = s.Ticket[key].OrderDate.Format("02.01.2006")
		s.Ticket[key].OrderTimeSep = s.Ticket[key].OrderDate.Format("15:04")
	}
}
