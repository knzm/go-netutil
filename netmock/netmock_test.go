package netmock_test

import (
	"io"
	"testing"

	"github.com/knzm/go-netutil/netmock"
)

func TestMockListener(t *testing.T) {
	listener := netmock.Listener{}
	conn, err := listener.Accept()
	if err != nil {
		t.Fatalf("listener.Accept() failed: %v", err)
	}

	b := make([]byte, 4096)
	if _, err = conn.Read(b); err != io.EOF {
		t.Fatalf("conn.Read() failed: %v", err)
	}
	if _, err = conn.Write(b); err != io.EOF {
		t.Fatalf("conn.Write() failed: %v", err)
	}
	if err = conn.Close(); err != nil {
		t.Fatalf("conn.Close() failed: %v", err)
	}
}
