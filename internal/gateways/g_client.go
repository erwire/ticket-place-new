package gateways

import (
	"encoding/json"
	"fmt"
	"fptr/internal/entities"
	apperr "fptr/internal/error_list"
	errorlog "fptr/pkg/error_logs"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	OrderURL         = "/api/order/%s"                        //структура Sell
	RefoundURL       = "/api/refund/%s"                       //структура Refound
	AuthorizationURL = "/api/auth/login?email=%s&password=%s" //структура
	LastConditionURL = "/api/print-requests/by-user/%d"       // структура Click
)

const (
	sellDumpName    = "./debug_info/sell/sell.json"
	refoundDumpName = "./debug_info/refound/refound.json"
	clickDumpName   = "./debug_info/click/click.json"
	loginDumpName   = "./debug_info/login/login.json"
)

type ClientGateway struct {
	client *http.Client
}

func NewClientGateway(client *http.Client) *ClientGateway {
	return &ClientGateway{
		client: client,
	}
}

func (l *ClientGateway) SetTimeout(timeout time.Duration) {
	l.client.Timeout = timeout
}

func (l *ClientGateway) Login(config entities.AppConfig) (*entities.SessionInfo, error) {
	var session entities.SessionInfo
	if len(config.Driver.Connection) == 0 {
		return nil, apperr.NewClientError(errorlog.EmptyURLErrorMessage, errorlog.EmptyURLDataError)
	}
	if !config.User.ValidateUser() {
		return nil, apperr.NewClientError(errorlog.IncorrectLoginOrPasswordErrorMessage, errorlog.InvalidLoginOrPassword)
	}

	url := config.Driver.Connection
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	requestURI := url + fmt.Sprintf(AuthorizationURL, config.User.Login, config.User.Password)

	request, err := http.NewRequest(http.MethodPost, requestURI, nil)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.CreateRequestErrorMessage, err)
	}

	response, err := l.client.Do(request)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.ProcessingRequestErrorMessage, err)
	}

	switch response.StatusCode {

	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	case http.StatusInternalServerError:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.InternalServerError, response.StatusCode)
	default:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.DefaultHttpError, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.ReadBodyErrorMessage, err)
	}
	err = json.Unmarshal(body, &session)
	if err != nil {
		l.writeProblemDataIntoJSONDump(body, loginDumpName)
		return nil, apperr.NewClientError(errorlog.JsonUnmarshallingErrorMessage, err)
	}
	return &session, nil
}

func (l *ClientGateway) GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, error) {
	var click entities.Click

	if len(connectionURL) == 0 || session.UserData.ID == 0 {
		return nil, apperr.NewClientError(errorlog.EmptyURLErrorMessage, errorlog.EmptyURLDataError)
	}

	url := connectionURL

	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	requestURI := url + fmt.Sprintf(LastConditionURL, session.UserData.ID)

	request, err := http.NewRequest(http.MethodGet, requestURI, nil)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.CreateRequestErrorMessage, err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("%s %s", session.TokenType, session.AccessToken))

	response, err := l.client.Do(request)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.ProcessingRequestErrorMessage, err)
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	case http.StatusInternalServerError:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.InternalServerError, response.StatusCode)
	default:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.DefaultHttpError, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		l.writeProblemDataIntoJSONDump(body, clickDumpName)
		return nil, apperr.NewClientError(errorlog.ReadBodyErrorMessage, err)
	}
	err = json.Unmarshal(body, &click)
	if err != nil {
		return nil, apperr.NewClientError(errorlog.JsonUnmarshallingErrorMessage, err)
	}
	return &click, nil
}

func (l *ClientGateway) GetSell(info entities.Info, sellID string) (*entities.Sell, error) {
	var sell entities.Sell

	if len(info.AppConfig.Driver.Connection) == 0 || len(sellID) == 0 {
		return nil, apperr.NewClientError(errorlog.EmptyURLErrorMessage, errorlog.EmptyURLDataError)
	}

	url := info.AppConfig.Driver.Connection

	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	requestURI := url + fmt.Sprintf(OrderURL, sellID)

	request, err := http.NewRequest(http.MethodGet, requestURI, nil)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.CreateRequestErrorMessage, err)
	}

	request.Header.Add("Authorization", fmt.Sprintf("%s %s", info.Session.TokenType, info.Session.AccessToken))

	response, err := l.client.Do(request)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.ProcessingRequestErrorMessage, err)
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	case http.StatusInternalServerError:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.InternalServerError, response.StatusCode)
	default:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.DefaultHttpError, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.ReadBodyErrorMessage, err)
	}
	err = json.Unmarshal(body, &sell)
	if err != nil {
		l.writeProblemDataIntoJSONDump(body, sellDumpName)
		return nil, apperr.NewClientError(errorlog.JsonUnmarshallingErrorMessage, err)
	}
	return &sell, nil
}

func (l *ClientGateway) GetRefound(info entities.Info, refoundID string) (*entities.Refound, error) {
	var refound entities.Refound

	if len(info.AppConfig.Driver.Connection) == 0 || len(refoundID) == 0 {
		return nil, apperr.NewClientError(errorlog.EmptyURLErrorMessage, errorlog.EmptyURLDataError)
	}

	url := info.AppConfig.Driver.Connection

	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	requestURI := url + fmt.Sprintf(RefoundURL, refoundID)

	request, err := http.NewRequest(http.MethodGet, requestURI, nil)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.CreateRequestErrorMessage, err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("%s %s", info.Session.TokenType, info.Session.AccessToken))

	response, err := l.client.Do(request)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.ProcessingRequestErrorMessage, err)
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	case http.StatusInternalServerError:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.InternalServerError, response.StatusCode)
	default:
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.DefaultHttpError, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, apperr.NewClientError(errorlog.ReadBodyErrorMessage, err)
	}
	err = json.Unmarshal(body, &refound)
	if err != nil {
		l.writeProblemDataIntoJSONDump(body, refoundDumpName)
		return nil, apperr.NewClientError(errorlog.JsonUnmarshallingErrorMessage, err)
	}
	return &refound, nil
}

func (l *ClientGateway) writeProblemDataIntoJSONDump(body []byte, filepath string) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		// return err -> в дальнейшем избавиться от функций без обработчика ошибок. Можно не возвращать в handler, но в сервисе быть должны однозначно для логирования
	}
	if err := file.Truncate(0); err != nil {
		log.Println(err.Error())
	}

	//message := l.messageDecode(string(body))
	if _, err = file.WriteString(string(body)); err != nil {
		// return err -> в дальнейшем избавиться от функций без обработчика ошибок. Можно не возвращать в handler, но в сервисе быть должны однозначно для логирования
	} else {
		// return err -> в дальнейшем избавиться от функций без обработчика ошибок. Можно не возвращать в handler, но в сервисе быть должны однозначно для логирования
	}

	file.Close()
}
