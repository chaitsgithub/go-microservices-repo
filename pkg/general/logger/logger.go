package logger

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	// Import the go-logstash hook.
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
)

var Logger *logrus.Logger

const (
	logstashAddr = "localhost:5000"
)

func Init(serviceName string) {
	if Logger != nil {
		return
	}
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: false, // Set to true for local development readability, false for production.
	})

	// Common levels are Debug, Info, Warn, Error, and Fatal.
	Logger.SetLevel(logrus.InfoLevel)
	conn, err := net.Dial("tcp", logstashAddr)
	if err != nil {
		Logger.WithError(err).Errorln("Error connecting to logstash")
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "go-microservices-repo", "service": serviceName}))
	Logger.AddHook(hook)
}
