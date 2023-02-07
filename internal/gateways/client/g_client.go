package listener

import (
	"encoding/json"
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

type ClientGateway struct {
	info *entities.Info
}

func (l *ClientGateway) Authorization() error {
	return nil
}

func (l *ClientGateway) MakeRequest(url string, method string, structure interface{}, data ...interface{}) error {
	if l.info.AppConfig.Driver.Connection == "" { // не добавлено поле URL в настройках программы
		return fmt.Errorf("%w", errorlog.EmptyURLDataError)

	}
	var orderID = ""
	for key, value := range data { //парсинг дополнительных данных
		if key == 0 {
			orderID = fmt.Sprint(value) //добавить проверку на число
		}
		continue
	}

	var requestURL = l.info.AppConfig.Driver.Connection

	switch url {
	case refoundURL, orderURL:
		if orderID == "" {
			return fmt.Errorf("%w URL: %s", errorlog.DataNilError, requestURL)

		}
		requestURL += fmt.Sprintf(url, orderID)

	case authorizationURL:
		if l.info.AppConfig.User.Login == "" || l.info.AppConfig.User.Password == "" {
			return fmt.Errorf("%w URL: %s", errorlog.LoginOrPasswordNilError, requestURL)
		}
		requestURL += fmt.Sprintf(url, l.info.AppConfig.User.Login, l.info.AppConfig.User.Password)

	case lastConditionURL:
		requestURL += fmt.Sprintf(url, l.info.Session.UserData.ID)
	default:
		return fmt.Errorf("%w URL: %s", errorlog.IncorrectURLDataError, requestURL)
	}

	client := http.Client{Timeout: 2 * time.Second}
	request, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		return fmt.Errorf("%w URL: %s", errorlog.RequestCreatingError, requestURL)
	}

	switch url {

	case refoundURL, orderURL:
		request.Header.Add("Authorization", fmt.Sprintf("%s %s", l.info.Session.Token.TokenType, l.info.Session.Token.AccessToken))
	case lastConditionURL:
		request.Header.Add("Authorization", fmt.Sprintf("%s %s", l.info.Session.Token.TokenType, l.info.Session.Token.AccessToken))
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusAccepted, http.StatusCreated, http.StatusOK:
	default:
		return fmt.Errorf("%w: URL: %s, Status Code: %d\n", errorlog.ResponseError, requestURL, response.StatusCode)
	}

	defer response.Body.Close()
	var bodyBytes []byte
	_, err = response.Body.Read(bodyBytes)

	err = json.Unmarshal(bodyBytes, &structure)
	if err != nil {
		return fmt.Errorf("%w URL: %s", errorlog.JsonUnmarshalError, requestURL)
	}
	return nil
}
