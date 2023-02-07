package error_logs

import "errors"

var (
	ResponseError           = errors.New("error response from the server")
	RequestCreatingError    = errors.New("error during creating a request")
	JsonUnmarshalError      = errors.New("an error occurred during unmarshalling")
	AuthorizationError      = errors.New("an error occurred during authorization")
	LoginOrPasswordNilError = errors.New("login and password data in the structure are not filled in")
	DataNilError            = errors.New("there is no data for the request that needs an additional parameter")
	IncorrectURLDataError   = errors.New("unhandled URL-endpoint")
	EmptyURLDataError       = errors.New("empty URL field")
) // + Ошибки, связанные с работой модуля клиента

var (
	DecodingTomlError = errors.New("error during decoding data from toml file")
	EncodingTomlError = errors.New("error during encoding data into toml file")
) // + Ошибки, связанные с работой модуля TOML-кодирования
