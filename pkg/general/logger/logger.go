package logger

import "net"

// LogstashHook is a logrus hook that sends logs to a Logstash server over TCP.
type LogstashHook struct {
	conn net.Conn
	addr string
	// The mutex is not strictly necessary for this simple example as logrus
	// handles concurrent writes to the hook. However, it's good practice
	// for more complex hooks that manage shared resources.
}
