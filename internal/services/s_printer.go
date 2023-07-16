package services

import (
	"fptr/internal/entities"
	"fptr/internal/gateways"
	"github.com/google/logger"
)

type PrinterInterface interface {
	Print(config entities.DriverInfo, dto entities.OrderDTO, pp entities.PageParamsDTO) error
	Ping(config entities.DriverInfo) error
	GetListOfPrinters(config entities.DriverInfo) ([]string, error)
	StartService(config entities.DriverInfo) error
	//StopService(config entities.DriverInfo)
}
type PrinterService struct {
	gw gateways.PrinterInterface
	*logger.Logger
}

func NewPrinterService(logg *logger.Logger) *PrinterService {
	return &PrinterService{
		gw:     gateways.NewPrinter(),
		Logger: logg,
	}
}

func (ps *PrinterService) StartService(config entities.DriverInfo) error {
	if err := ps.gw.StartService(config); err != nil {
		ps.Errorf("%s: %v", "Ошибка при запуске службы печати", err)
		return err
	}
	return nil
}

func (ps *PrinterService) Print(config entities.DriverInfo, dto entities.OrderDTO, pp entities.PageParamsDTO) error {
	if err := ps.gw.Print(config, dto, pp); err != nil {
		ps.Errorf("%s: %v", "Ошибка при печати билета", err)
		return err
	}
	return nil
}

func (ps *PrinterService) Ping(config entities.DriverInfo) error {
	if err := ps.gw.Ping(config); err != nil {
		ps.Errorf("%s: %v", "Ошибка при проверке работы службы печати", err)
		return err
	}
	return nil
}

func (ps *PrinterService) GetListOfPrinters(config entities.DriverInfo) ([]string, error) {
	printers, err := ps.gw.GetPrinterList(config)
	if err != nil {
		ps.Errorf("%s: %v", "Ошибка при проверке работы службы печати", err)
		return nil, err
	}
	return printers, nil
}
