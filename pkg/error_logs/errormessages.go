package error_logs

// # client error message
var (
	StatusCodeErrorMessage               = "Некорректный код статуса ответа"
	JsonUnmarshallingErrorMessage        = "Ошибка во время десериализации"
	ReadBodyErrorMessage                 = "Ошибка при чтении тела ответа"
	CreateRequestErrorMessage            = "Ошибка в создании запроса"
	EmptyURLErrorMessage                 = "Отсутствует URL в данных конфигурации"
	IncorrectLoginOrPasswordErrorMessage = "Введены некорректные логин или пароль"
	ProcessingRequestErrorMessage        = "Ошибка при выполнении запроса"
)

//@ Error Handling: Message -> Err -> Status Code (IE) -> Dependency (EXT)
//@ Формирование сообщения: Во время работы с Dependence произошла ошибка: \nMessage
