package listener

import (
	"fptr/internal/entities"
)

const (
	orderURL         = "/api/order/%s"
	refoundURL       = "/api/refund/%s"
	authorizationURL = "/api/auth/login?email=%s&password=%s"
	lastConditionURL = "/api/print-requests/by-user/%s"
)

type Listen struct {
	info *entities.Info
}

func (l *Listen) Listen() {

}

/*func (l *Listen) MakeRequest(url string, method string, structure interface{}, data ...interface{}) error {
	if l.info.Driver.Connection == "" {
		return errors.New(error_log.EmptyURLDataError)
	}
	var orderID = ""
	for key, value := range data {
		if key == 0 {
			orderID = fmt.Sprint(value)
		}
		continue
	}

	var requestURL = l.info.Driver.Connection

	switch url {
	case refoundURL, orderURL:
		if orderID == "" {
			return errors.New(error_log.DataNilError)
		}
		requestURL += fmt.Sprintf(url, orderID)

	case authorizationURL:

	case lastConditionURL:

	default:
		return errors.New(error_log.IncorrectURLDataError)
	}
}*/
