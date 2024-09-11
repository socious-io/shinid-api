package lib

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerFormatterType string

const (
	LOGGER_JSON_FORMATTER LoggerFormatterType = "LOGGER_JSON_FORMATTER"
	LOGGER_TEXT_FORMATTER LoggerFormatterType = "LOGGER_TEXT_FORMATTER"
)

var (
	jsonFormatter *logrus.JSONFormatter = &logrus.JSONFormatter{}
	textFormat    *logrus.TextFormatter = &logrus.TextFormatter{}
)

// General Logger
type Logger struct {
	lib *logrus.Logger
}

func (logger *Logger) Trace(log string) {
	logger.lib.Trace(log)
}
func (logger *Logger) Debug(log string) {
	logger.lib.Debug(log)
}
func (logger *Logger) Info(log string) {
	logger.lib.Info(log)
}
func (logger *Logger) Warn(log string) {
	logger.lib.Warn(log)
}
func (logger *Logger) Error(log string) {
	logger.lib.Error(log)
}
func (logger *Logger) Fatal(log string) {
	logger.lib.Fatal(log) // Calls os.Exit(1) after logging
}
func (logger *Logger) Panic(log string) {
	logger.lib.Panic(log) // Calls panic() after logging
}

func CreateLogger(output *os.File, formatter LoggerFormatterType) *Logger {
	logger := logrus.New()
	logger.Out = output

	if formatter == LOGGER_JSON_FORMATTER {
		logger.SetFormatter(jsonFormatter)
	}

	l := new(Logger)
	l.lib = logger

	return l
}

// Gin Logger
type GinLogger struct {
	Logger
}
type GinLogFields struct {
	IP             string
	Method         string
	Path           string
	Query          string
	StatusCode     int
	RequestHeaders http.Header
	Headers        http.Header
	RequestBody    *bytes.Buffer
	Body           *bytes.Buffer
	Duration       time.Duration
}

func (ginLogger *GinLogger) Auto(id string, fields GinLogFields) {
	//Automatic Logging
	reqLog := fmt.Sprintf("Request | %s | %s | %s | %s | %s", id, fields.Method, fields.Path, fields.RequestHeaders, fields.RequestBody.String())
	resLog := fmt.Sprintf("Response | %s | %s | %s | %s | %s | %d | %s | %s", id, fields.Method, fields.Path, fields.Headers, fields.StatusCode, fields.Body.String(), fields.Duration)

	ginLogger.Info(reqLog)
	if fields.StatusCode >= 100 {
		ginLogger.Info(resLog)
	}
	if fields.StatusCode >= 400 {
		ginLogger.Warn(resLog)
	}
	if fields.StatusCode >= 500 {
		ginLogger.Error(resLog)
	}
}

func CreateGinLogger(output *os.File, formatter LoggerFormatterType) *GinLogger {
	logger := logrus.New()
	logger.Out = output

	if formatter == LOGGER_JSON_FORMATTER {
		logger.SetFormatter(jsonFormatter)
	}

	l := new(GinLogger)
	l.lib = logger

	return l
}
