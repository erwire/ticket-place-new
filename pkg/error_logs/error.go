package error_logs

const (
	EmptyURLDataError       = "поле URL пустое"
	IncorrectURLDataError   = "некорректный адрес URL"
	DataNilError            = "для данного запроса необходимы параметры"
	LoginOrPasswordNilError = "не заполнены данные по логину и паролю в структуре"
	AuthorizationError      = "во время выполнения авторизации произошла ошибка: %s"
	JsonUnmarshalError      = "во время десериализации произошла ошибка: %s"
	ResponseCodeError       = "ответ от сервера по адресу: %s. статус код: %d"
)
