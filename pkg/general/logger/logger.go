package logger

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

var logger *LogstashHook

const (
	logstashAddr = "localhost:5000"
)

// LogstashHook is a logrus hook that sends logs to a Logstash server over TCP.
type LogstashHook struct {
	conn net.Conn
	addr string
	mu   sync.RWMutex
}

type LogrusLoggers struct {
	Json *logrus.Logger
	Text *logrus.Logger
}

func NewLogrusLoggers() *LogrusLoggers {
	l := &LogrusLoggers{}
	hook, err := NewLogstashHook()
	if err != nil {
		log.Fatal(err)
	}
	l.Json = logrus.New()
	l.Json.AddHook(hook)
	l.Json.SetFormatter(&logrus.JSONFormatter{})

	l.Text = logrus.New()
	l.Text.AddHook(hook)
	l.Text.SetFormatter(&logrus.TextFormatter{})
	return l
}

func NewLogstashHook() (*LogstashHook, error) {
	conn, err := net.Dial("tcp", logstashAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial Logstash: %w", err)
	}
	logger = &LogstashHook{
		conn: conn,
		addr: logstashAddr,
	}
	return logger, nil
}

func GetLogstashHookInstance() (*LogstashHook, error) {
	if logger == nil {
		return nil, fmt.Errorf("need to create a new logstashHook before trying to get an Instance of it")
	}
	return logger, nil
}

// Fire is the logrus hook method. It's called for every log entry.
func (hook *LogstashHook) Fire(entry *logrus.Entry) error {
	// Marshal the log entry to JSON format. logrus.JSONFormatter already
	// structures the data for us.
	serialized, err := entry.Bytes()
	if err != nil {
		return fmt.Errorf("failed to marshal log entry to JSON: %w", err)
	}

	// Write the JSON bytes to the TCP connection.
	if _, err := hook.conn.Write(serialized); err != nil {
		// If the write fails, assume the connection is broken and close it.
		// The main program would need a reconnection strategy for robustness.
		hook.conn.Close()
		return fmt.Errorf("failed to write to Logstash: %w", err)
	}
	return nil
}

// Levels returns the log levels this hook will fire for.
func (hook *LogstashHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Close gracefully closes the connection.
func (hook *LogstashHook) Close() {
	if hook.conn != nil {
		hook.conn.Close()
	}
}
