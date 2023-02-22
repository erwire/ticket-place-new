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

func (c *ClientError) Error() string {
	if c.ClientErrorType == &ResponseError {
		return fmt.Sprintf("Status Code: %d Message: %s Error: %s", c.StatusCode, c.Message, c.ClientError)
	}
	return fmt.Sprintf("Message: %s Error: %s", c.Message, c.ClientError)
}
