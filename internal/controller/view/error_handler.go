package view

import (
	"errors"
	"fmt"
	apperr "fptr/internal/error_list"
	errorlog "fptr/pkg/error_logs"
	"fptr/pkg/fptr10"
	"fptr/pkg/toml"
	"net/http"
	"strings"
)

var WarningBeep = "warning_beep"
var ErrorBeep = "error_beep"

type ResponsibilityType string

var SellResponsibility = "SellResponsibility"
var RefoundResponsibility = "RefoundResponsibility"
var LoginResponsibility = "LoginResponsibility"
var ClickResponsibility = "ClickResponsibility"
var FunctionResponsibility = "FunctionResponsibility"
var NewUser = "NewUser"

func (f *FyneApp) ErrorHandler(err error, dependence string) {

	switch err.(type) {
	case *apperr.ClientError:
		f.ClientErrorHandler(err.(*apperr.ClientError), dependence)
	case *fptr10.Error:
		f.FPTRErrorHandler(err.(*fptr10.Error), dependence)
	case *toml.TomlError:
		go f.Beep(ErrorBeep)
		f.TomlErrorHandler(err.(*toml.TomlError))
	case *apperr.BusinessError:
		//f.Beep(WarningBeep)
		f.BusinessErrorHandler(err.(*apperr.BusinessError))
	default:
		go f.Beep(WarningBeep)
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

func (f *FyneApp) ClientErrorHandler(err *apperr.ClientError, dependence string) {
	switch err.ClientErrorType {
	case &apperr.RequestError:
		go f.Beep(ErrorBeep)
		f.RequestErrorHandler(err, dependence)
	case &apperr.ResponseError:
		f.ResponseErrorHandler(err, dependence)
	default:
		//! неклассифицированные ошибки
	}
} //# Обработка ошибок клиента

func (f *FyneApp) RequestErrorHandler(err *apperr.ClientError, dependence string) {
	switch err.Message {
	case errorlog.CreateRequestErrorMessage:
		if dependence == LoginResponsibility {
			f.Logout()
			f.ShowWarning("Ошибка в создании запроса")
		}
	case errorlog.EmptyURLErrorMessage:
		if dependence == LoginResponsibility {
			f.Logout()
			f.ShowWarning("Не заполнено поле адреса сервера. Пожалуйста, заполните данные по адресу.")
		}
	case errorlog.IncorrectLoginOrPasswordErrorMessage:
		if dependence == LoginResponsibility {
			f.Logout()
			f.ShowWarning("Некорректно заполнены поля логина или пароля")
		}
	case errorlog.ReadBodyErrorMessage:
		if dependence == LoginResponsibility {
			f.Logout()
			f.ShowWarning("Ошибка чтения ответа от сервера")
		}
	case errorlog.JsonUnmarshallingErrorMessage:
		if dependence == LoginResponsibility {
			f.Logout()
			f.ShowWarning("Ошибка десериализации")
		}
	case errorlog.ProcessingRequestErrorMessage:
		//f.Logout()
		switch true {
		case strings.Contains(errors.Unwrap(err.ClientError).Error(), "A socket operation was attempted to an unreachable network"):
			f.Logout()
			f.ShowWarning("Потерян доступ в интернет. Проверьте подключение к сети.")
		case strings.Contains(errors.Unwrap(err.ClientError).Error(), "Client.Timeout exceeded while awaiting headers"):
			//f.ShowWarning("Неправильный адрес сервера или сервер недоступен.")
			f.ShowProgresser()
		case strings.Contains(errors.Unwrap(err.ClientError).Error(), "No connection could be made because the target machine actively refused it"):
			f.ShowProgresser()
		case strings.Contains(errors.Unwrap(err.ClientError).Error(), "no such host"):
			f.ShowWarning("Нет подключения к интернету или адрес недоступен. Проверьте параметры настроек и подключение к интернету и авторизуйтесь повторно.")
			f.Logout()
		default:
			f.ShowWarning("Неправильный адрес сервера или сервер недоступен. Попробуйте сменить адрес сервера или проверьте подключение к интернету.")
		}

	}
} //# Обработка ошибок реквест-типа

func (f *FyneApp) ResponseErrorHandler(err *apperr.ClientError, dependence string) {
	switch err.Message {
	case errorlog.StatusCodeErrorMessage:
		f.ResponseStatusCodeErrorHandler(err, dependence)
	}
} //# Обработка ошибок ответа

func (f *FyneApp) ResponseStatusCodeErrorHandler(err *apperr.ClientError, dependence string) {
	switch err.StatusCode {
	case http.StatusNotFound:
		if dependence == ClickResponsibility {
			return
		}
		go f.Beep(ErrorBeep)
		f.ShowWarning("По запросу не найден заказ")

	case http.StatusUnprocessableEntity:
		go f.Beep(ErrorBeep)
		f.ShowWarning("В вашем запросе присутствуют данные, которые не могут быть обработаны сервером")
	case http.StatusForbidden:
		go f.Beep(ErrorBeep)
		switch dependence {
		case LoginResponsibility:
			f.ShowWarning("Неправильный логин или пароль")
		}
	case http.StatusBadRequest:
		go f.Beep(ErrorBeep)
		f.ShowWarning("Некорректный запрос")
	case http.StatusInternalServerError:
		go f.Beep(ErrorBeep)
		f.ShowWarning("Сервер недоступен")
		f.Logout()
	case http.StatusBadGateway:
		go f.Beep(ErrorBeep)
		f.ShowProgresser()
	default:
		go f.Beep(ErrorBeep)
		switch dependence {
		case LoginResponsibility:
			f.ShowWarning("Сервер прислал необрабатываемую ошибку, обратитесь к системному администратору")
			f.Logout()
		default:
			f.ShowWarning("Сервер прислал необрабатываемую ошибку, обратитесь к системному администратору")
		}

	}

} //# Обработка статус-кодов ответа сервера

//# Раздел обработки ККТ-ошибок

func (f *FyneApp) FPTRErrorHandler(err *fptr10.Error, dependence string) {
	switch err.ErrorDescription {
	case apperr.LibfptrErrorPortNotAvailable:
		f.Logout()
		f.ShowWarning(fmt.Sprintf("Ошибка работы кассы: %s", err.ErrorDescription))
	case apperr.LibfptrErrorPortBusy:
		f.Logout()
		f.ShowWarning(fmt.Sprintf("Ошибка работы кассы: %s", err.ErrorDescription))
	case apperr.LibfptrErrorNoConnection:
		f.Logout()
		fmt.Println("!")
		f.ShowWarning(fmt.Sprintf("Прервалось соединение с кассой. Пожалуйста, проверьте подключение кассы!"))
	case apperr.LibfptrErrorShiftExpired:
		//f.Reconnect()
		f.ShowWarning("Смена истекла, вы были авторизованы снова")
	case apperr.LibfptrErrorDeniedInClosedShift:
		//f.Reconnect()
		f.ShowWarning("Смена не открыта. Пожалуйста, авторизуйтесь повторно")
	case apperr.LibfptrErrorConnectionDisabled:
		f.Logout()
		f.ShowWarning("Соединение с кассой не установлено. Включите кассу и попробуйте попытку снова.")
	case apperr.LibfptrErrorConnectionLost:
		f.Logout()
		f.ShowWarning("Соединение с кассой потеряно. После подключения кассы нажмите кнопку \"Восстановить соединение с кассой\"")

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

func (f *FyneApp) Beep(beepType string) {
	if f.flag.SoundError {
		f.service.Beep(beepType)
	}
}
