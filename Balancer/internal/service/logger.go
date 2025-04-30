package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const LOGDIR = "../log"

var (
	AppLogger   *log.Logger
	ErrorLogger *log.Logger
	appLogFile  *os.File
	errLogFile  *os.File
)

func InitLogger() error {
	if err := os.MkdirAll(LOGDIR, 0755); err != nil { // Создание директории для логов. Права drwxr-xr-x
		return fmt.Errorf("Ошибка создания директории логов: %w", err)
	}

	logPath := filepath.Join(LOGDIR, "app.log") // Название файла логов

	var err error

	appLogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) // -rw-r--r--

	if err != nil {
		return fmt.Errorf("Ошибка открытия файла логов: %w", err)
	}

	errorLogPath := filepath.Join(LOGDIR, "error.log")

	errLogFile, err = os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		return fmt.Errorf("Ошибка открытия файла логов: %w", err)
	}

	AppLogger = log.New(appLogFile, "APP: ", log.LstdFlags|log.Llongfile) // Флаги, отвечающие за формат выводов лога и путь, где будет вызываться лог
	ErrorLogger = log.New(errLogFile, "ERROR: ", log.LstdFlags|log.Llongfile)

	return nil
}

func Close() {
	if appLogFile != nil {
		appLogFile.Close()
	}
	if errLogFile != nil {
		errLogFile.Close()
	}
}
