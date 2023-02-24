package toml

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
)

const (
	DriverInfoPath = "./cookie/appconfig/.toml"
	SessionPath    = "./cookie/session/.toml"
	UserInfoPath   = "./cookie/userdata/.toml"
	ClickPath      = "./cookie/click/.toml"
)

var (
	DecodingErrorMessage = "Ошибка декодирования"
	OpenFileErrorMessage = "Ошибка открытия файла для записи"
	EncodingErrorMessage = "Ошибка кодирования"
)

type TomlError struct {
	Message string
	Err     error
}

func NewTomlError(message string, err error) *TomlError {
	return &TomlError{Message: message, Err: err}
}

func (e *TomlError) Error() string {
	return fmt.Sprintf("Message: %s Error: %s", e.Message, e.Err)
}

func ReadToml(path string, structure interface{}) error {
	if _, err := os.Stat(path); err != nil {
		os.Create(path)
	}

	_, err := toml.DecodeFile(path, structure)

	if err != nil {
		return NewTomlError(DecodingErrorMessage, fmt.Errorf("data type: %v, error: %w", reflect.TypeOf(structure).String(), err))
	}

	return nil
}

func WriteToml(path string, structure interface{}) error {
	if _, err := os.Stat(path); err != nil {
		os.Create(path)
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0777)
	defer file.Close()

	if err != nil {
		return NewTomlError(OpenFileErrorMessage, fmt.Errorf("data type: %v, error: %w", reflect.TypeOf(structure).String(), err))
	}

	file.Truncate(0)
	err = toml.NewEncoder(file).Encode(structure)

	if err != nil {
		return NewTomlError(EncodingErrorMessage, fmt.Errorf("data type: %v, error: %w", reflect.TypeOf(structure).String(), err))
	}

	return nil
}
