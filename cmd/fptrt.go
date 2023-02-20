package main

import (
	"errors"
	"fmt"
	apperr "fptr/internal/error_list"
)

func main() {

	err := apperr.NewClientError("!", errors.New("!!!"))
	fmt.Println(err)
	fmt.Println(err.Error())
	fmt.Println(errors.Unwrap(err))

	//err := errors.New("json unmarshalling error")
	//err = NewClientError("Получены некорректные данные от сервера", err, http.StatusUnprocessableEntity)
	//if err != nil {
	//	log.Printf("Ошибка при выполнении печати чека: %s", err.Error()) //на уровень сервиса
	//	ErrorHandler(err)                                                //на уровень контроллера
	//}
	//err2 := errors.New("request create error")
	//err = error(NewClientError(RequestErrorMessage, err2))
	//if err != nil {
	//	log.Printf("Ошибка при выполнении печати чека: %s", err.Error())
	//	ErrorHandler(err)
	//}
}

//func ErrorHandler(err error) {
//	switch err.(type) {
//	case *ClientError:
//		ClientErrorHandler(err.(*ClientError))
//	}
//}
//
//func ClientErrorHandler(err *ClientError) {
//
//	switch err.ClientErrorType {
//	case &ResponseError:
//		ResponseErrorHandler(err)
//	case &RequestError:
//		RequestErrorHandler(err)
//	}
//
//}
//
//var RequestErrorMessage = "some request error message"
//
//func RequestErrorHandler(err *ClientError) {
//	switch err.Message {
//	case RequestErrorMessage:
//		fmt.Println("Ошибка при создании запроса на сервер")
//	}
//}
//
//func ResponseErrorHandler(err *ClientError) {
//	switch err.StatusCode {
//	case http.StatusBadRequest:
//		fmt.Println("На сервер отправлены некорректные данные")
//	case http.StatusUnprocessableEntity:
//		fmt.Println("От сервера получены некорректные данные, операция невозможна")
//	}
//
//}
