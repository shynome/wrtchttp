package wrtchttp

import (
	"net"
	"time"

	"github.com/pion/datachannel"
)

type Conn struct {
	Conn     datachannel.ReadWriteCloser
	Listener net.Listener
}

var _ net.Conn = &Conn{}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.Conn.Read(b)
}
func (c *Conn) Write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}
func (c *Conn) Close() error {
	return c.Conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.Listener.Addr()
}
func (c *Conn) RemoteAddr() net.Addr {
	return &Addr{}
}

func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}

type Addr struct {
	ID string
}

var _ net.Addr = &Addr{}

func (a *Addr) Network() string {
	return "wrtc"
}
func (a *Addr) String() string {
	return a.ID
}
