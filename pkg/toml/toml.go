package toml

import (
	"fmt"
	"fptr/pkg/error_logs"
	"github.com/BurntSushi/toml"
	"os"
	"reflect"
)

const (
	DriverInfoPath = "./cookie/appconfig/.toml"
	SessionPath    = "./cookie/session/.toml"
	UserInfoPath   = "./cookie/userdata/.toml"
)

func ReadToml(path string, structure interface{}) error {
	if _, err := os.Stat(path); err != nil {
		os.Create(path)
	}

	_, err := toml.DecodeFile(path, structure)

	if err != nil {
		return fmt.Errorf("%w, data type: %v, path: %s, error_description: %s", error_logs.DecodingTomlError, reflect.TypeOf(structure).String(), path, err.Error())
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
		return fmt.Errorf("%w, data type: %v, path: %s, error_description: %s", error_logs.EncodingTomlError, reflect.TypeOf(structure).String(), path, err.Error())
	}

	err = toml.NewEncoder(file).Encode(structure)

	if err != nil {
		return fmt.Errorf("ERROR: [error type: \"%w\", data type: \"%v\", path: \"%s\", error_description: \"%s\"]\n", error_logs.EncodingTomlError, reflect.TypeOf(structure).String(), path, err.Error())
	}

	return nil
}
