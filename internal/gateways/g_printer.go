package gateways

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"fptr/internal/entities"
	apperr "fptr/internal/error_list"
	errorlog "fptr/pkg/error_logs"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const PrinterJarPath = "./modules/print_service/Printer_Service.jar"

const ExamplePath = "./content/system/html_template/ticket.html"

type PrinterInterface interface {
	StartService(config entities.DriverInfo) error
	Ping(config entities.DriverInfo) error
	Print(config entities.DriverInfo, dto entities.OrderDTO, pp entities.PageParamsDTO) error
	StopPrinterServer()
	GetPrinterList(config entities.DriverInfo) ([]string, error)
}

type MessageDTO struct {
}

type Printer struct {
	client *http.Client
}

func NewPrinter() *Printer {
	return &Printer{client: &http.Client{Timeout: 10 * time.Second}}
}

func (p *Printer) Ping(config entities.DriverInfo) error {
	address := fmt.Sprintf("http://localhost:%s/ping", config.PrinterServiceAddress)
	log.Println(address)
	request, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return err
	}
	response, err := p.client.Do(request)
	if err != nil {
		return err
	}
	fmt.Println(response.StatusCode)
	if response.StatusCode != 200 {
		return apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.InternalServerError, response.StatusCode)
	}

	return nil
}

func (p *Printer) Print(config entities.DriverInfo, dto entities.OrderDTO, pp entities.PageParamsDTO) error {
	data, err := p.formHtmlFromData(dto)

	if err != nil {
		return err
	}

	requestData := entities.NewRequestToPrintDataDTO()

	requestData.Ticket = data
	requestData.PrinterName = config.PrinterName
	requestData.PageOrientation = pp.PageOrientation
	requestData.PageSize = pp.PageSize

	log.Println()

	message, err := json.Marshal(requestData)

	if err != nil {
		return err
	}
	body := bytes.NewReader(message)
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/print", config.PrinterServiceAddress), body)
	if err != nil {
		return err
	}
	response, err := p.client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.InternalServerError, response.StatusCode)
	}

	return nil
}

func (p *Printer) formHtmlFromData(dto entities.OrderDTO) (string, error) {
	err := dto.InitQR()
	if err != nil {
		return "", err
	}

	dto.InitDateSeparate()
	dto.InitAdditionalSeparator()
	tmpl, err := template.ParseFiles(ExamplePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, dto)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *Printer) StartService(config entities.DriverInfo) error {
	str, err := filepath.Abs(PrinterJarPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("javaw", "-jar", str)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("port=%s", config.PrinterServiceAddress))

	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (p *Printer) GetPrinterList(config entities.DriverInfo) ([]string, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/printer", config.PrinterServiceAddress), nil)
	if err != nil {
		return nil, err
	}
	do, err := p.client.Do(request)
	if err != nil {
		return nil, err
	}
	if do.StatusCode != 200 {
		return nil, apperr.NewClientError(errorlog.StatusCodeErrorMessage, errorlog.InternalServerError, do.StatusCode)
	}

	dto := NewResponseDTO()
	respBody, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, &dto)
	if err != nil {
		return nil, err
	}

	if dto.Error != "null" {
		return nil, errors.New("Ошибка при формировании списка принтеров")
	}
	fmt.Println(dto.Message)

	var array []string

	for _, value := range dto.Message.([]interface{}) {
		array = append(array, value.(string))
	}

	return array, nil
}

func (p *Printer) StopPrinterServer() {
	//TODO implement me
	panic("implement me")
}
