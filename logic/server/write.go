package server

import (
	"net"
)

type Write struct {
	Conn net.Conn
	Body []byte
}

func (w *Write) Out() error {
	_, err := w.Conn.Write(w.Body)
	return err
}
