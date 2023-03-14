package middleware

import (
	"github.com/google/logger"
	"log"
	"os"
)

var Directories = [...]string{
	"./log",
	"./debug_info",
	"./debug_info/sell",
	"./debug_info/refound",
	"./debug_info/click",
	"./debug_info/login",
	"./content",
	"./content/system",
	"./content/system/icon",
	"./cookie",
	"./cookie/appconfig",
	"./cookie/click",
	"./cookie/session",
	"./cookie/userdata",
	"./db",
}

var Files = [...]string{
	"./content/system/icon/logo.png",
	"./content/system/icon/main.png",
}

var FilesPathsMap = map[string]string{}

type Middleware struct {
	logf *logger.Logger
	*SystemStatus
}

func NewMiddleware(logg *logger.Logger) *Middleware {
	return &Middleware{logf: logg, SystemStatus: newSystemStatus()}
}

type SystemStatus struct {
	KKTIsActive bool
	FileStatus  map[string]bool
}

func newSystemStatus() *SystemStatus {
	return &SystemStatus{}
}

func (m *Middleware) Initialize() {
	m.CreateAppDirectories()
	m.FileStatus = make(map[string]bool)
	for _, value := range Files {
		m.FileStatus[value] = false
	}
}

func (m *Middleware) CheckAllFiles() {
	for key := range m.FileStatus {
		if _, err := os.Stat(key); err != nil && os.IsNotExist(err) {
			m.FileStatus[key] = false
			m.logf.Warningf("Отсутствует файл %s", key)
		}
	}
}

func (m *Middleware) PullAllNonExistingFiles() {
	for key, value := range m.FileStatus {
		if !value {
			m.logf.Warningf("Происходит загрузка отсутствующего файла - %s", key)
		}
	}
}

func (m *Middleware) CreateAppDirectories() {
	for _, path := range Directories {
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			log.Printf("Создаем папку %s\n", path)
			_ = os.Mkdir(path, 0660)
		}
	}
}

func (m *Middleware) BasicMiddleware() {
	m.Initialize()
	m.CheckAllFiles()
	m.PullAllNonExistingFiles()
}
