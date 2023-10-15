package middleware

import (
	apperr "fptr/internal/error_list"
	errorlog "fptr/pkg/error_logs"
	"github.com/google/logger"
	"io"
	"net/http"
	"os"
)

const DatabasePath = "./db/sqlite.db"

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
	"./content/system/icon/main.ico",
}

var FilesPathsMap = map[string]string{
	"./content/system/icon/logo.png": "https://raw.githubusercontent.com/JahnGeor/ticket-place/main/content/system/icon/logo.png",
	"./content/system/icon/main.png": "https://raw.githubusercontent.com/JahnGeor/ticket-place/main/content/system/icon/main.png",
	"./content/system/icon/main.ico": "https://raw.githubusercontent.com/JahnGeor/ticket-place/main/content/system/icon/main.ico",
}

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
		m.FileStatus[value] = true
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
			err := m.DownloadIcon(FilesPathsMap[key], key)
			if err != nil {
				m.logf.Warningf("Ошибка при загрузке файла %s: %v", key, err)
				continue
			}
			m.logf.Infof("Файл скачался")
		}
	}
}

func (m *Middleware) CreateAppDirectories() {
	for _, path := range Directories {
		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			m.logf.Infof("Создаем папку %s\n", path)
			_ = os.Mkdir(path, 0660)
		}
	}
}

func (m *Middleware) DownloadIcon(reqPath string, path string) error {
	resp, err := http.Get(reqPath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return apperr.NewClientError("Ошибка при скачивании файла", errorlog.ResponseError, resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (m *Middleware) BasicMiddleware() {
	m.Initialize()
	m.CheckAllFiles()
	m.PullAllNonExistingFiles()

}
