package listener

import (
	"encoding/json"
	"errors"
	"fmt"
	"fptr/internal/entities"
	errorlog "fptr/pkg/error_logs"
	"net/http"
	"time"
)

const (
	orderURL         = "/api/order/%s"                        //структура Sell
	refoundURL       = "/api/refund/%s"                       //структура Refound
	authorizationURL = "/api/auth/login?email=%s&password=%s" //структура
	lastConditionURL = "/api/print-requests/by-user/%s"       // структура Click
)

type Listen struct {
	info *entities.Info
}

func (l *Listen) Listen() {
	listenedData := &entities.Click{}
	l.MakeRequest(lastConditionURL, http.MethodGet, &listenedData)
}

func (l *Listen) MakeRequest(url string, method string, structure interface{}, data ...interface{}) error {
	if l.info.Driver.Connection == "" { // не добавлено поле URL в настройках программы
		return errors.New(errorlog.EmptyURLDataError)
	}
	var orderID = ""
	for key, value := range data { //парсинг дополнительных данных
		if key == 0 {
			orderID = fmt.Sprint(value) //добавить проверку на число
		}
		continue
	}

	var requestURL = l.info.Driver.Connection

	switch url {
	case refoundURL, orderURL:
		if orderID == "" {
			return errors.New(errorlog.DataNilError)
		}
		requestURL += fmt.Sprintf(url, orderID)

	case authorizationURL:
		if l.info.User.Login == "" || l.info.User.Password == "" {
			return errors.New(errorlog.LoginOrPasswordNilError)
		}
		requestURL += fmt.Sprintf(url, l.info.User.Login, l.info.User.Password)

	case lastConditionURL:
		requestURL += fmt.Sprintf(url, l.info.Logined.UserData.ID)
	default:
		return errors.New(errorlog.IncorrectURLDataError)
	}

	client := http.Client{Timeout: 2 * time.Second}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	switch url {

	case refoundURL, orderURL:
		request.Header.Add("Authorization", fmt.Sprintf("%s %s", l.info.Logined.TokenType, l.info.Logined.AccessToken))
	case lastConditionURL:
		request.Header.Add("Authorization", fmt.Sprintf("%s %s", l.info.Logined.TokenType, l.info.Logined.AccessToken))
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	switch response.StatusCode {

	case http.StatusForbidden:
		return errors.New(fmt.Sprintf(errorlog.ResponseCodeError, l.info.Driver.Connection, response.StatusCode))
	case http.StatusBadGateway:
		return errors.New(fmt.Sprintf(errorlog.ResponseCodeError, l.info.Driver.Connection, response.StatusCode))
	case http.StatusNotFound:
		return errors.New(fmt.Sprintf(errorlog.ResponseCodeError, l.info.Driver.Connection, response.StatusCode))
	case http.StatusInternalServerError:
		return errors.New(fmt.Sprintf(errorlog.ResponseCodeError, l.info.Driver.Connection, response.StatusCode))
	}
	defer response.Body.Close()
	var bodyBytes []byte
	_, err = response.Body.Read(bodyBytes)

	err = json.Unmarshal(bodyBytes, &structure)
	if err != nil {
		return errors.New(fmt.Sprintf(errorlog.JsonUnmarshalError, err.Error()))
	}
	return nil
}
