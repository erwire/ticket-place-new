package entities

import (
	"net/mail"
	"time"
)

type Info struct {
	AppConfig AppConfig   `toml:"app_config"`
	Session   SessionInfo `toml:"session_info"`
}

type AppConfig struct {
	User   UserInfo   `toml:"user_info"`
	Driver DriverInfo `toml:"driver_info"`
}

type UserInfo struct {
	Login     string    `toml:"login"`
	Password  string    `toml:"password"`
	TaxesInfo TaxesInfo `toml:"taxes"`
}

type DriverInfo struct {
	Path          string        `toml:"path"`
	Com           string        `toml:"com"`
	Time          string        `toml:"time"`
	Connection    string        `toml:"baseurl"`
	PollingPeriod time.Duration `toml:"polling_period"`
	TimeoutPeriod time.Duration `toml:"timeout_duration"`
	UpdatePath    string        `toml:"update_path"`
}

type TaxesCalculationType int

const (
	UndefinedTaxes TaxesCalculationType = iota
	NoTaxes
	TaxesValue0
	TaxesValue5
	TaxesValue7
	TaxesValue10
	TaxesValue20
	TaxesValue105
	TaxesValue107
	TaxesValue110
	TaxesValue120
)

func NewCalculationTypeList() []TaxesCalculationType {
	return []TaxesCalculationType{
		NoTaxes,
		TaxesValue0,
		TaxesValue5,
		TaxesValue7,
		TaxesValue10,
		TaxesValue20,
		TaxesValue105,
		TaxesValue107,
		TaxesValue110,
		TaxesValue120,
	}
}

func NewCalculationType(t string) TaxesCalculationType {
	switch t {
	case "Без НДС":
		return NoTaxes
	case "0% НДС":
		return TaxesValue0
	case "5% НДС":
		return TaxesValue5
	case "7% НДС":
		return TaxesValue7
	case "10% НДС":
		return TaxesValue10
	case "20% НДС":
		return TaxesValue20
	case "НДС рассчитанный 5/105":
		return TaxesValue105
	case "НДС рассчитанный 7/107":
		return TaxesValue107
	case "НДС рассчитанный 10/110":
		return TaxesValue110
	case "НДС рассчитанный 20/120":
		return TaxesValue120
	}

	return TaxesCalculationType(0)
}

func (t TaxesCalculationType) String() string {
	switch t {
	case NoTaxes:
		return "Без НДС"
	case TaxesValue0:
		return "0% НДС"
	case TaxesValue5:
		return "5% НДС"
	case TaxesValue7:
		return "7% НДС"
	case TaxesValue10:
		return "10% НДС"
	case TaxesValue20:
		return "20% НДС"
	case TaxesValue105:
		return "НДС рассчитанный 5/105"
	case TaxesValue107:
		return "НДС рассчитанный 7/107"
	case TaxesValue110:
		return "НДС рассчитанный 10/110"
	case TaxesValue120:
		return "НДС рассчитанный 20/120"
	default:
		return "Неизвестный"
	}
}

type TaxesInfo struct {
	Taxes TaxesCalculationType `toml:"taxes_type"`
}

type SessionInfo struct {
	AccessToken string    `json:"accessToken" toml:"access_token"`
	TokenType   string    `json:"token_type" toml:"token_type"`
	CreatedAt   time.Time `toml:"created_at"`

	UserData struct {
		ID       int    `json:"id" toml:"id"`
		FullName string `json:"fullName" toml:"full_name"`
		Username string `json:"username" toml:"user_name"`
		Avatar   string `json:"avatar" toml:"avatar"`
		Email    string `json:"email" toml:"email"`
		Role     string `json:"role" toml:"role"`
		Inn      uint64 `json:"inn" toml:"inn"`
	} `json:"userData" toml:"user_data"`
}

func (u *UserInfo) ValidateUser() bool {
	_, err := mail.ParseAddress(u.Login)
	return err == nil && u.Password != ""
}

func (s *SessionInfo) IsDead() bool {
	return time.Since(s.CreatedAt).Hours() > 24
}
