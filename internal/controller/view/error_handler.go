package view

import (
	"errors"
	"fmt"
	apperr "fptr/internal/error_list"
	errorlog "fptr/pkg/error_logs"
	"fptr/pkg/fptr10"
	"fptr/pkg/toml"
	"log"
	"net/http"
	"strings"
)

var WarningBeep = "warning_beep"
var ErrorBeep = "error_beep"

type ResponsibilityType string

var SellResponsibility ResponsibilityType
var RefoundResponsibility ResponsibilityType
var LoginResponsibility ResponsibilityType
var ClickResponsibility ResponsibilityType
var FunctionResponsibility ResponsibilityType

func (f *FyneApp) ErrorHandler(err error, dependence ResponsibilityType) {

	switch err.(type) {
	case *apperr.ClientError:
		f.service.Beep(ErrorBeep)
		f.ClientErrorHandler(err.(*apperr.ClientError), dependence)
		log.Printf("Ошибка клиента: %v\n", err)
	case *fptr10.Error:
		f.service.Beep(ErrorBeep)
		log.Printf("Ошибка ККТ: %v\n", err)
		f.FPTRErrorHandler(err.(*fptr10.Error), dependence)
	case *toml.TomlError:
		f.service.Beep(ErrorBeep)
		log.Printf("Ошибка работы с файлами кэша")
		f.TomlErrorHandler(err.(*toml.TomlError))
	case *apperr.BusinessError:
		f.service.Beep(WarningBeep)
		log.Printf("Ошибка работы бизнес-логики")
		f.BusinessErrorHandler(err.(*apperr.BusinessError))
	default:
		f.service.Beep(WarningBeep)
		log.Printf("Ошибка не классифицирована: %v\n", err)
	}
}

func (f *FyneApp) TomlErrorHandler(err *toml.TomlError) {
	switch err.Message {
	case toml.EncodingErrorMessage:
		f.ShowWarning("Ошибка кодирования данных для записи в файлы кэша")
	case toml.DecodingErrorMessage:
		f.ShowWarning("Ошибка декодирования данных для чтения из файлов кэша")
	case toml.OpenFileErrorMessage:
		f.ShowWarning("Ошибка открытия файла кэш-данных")
	}
} //# Обработка ошибки TOML

func (f *FyneApp) ClientErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.ClientErrorType {
	case &apperr.RequestError:
		f.RequestErrorHandler(err, dependence)
	case &apperr.ResponseError:
		f.ResponseErrorHandler(err, dependence)
	default:
		//! неклассифицированные ошибки
	}
} //# Обработка ошибок клиента

func (f *FyneApp) RequestErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.Message {
	case errorlog.CreateRequestErrorMessage:
		if &dependence == &LoginResponsibility {
			f.Logout()
			f.ShowWarning("Ошибка в создании запроса")
		}
	case errorlog.EmptyURLErrorMessage:
		if &dependence == &LoginResponsibility {
			f.Logout()
			f.ShowWarning("Не заполнено поле адреса сервера. Пожалуйста, заполните данные по адресу.")
		}
	case errorlog.IncorrectLoginOrPasswordErrorMessage:
		if &dependence == &LoginResponsibility {
			f.Logout()
			f.ShowWarning("Некорректно заполнены поля логина или пароля")
		}
	case errorlog.ReadBodyErrorMessage:
		if &dependence == &LoginResponsibility {
			f.Logout()
			f.ShowWarning("Ошибка чтения ответа от сервера")
		}
	case errorlog.JsonUnmarshallingErrorMessage:
		if &dependence == &LoginResponsibility {
			f.Logout()
			f.ShowWarning("Ошибка десериализации")
		}
	case errorlog.ProcessingRequestErrorMessage:
		f.Logout()
		switch true {
		case strings.Contains(errors.Unwrap(err.ClientError).Error(), "A socket operation was attempted to an unreachable network"):
			f.ShowWarning("Потерян доступ в интернет. Проверьте подключение к сети.")
		case strings.Contains(errors.Unwrap(err.ClientError).Error(), "Client.Timeout exceeded while awaiting headers"):
			f.ShowWarning("Неправильный адрес сервера или сервер недоступен.")
		default:
			f.ShowWarning("Неправильный адрес сервера или сервер недоступен. Попробуйте сменить адрес сервера или проверьте подключение к интернету.")
		}

	}
} //# Обработка ошибок реквест-типа

func (f *FyneApp) ResponseErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.Message {
	case errorlog.StatusCodeErrorMessage:
		f.ResponseStatusCodeErrorHandler(err, dependence)
	}
} //# Обработка ошибок ответа

func (f *FyneApp) ResponseStatusCodeErrorHandler(err *apperr.ClientError, dependence ResponsibilityType) {
	switch err.StatusCode {
	case http.StatusNotFound:
		f.ShowWarning("По запросу не найден заказ")
	case http.StatusUnprocessableEntity:
		f.ShowWarning("В вашем запросе присутствуют данные, которые не могут быть обработаны сервером")
	case http.StatusForbidden:
		switch &dependence {
		case &LoginResponsibility:
			f.ShowWarning("Неправильный логин или пароль")
		}
	case http.StatusBadRequest:
		f.ShowWarning("Некорректный запрос")
	case http.StatusInternalServerError:
		f.ShowWarning("Сервер недоступен")
		f.Logout()
	default:
		switch &dependence {
		case &LoginResponsibility:
			f.ShowWarning("Сервер прислал необрабатываемую ошибку, обратитесь к системному администратору")
			f.Logout()
		default:
			f.ShowWarning("Сервер прислал необрабатываемую ошибку, обратитесь к системному администратору")
		}

	}

} //# Обработка статус-кодов ответа сервера

//# Раздел обработки ККТ-ошибок

func (f *FyneApp) FPTRErrorHandler(err *fptr10.Error, dependence ResponsibilityType) {
	switch err.ErrorDescription {
	case apperr.LibfptrErrorPortNotAvailable:
		f.Logout()
		f.ShowWarning(fmt.Sprintf("Ошибка работы кассы: %s", err.ErrorDescription))
	case apperr.LibfptrErrorPortBusy:
		f.Logout()
		f.ShowWarning(fmt.Sprintf("Ошибка работы кассы: %s", err.ErrorDescription))
	case apperr.LibfptrErrorNoConnection:
		f.Logout()
		f.ShowWarning(fmt.Sprintf("Прервалось соединение с кассой. Пожалуйста, проверьте подключение кассы!"))
	case apperr.LibfptrErrorShiftExpired:
		f.LogoutWS()
		f.ShowWarning("Смена истекла, пожалуйста, авторизуйтесь в системе снова")
	case apperr.LibfptrErrorDeniedInClosedShift:
		f.Logout()
		f.ShowWarning("Смена не открыта. Пожалуйста, авторизуйтесь повторно")
	default:
		f.ShowWarning(fmt.Sprintf("Необрабатываемая ошибка: %s", err.ErrorDescription))
	}
}

func (f *FyneApp) BusinessErrorHandler(err *apperr.BusinessError) {
	switch err.BusinessErr {
	case errorlog.ValidateError:
		//inform := dialog.NewInformation("Информация", err.Message, f.MainWindow)
		//inform.Show()
	}
}
