package tcpwriter

import (
	"io"
	"net"
	"time"
)

type tcpWriter struct {
	conn net.Conn
}

func (w *tcpWriter) Write(p []byte) (n int, err error) {
	p = append(p, '\n')
	return w.conn.Write(p)
}

func NewTCPWriter(addr string) (io.WriteCloser, error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &tcpWriter{conn: conn}, nil
}

func (w *tcpWriter) Close() error {
	return w.conn.Close()
}
