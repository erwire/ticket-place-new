package toml

import (
	"fptr/internal/entities"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

const (
	tomlPath = "./toml/.toml"
)

func InitializeToml() (entities.Info, error) {
	if _, err := os.Stat(tomlPath); err != nil {
		os.Create(tomlPath)
	}
	var info entities.Info

	_, err := toml.DecodeFile(tomlPath, &info)

	if err != nil {
		log.Println("Ошибка при декодировании конфигурации: " + err.Error())
		return info, err
	}

	return info, nil
}

func GetDataFromToml() {
	
}

func PutDataIntoToml() {

}
