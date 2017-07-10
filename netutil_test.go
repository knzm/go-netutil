package netutil_test

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/knzm/go-netutil"
)

func TestBrokenPipeErrorIsNotFatal(t *testing.T) {
	if testing.Short() {
		t.Skip("skip a slow test using real network connections")
	}

	// A problematic client
	client := func(conn net.Conn) {
		// send a request
		_, err := conn.Write([]byte("ping\n"))
		if err != nil {
			t.Fatal(err)
		}

		// then disconnect without reading any response
		err = conn.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	// open a server side connection
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	addr := listener.Addr()
	t.Logf("address: %s\n", addr)

	go func() {
		// open a client side connection
		conn, err := net.Dial("tcp", addr.String())
		if err != nil {
			t.Fatal(err)
		}
		client(conn)
	}()

	conn, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}

	// read from the socket
	r := bufio.NewReader(conn)
	_, err = r.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	// write to the socket that is already closed for a second repeatedly
	err = func(conn net.Conn) error {
		timer := time.After(1 * time.Second)
		for {
			select {
			case <-timer:
				return nil
			default:
			}

			b := make([]byte, 1)
			_, err = conn.Write(b)
			if err != nil {
				return err
			}
		}
	}(conn)

	if err == nil {
		t.Fatal("An error should be occured.")
	}

	if netutil.IsNetworkErrorFatal(err) {
		t.Fatalf("Fatal error: %s\n", err)
	}

	t.Logf("non fatal error: %v\n", err)

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}
