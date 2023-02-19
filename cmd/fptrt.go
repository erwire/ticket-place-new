package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type ClientErrorType string

var ResponseError ClientErrorType
var RequestError ClientErrorType

type ClientError struct {
	*ClientErrorType
	Message     string
	ClientError error
	StatusCode  int
}

func NewClientError(message string, err error, code ...int) *ClientError {
	if len(code) == 0 {
		return &ClientError{ClientErrorType: &RequestError, Message: message, ClientError: err}
	}
	return &ClientError{ClientErrorType: &ResponseError, Message: message, StatusCode: code[0], ClientError: err}
}

func (c *ClientError) Error() string {
	if c.ClientErrorType == &ResponseError {
		return fmt.Sprintf("Status Code: %d Message: %s Error: %s", c.StatusCode, c.Message, c.ClientError)
	}
	return fmt.Sprintf("Message: %s Error: %s", c.Message, c.ClientError)
}

func main() {
	err := errors.New("json unmarshalling error")
	err = NewClientError("Получены некорректные данные от сервера", err, http.StatusUnprocessableEntity)
	if err != nil {
		log.Printf("Ошибка при выполнении печати чека: %s", err.Error()) //на уровень сервиса
		ErrorHandler(err)                                                //на уровень контроллера
	}
	err2 := errors.New("request create error")
	err = error(NewClientError(RequestErrorMessage, err2))
	if err != nil {
		log.Printf("Ошибка при выполнении печати чека: %s", err.Error())
		ErrorHandler(err)
	}
}

func ErrorHandler(err error) {
	switch err.(type) {
	case *ClientError:
		ClientErrorHandler(err.(*ClientError))
	}
}

func ClientErrorHandler(err *ClientError) {

	switch err.ClientErrorType {
	case &ResponseError:
		ResponseErrorHandler(err)
	case &RequestError:
		RequestErrorHandler(err)
	}

}

var RequestErrorMessage = "some request error message"

func RequestErrorHandler(err *ClientError) {
	switch err.Message {
	case RequestErrorMessage:
		fmt.Println("Ошибка при создании запроса на сервер")
	}
}

func ResponseErrorHandler(err *ClientError) {
	switch err.StatusCode {
	case http.StatusBadRequest:
		fmt.Println("На сервер отправлены некорректные данные")
	case http.StatusUnprocessableEntity:
		fmt.Println("От сервера получены некорректные данные, операция невозможна")
	}

}
