package oldlogger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"chaits.org/go-microservices-repo/pkg/network/tcpwriter"
)

var (
	instance *Logger
	mu       sync.RWMutex
)

type Logger struct {
	Service    string
	Env        string
	Context    ContextData
	jsonLogger *slog.Logger
	textLogger *slog.Logger
	tcpWriter  io.WriteCloser
}

type ContextData struct {
	RequestID    string
	TraceID      string
	SpanID       string
	Method       string
	Path         string
	RequestBody  string
	ResponseBody string
	Error        error
	Status       int
	LatencyMs    int64
}

type LogParms struct {
	LogLevel       string
	LogDestination string
}

func NewLogger(service, environment string) *Logger {

	writer, err := tcpwriter.NewTCPWriter("localhost:5000")
	if err != nil {
		slog.Error("Failed to connect to logstash service. Falling back to text logging")
	}

	l := &Logger{
		Service:    service,
		Env:        environment,
		jsonLogger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		textLogger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		tcpWriter:  writer,
	}

	mu.Lock()
	instance = l
	mu.Unlock()

	return l
}

func GetInstance() *Logger {
	mu.RLock()
	defer mu.RUnlock()

	if instance == nil {
		panic("Logger instance not initialized. Call NewLogger() first.")
	}

	return instance
}

func (l *Logger) formatLogEntry(msg, logLevel string) []byte {
	logEntry := LogEntry{
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Level:        logLevel,
		Service:      l.Service,
		Env:          l.Env,
		Message:      msg,
		RequestID:    l.Context.RequestID,
		TraceID:      l.Context.TraceID,
		SpanID:       l.Context.SpanID,
		Method:       l.Context.Method,
		Path:         l.Context.Path,
		Status:       l.Context.Status,
		LatencyMs:    l.Context.LatencyMs,
		RequestBody:  l.Context.RequestBody,
		ResponseBody: l.Context.ResponseBody,
	}

	jsonBytes, err := json.Marshal(logEntry)
	if err != nil {
		log.Fatalf("Failed to marshal log entry: %v", err)
	}

	return jsonBytes
}

func (l *Logger) LogMessage(msg string, logParms LogParms) {
	formattedLog := l.formatLogEntry(msg, logParms.LogLevel)

	if logParms.LogDestination == "" {
		logParms.LogDestination = DEFAULT_LOG_DESTINATION
	}
	if logParms.LogLevel == "" {
		logParms.LogLevel = DEFAULT_LOG_LEVEL
	}

	switch logParms.LogDestination {
	case LOG_TO_CONSOLE_JSON:
		l.jsonLogger.Info(string(formattedLog))
	case LOG_TO_CONSOLE_TEXT:
		l.textLogger.Info(string(formattedLog))
	case LOG_TO_CONSOLE_DIRECT:
		slog.Log(context.Background(), slog.LevelInfo, msg)
	case LOG_TO_LOGSTASH:
		log.Println("Writing to Logstash")
		l.tcpWriter.Write(formattedLog)
	}
}

func (l *Logger) LogHttpMessage(logMessage string, resp *http.Response, contextData ContextData) {
	if resp != nil {
		l.Context.Method = resp.Request.Method
		if resp.Request.Body != nil {
			reqBody, _ := io.ReadAll(resp.Request.Body)
			l.Context.RequestBody = string(reqBody)
		}
		if resp.Body != nil {
			respBody, _ := io.ReadAll(resp.Body)
			l.Context.ResponseBody = string(respBody)
		}
	}
	l.LogMessage(logMessage, LogParms{})
}
