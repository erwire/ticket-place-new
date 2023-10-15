package entities

import (
	"strconv"
	"time"
)

const (
	DoneStatus  = "Выполнен"
	ErrorStatus = "Не выполнен"
)

type SellsDTO struct {
	SellID string `json:"sell_id,omitempty" db:"sell_id"`
	Date   string `json:"date,omitempty" db:"date"`
	Amount uint64 `json:"amount,omitempty" db:"amount"`
	Status string `json:"status,omitempty" db:"status"`
	Error  string `json:"error,omitempty" db:"error"`
	Event  string `json:"event,omitempty" db:"event"`
}

func GetSellsDTO(sell Sell, status string, err string) *SellsDTO {
	amount := 0
	for _, ticket := range sell.Data.Tickets {
		amount += ticket.Amount
	}

	return &SellsDTO{
		SellID: strconv.Itoa(sell.Data.OrderId),
		Date:   time.Now().Format("02.01.2006 15:04:05"),
		Amount: uint64(amount),
		Status: status,
		Error:  err,
		Event:  sell.Data.Event.Name,
	}
}

func NewSellsDTO() *SellsDTO {
	return &SellsDTO{}
}
