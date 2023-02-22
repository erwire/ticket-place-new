package view

import (
	apperr "fptr/internal/error_list"
	errorlog "fptr/pkg/error_logs"
	"fptr/pkg/fptr10"
	"net/http"
)

type ResponsibilityType string

var SellResponsibility ResponsibilityType
var RefoundResponsibility ResponsibilityType
var LoginResponsibility ResponsibilityType
var ClickResponsibility ResponsibilityType

func (f *FyneApp) ErrorHandler(err error, dependence ResponsibilityType) {
	switch err.(type) {
	case *apperr.ClientError:
		f.ClientErrorHandler(err.(*apperr.ClientError), dependence)
	case *fptr10.Error:
		f.FPTRErrorHandler(err.(*fptr10.Error), dependence)
	}
}

//# Обработка ошибок клиента

func (f *FyneApp) ClientErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.ClientErrorType {
	case &apperr.RequestError:
		f.RequestErrorHandler(err, dependence)
	case &apperr.ResponseError:
		f.ResponseErrorHandler(err, dependence)
	default:
		//! неклассифицированные ошибки
	}
}

func (f *FyneApp) RequestErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.Message {

	}
}

func (f *FyneApp) ResponseErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.Message {
	case errorlog.StatusCodeErrorMessage:

	}
}

func (f *FyneApp) ResponseStatusCodeErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.StatusCode {
	case http.StatusNotFound:
	case http.StatusUnprocessableEntity:
	}

}

//# Раздел обработки ККТ-ошибок

func (f *FyneApp) FPTRErrorHandler(err *fptr10.Error, dependence ResponsibilityType) {
	switch err.ErrorCode {

	}
}
