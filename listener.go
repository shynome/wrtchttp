package wrtchttp

import (
	"net"

	"github.com/shynome/wrtchttp/adapter"
)

type Listener struct {
	Adapter adapter.Adapter
}

var _ net.Listener = &Listener{}

func NewListener() (l *Listener) {
	l = &Listener{}
	return
}

func (l *Listener) Accept() (conn net.Conn, err error) {
	rw, err := l.Adapter.Accept()
	if err != nil {
		return
	}
	return &Conn{Conn: rw, Listener: l}, nil
}

func (l *Listener) Close() error {
	return nil
}

func (l *Listener) Addr() net.Addr {
	return &Addr{}
}
