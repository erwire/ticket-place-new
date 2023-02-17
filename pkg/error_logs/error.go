package error_logs

import "errors"

var (
	AuthorizationError = errors.New("an error occurred during authorization")
) //+ Ошибки по модулям

var (
	ReadingBodyError       = errors.New("error during body reading")
	ResponseError          = errors.New("error response from the server")
	RequestCreatingError   = errors.New("error during creating a request")
	JsonUnmarshalError     = errors.New("an error occurred during unmarshalling")
	InvalidLoginOrPassword = errors.New("login and password data in the structure are not filled in or filled incorrect")
	DataNilError           = errors.New("there is no data for the request that needs an additional parameter")
	IncorrectURLDataError  = errors.New("unhandled URL-endpoint")
	EmptyURLDataError      = errors.New("empty URL field")
) // + Ошибки, связанные с работой модуля клиента

var (
	DecodingTomlError = errors.New("error during decoding data from toml file")
	EncodingTomlError = errors.New("error during encoding data into toml file")
) // + Ошибки, связанные с работой модуля TOML-кодирования

var (
	ShiftIsNotOpenedError   = errors.New("shift is not opened")
	ShiftIsExpired          = errors.New("shift is expired")
	ShiftIsOpenError        = errors.New("shift is opened")
	BoxOfficeIsOpenError    = errors.New("box office is open")
	BoxOfficeIsNotOpenError = errors.New("box office is not opened")
	OpenReceiptError        = errors.New("error during open receipt")
	DocumentNotClosed       = errors.New("error during process of closing document")
	CantCancelReceipt       = errors.New("error during closing receipt")
)
