package logger

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func TestGetLogstashHook(t *testing.T) {
	hook, err := NewLogstashHook()
	if err != nil {
		t.Fatal(err)
	}
	log = logrus.New()
	log.AddHook(hook)
	// Log some sample data to trigger the hook.
	log.WithFields(logrus.Fields{
		"message": "TestGetLogstashHook successful",
		"user_id": "12345",
		"event":   "login",
	}).Info("Login event")

	time.Sleep(2 * time.Second)
}

func TestGetLogstashHookInstance(t *testing.T) {
	hook, err := GetLogstashHookInstance()
	if err != nil {
		t.Fatal(err)
	}
	log = logrus.New()
	log.AddHook(hook)
	log.SetFormatter(&logrus.JSONFormatter{})
	// Log some sample data to trigger the hook.
	log.WithFields(logrus.Fields{
		"message": "TestGetLogstashHookInstance successful",
		"user_id": "12345",
		"event":   `{"hello": "world","hi":"there"}`,
	}).Info("Login event")

	time.Sleep(2 * time.Second)
}
