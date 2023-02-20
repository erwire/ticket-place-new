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
)

const (
	OrderURL         = "/api/order/%s"                        //структура Sell
	RefoundURL       = "/api/refund/%s"                       //структура Refound
	AuthorizationURL = "/api/auth/login?email=%s&password=%s" //структура
	LastConditionURL = "/api/print-requests/by-user/%d"       // структура Click
)

type ClientGateway struct {
	client *http.Client
}

func NewClientGateway(client *http.Client) *ClientGateway {
	return &ClientGateway{
		client: client,
	}
}

func (l *ClientGateway) Login(config entities.AppConfig) (*entities.SessionInfo, error) {
	var session entities.SessionInfo
	if len(config.Driver.Connection) == 0 {
		return nil, apperr.NewClientError(errorlog.EmptyURLErrorMessage, errorlog.EmptyURLDataError)
	}
	if !config.User.ValidateUser() {
		return nil, apperr.NewClientError(errorlog.IncorrectLoginOrPasswordErrorMessage, errorlog.InvalidLoginOrPassword)
	}

	requestURI := config.Driver.Connection + fmt.Sprintf(AuthorizationURL, config.User.Login, config.User.Password)

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
		return nil, apperr.NewClientError(errorlog.JsonUnmarshallingErrorMessage, err, http.StatusUnprocessableEntity)
	}
	return &session, nil
}

func (l *ClientGateway) GetLastReceipt(connectionURL string, session entities.SessionInfo) (*entities.Click, error) {
	var click entities.Click

	if len(connectionURL) == 0 || session.UserData.ID == 0 {
		return nil, errorlog.EmptyURLDataError
	}

	requestURI := connectionURL + fmt.Sprintf(LastConditionURL, session.UserData.ID)

	request, err := http.NewRequest(http.MethodGet, requestURI, nil)

	if err != nil {
		return nil, errorlog.RequestCreatingError
	}

	request.Header.Add("Authorization", fmt.Sprintf("%s %s", session.TokenType, session.AccessToken))

	response, err := l.client.Do(request)

	if err != nil {
		return nil, errorlog.ResponseError
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	default:
		return nil, fmt.Errorf("%w, status code: %d", errorlog.ResponseError, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, errorlog.ReadingBodyError
	}
	err = json.Unmarshal(body, &click)
	if err != nil {
		return nil, errorlog.JsonUnmarshalError
	}
	return &click, nil
}

func (l *ClientGateway) GetSell(info entities.Info, sellID string) (*entities.Sell, error) {
	var sell entities.Sell

	if len(info.AppConfig.Driver.Connection) == 0 || len(sellID) == 0 {
		return nil, errorlog.EmptyURLDataError
	}

	requestURI := info.AppConfig.Driver.Connection + fmt.Sprintf(OrderURL, sellID)

	request, err := http.NewRequest(http.MethodGet, requestURI, nil)

	if err != nil {
		return nil, errorlog.RequestCreatingError
	}

	request.Header.Add("Authorization", fmt.Sprintf("%s %s", info.Session.TokenType, info.Session.AccessToken))

	response, err := l.client.Do(request)

	if err != nil {
		return nil, errorlog.ResponseError
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	default:
		return nil, fmt.Errorf("%w, status code: %d", errorlog.ResponseError, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, errorlog.ReadingBodyError
	}
	err = json.Unmarshal(body, &sell)
	if err != nil {
		return nil, errorlog.JsonUnmarshalError
	}
	return &sell, nil
}

func (l *ClientGateway) GetRefound(info entities.Info, refoundID string) (*entities.Refound, error) {
	var refound entities.Refound

	if len(info.AppConfig.Driver.Connection) == 0 || len(refoundID) == 0 {
		return nil, errorlog.EmptyURLDataError
	}

	requestURI := info.AppConfig.Driver.Connection + fmt.Sprintf(RefoundURL, refoundID)

	request, err := http.NewRequest(http.MethodGet, requestURI, nil)

	if err != nil {
		return nil, errorlog.RequestCreatingError
	}
	request.Header.Add("Authorization", fmt.Sprintf("%s %s", info.Session.TokenType, info.Session.AccessToken))

	response, err := l.client.Do(request)

	if err != nil {
		return nil, errorlog.ResponseError
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	default:
		return nil, fmt.Errorf("%w, status code: %d", errorlog.ResponseError, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	file, err := os.OpenFile("./cookie/refound/.json", os.O_WRONLY, 0660)
	if err != nil {
		log.Println(err.Error())
	}

	if _, err = file.Write(body); err != nil {
		log.Println(err.Error())
	}

	file.Close()

	if err != nil {
		return nil, errorlog.ReadingBodyError
	}
	err = json.Unmarshal(body, &refound)
	if err != nil {
		return nil, fmt.Errorf("%w %v", errorlog.JsonUnmarshalError, err)
	}
	return &refound, nil
}
