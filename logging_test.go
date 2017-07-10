package netutil_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/knzm/go-netutil"
	"github.com/knzm/go-netutil/netmock"
)

func TestLoggingConn(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		s := "The quick brown fox jumps over the lazy dog"
		conn := &netmock.Conn{
			ReadFunc: func(b []byte) (n int, err error) {
				r := bytes.NewReader([]byte(s))
				return r.Read(b)
			},
		}
		lconn := netutil.LoggingConn{Conn: conn}
		b := make([]byte, 4096)
		n, err := lconn.Read(b)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if n != len(s) {
			t.Errorf("Expected %d, got %d", len(s), n)
		}
	})
	t.Run("Write", func(t *testing.T) {
		conn := &netmock.Conn{
			WriteFunc: func(b []byte) (n int, err error) {
				return len(b), nil
			},
		}
		lconn := netutil.LoggingConn{Conn: conn}
		b := make([]byte, 4096)
		n, err := lconn.Write(b)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if n != len(b) {
			t.Errorf("Expected %d, got %d", len(b), n)
		}
	})
	t.Run("Close", func(t *testing.T) {
		conn := &netmock.Conn{
			CloseFunc: func() error {
				return errors.New("dummy")
			},
		}
		lconn := netutil.LoggingConn{Conn: conn}
		err := lconn.Close()
		if err.Error() != "dummy" {
			t.Errorf("Unexpected error: %s", err)
		}
	})
}
