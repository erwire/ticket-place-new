package services

import (
	"fmt"
	"github.com/google/logger"
	"os"
	"strings"
	"time"
)

const (
	logExt  = ".log"
	logPath = "./log/"
)

type LoggerService struct {
	file *os.File
	*logger.Logger
	logVerbose bool
}

func NewLogger(logVerbose bool) *LoggerService {
	return &LoggerService{logVerbose: logVerbose}
}

func (l *LoggerService) InitLog() error {
	if _, err := os.Stat(logPath); err != nil && os.IsNotExist(err) {
		err = os.Mkdir(logPath, 0660)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(logPath+time.Now().Format("2006-01-02")+logExt, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	l.file = file
	l.Logger = logger.Init("Logger", l.logVerbose, true, file)
	return nil
}

func (l *LoggerService) Reinit() error {
	info, err := l.file.Stat()
	if err != nil {
		return err
	}
	if info.Name() == logPath+time.Now().Format("2006-01-02")+logExt {
		return fmt.Errorf("уже ведется запись в файл с таким именем")
	}

	if err := l.file.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия файла: %w", err)
	}
	if err := l.InitLog(); err != nil {
		return fmt.Errorf("ошибка связки системы логирования и файла: %w", err)
	}

	return nil
}
func (l *LoggerService) CurrentTime() (time.Time, error) {
	info, err := l.file.Stat()
	if err != nil {
		err = fmt.Errorf("ошибка чтения информации о текущем log-файле: %w", err)
		return time.Time{}, err
	}
	currentTime, err := time.Parse("2006-01-02", strings.ReplaceAll(strings.ReplaceAll(info.Name(), logPath, ""), logExt, ""))
	if err != nil {
		err = fmt.Errorf("ошибка чтения даты из текущего log-файла: %w", err)
		return time.Time{}, err
	}
	return currentTime, nil
}
func (l *LoggerService) InitLogDebugger(duration time.Duration) error {
	file, err := os.OpenFile(logPath+time.Now().Add(duration).Format("2006-01-02")+logExt, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		err = fmt.Errorf("ошибка открытия файла: %w", err)
		return err
	}
	l.file = file
	l.Logger = logger.Init("Logger", l.logVerbose, true, file)
	return nil
}

func (l *LoggerService) ReinitDebugger(duration time.Duration) error {
	info, err := l.file.Stat()
	if err != nil {
		return err
	}
	if info.Name() == logPath+time.Now().Add(duration).Format("2006-01-02")+logExt {
		return fmt.Errorf("уже ведется запись в файл с таким именем")
	}
	l.Logger.Infoln("Начат перенос в следующий файл")
	if err := l.file.Close(); err != nil {
		err = fmt.Errorf("ошибка закрытия файла: %w", err)
		return err
	}
	if err := l.InitLogDebugger(duration); err != nil {
		err := fmt.Errorf("ошибка связки системы логирования и файла: %w", err)
		return err
	}
	l.Logger.Infoln("Продолжение сессии приложения")
	return nil
}

func (l *LoggerService) Close() {
	l.Logger.Close()
}
