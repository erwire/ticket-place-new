package apperr

import "fmt"

type ClientErrorType string

var ResponseError ClientErrorType
var RequestError ClientErrorType

type ClientError struct {
	*ClientErrorType
	Message     string
	ClientError error
	StatusCode  int
}

func NewClientError(message string, err error, code ...int) *ClientError {
	if len(code) == 0 {
		return &ClientError{ClientErrorType: &RequestError, Message: message, ClientError: err}
	}
	return &ClientError{ClientErrorType: &ResponseError, Message: message, StatusCode: code[0], ClientError: err}
}

func (e *ClientError) Error() string {
	if e.ClientErrorType == &ResponseError {
		return fmt.Sprintf("Status Code: %d Message: %s Error: %s", e.StatusCode, e.Message, e.ClientError)
	}
	return fmt.Sprintf("Message: %s Error: %s", e.Message, e.ClientError)
}

type BusinessError struct {
	Message     string
	BusinessErr error
}

func NewBusinessError(message string, businessErr error) *BusinessError {
	return &BusinessError{Message: message, BusinessErr: businessErr}
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("Message: %s Error: %v", e.Message, e.BusinessErr)
}

type DatabaseError struct {
	Message string
	Err     error
}

func NewDatabaseError(message string, error error) *DatabaseError {
	return &DatabaseError{Message: message, Err: error}
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("Message: %s Error: %v", e.Message, e.Err)
}
