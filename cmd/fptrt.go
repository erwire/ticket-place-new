package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {

	myApp := app.New()
	w := myApp.NewWindow("Image")

	w.SetContent(container.NewVBox(boxCenter, widget.NewLabel("!!!")))

	w.ShowAndRun()

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
