package error_logs

//client

var ( //! error
	StatusCodeErrorMessage               = "Некорректный код статуса ответа"
	JsonUnmarshallingErrorMessage        = "Ошибка во время десериализации"
	ReadBodyErrorMessage                 = "Ошибка при чтении тела ответа"
	CreateRequestErrorMessage            = "Ошибка в создании запроса"
	EmptyURLErrorMessage                 = "Отсутствует URL в данных конфигурации"
	IncorrectLoginOrPasswordErrorMessage = "Введены некорректные логин или пароль"
	ProcessingRequestErrorMessage        = "Ошибка при выполнении запроса"
)
