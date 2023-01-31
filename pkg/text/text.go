package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const textPath = "./data.json"

func ReadText(path string) (string, error) {

	var dataStruct struct {
		Text string `json:"text"`
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &dataStruct)
	if err != nil {
		return "", err
	}

	return dataStruct.Text, nil
}

func WriteText(path string, text string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var output struct {
		Text string `json:"text"`
	}

	output.Text = text

	indJson, err := json.MarshalIndent(output, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.Write(indJson)

	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(path string, data interface{}) error {

	file, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadJSON(path string, data interface{}) error {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return err
	}
	fmt.Println(data)
	return nil
}
