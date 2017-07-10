package netmock

import (
	"io"
	"net"
	"time"
)

type Addr struct {
	NetworkFunc func() string
	StringFunc  func() string
}

func (a *Addr) Network() string {
	if a.NetworkFunc == nil {
		return ""
	}
	return a.NetworkFunc()
}

func (a *Addr) String() string {
	if a.StringFunc == nil {
		return ""
	}
	return a.StringFunc()
}

type Conn struct {
	ReadFunc             func(b []byte) (n int, err error)
	WriteFunc            func(b []byte) (n int, err error)
	CloseFunc            func() error
	LocalAddrFunc        func() net.Addr
	RemoteAddrFunc       func() net.Addr
	SetDeadlineFunc      func(t time.Time) error
	SetReadDeadlineFunc  func(t time.Time) error
	SetWriteDeadlineFunc func(t time.Time) error
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.ReadFunc == nil {
		return 0, io.EOF
	}
	return c.ReadFunc(b)
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.WriteFunc == nil {
		return 0, io.EOF
	}
	return c.WriteFunc(b)
}

func (c *Conn) Close() error {
	if c.CloseFunc == nil {
		return nil
	}
	return c.CloseFunc()
}

func (c *Conn) LocalAddr() net.Addr {
	if c.LocalAddrFunc == nil {
		return &Addr{}
	}
	return c.LocalAddrFunc()
}

func (c *Conn) RemoteAddr() net.Addr {
	if c.RemoteAddrFunc == nil {
		return &Addr{}
	}
	return c.RemoteAddrFunc()
}

func (c *Conn) SetDeadline(t time.Time) error {
	if c.SetDeadlineFunc == nil {
		return nil
	}
	return c.SetDeadlineFunc(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	if c.SetReadDeadlineFunc == nil {
		return nil
	}
	return c.SetReadDeadlineFunc(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	if c.SetWriteDeadlineFunc == nil {
		return nil
	}
	return c.SetWriteDeadlineFunc(t)
}

type Listener struct {
	AcceptFunc func() (net.Conn, error)
	CloseFunc  func() error
	AddrFunc   func() net.Addr
}

func (l *Listener) Accept() (net.Conn, error) {
	if l.AcceptFunc == nil {
		return &Conn{}, nil
	}
	return l.AcceptFunc()
}

func (l *Listener) Close() error {
	if l.CloseFunc == nil {
		return nil
	}
	return l.CloseFunc()
}

func (l *Listener) Addr() net.Addr {
	if l.AddrFunc == nil {
		return &Addr{}
	}
	return l.AddrFunc()
}
