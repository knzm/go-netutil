package httpmock_test

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/knzm/go-netutil/httpmock"
	"github.com/knzm/go-netutil/netmock"
)

func TestMockServerResponse(t *testing.T) {
	type responseAndError struct {
		resp *http.Response
		err  error
	}

	var respValue atomic.Value

	newConn := func(wg *sync.WaitGroup) (net.Conn, error) {
		serverSide, clientSide := net.Pipe()
		go func() {
			defer wg.Done()
			client := &http.Client{
				Transport: &http.Transport{
					Dial: func(network, addr string) (net.Conn, error) {
						return clientSide, nil
					},
				},
			}
			req, _ := http.NewRequest(
				"POST",
				"http://test.server/hello",
				bytes.NewReader([]byte("ping")),
			)
			resp, err := client.Do(req)
			respValue.Store(&responseAndError{resp, err})
		}()
		return serverSide, nil
	}

	h := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("pong"))
	}

	err := httpmock.RunOnce(&httpmock.Server{
		NewConn: newConn,
		Handler: http.HandlerFunc(h),
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	v := respValue.Load()
	if v == nil {
		t.Fatal("No response returned.")
	}

	re := v.(*responseAndError)
	if re.err != nil {
		t.Fatal(re.err)
	}

	if re.resp.StatusCode != 200 {
		t.Fatalf("Expected StatusCode is 200, got %d", re.resp.StatusCode)
	}

	actualBody, err := ioutil.ReadAll(re.resp.Body)
	if err != nil {
		t.Fatalf("Read error %v", err)
	}

	expectedBody := []byte("pong")
	if !bytes.Equal(actualBody, expectedBody) {
		t.Fatalf("Expected body is %q, got %q", expectedBody, actualBody)
	}
}

func TestMockServerOpError(t *testing.T) {
	// In case that a client sends a request and then disconnect
	// without reading any responses.
	newConn := func(wg *sync.WaitGroup) (net.Conn, error) {
		defer wg.Done()

		var buf bytes.Buffer
		buf.WriteString("GET /test HTTP/1.0\n")
		buf.WriteString("\n")
		requestReader := bytes.NewReader(buf.Bytes())

		conn := &netmock.Conn{
			ReadFunc: func(b []byte) (n int, err error) {
				return requestReader.Read(b)
			},
			WriteFunc: func(b []byte) (n int, err error) {
				err = &net.OpError{
					Op:  "write",
					Net: "tcp",
					Err: os.NewSyscallError("write", syscall.EPIPE),
				}
				return 0, err
			},
		}
		return conn, nil
	}

	var writerErr atomic.Value

	h := func(w http.ResponseWriter, req *http.Request) {
		// A response must be greater than the buffer size.
		b := bytes.Repeat([]byte("x"), 4000)
		_, err := w.Write(b)
		if err != nil {
			writerErr.Store(err)
		}
	}

	err := httpmock.RunOnce(&httpmock.Server{
		NewConn: newConn,
		Handler: http.HandlerFunc(h),
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	v := writerErr.Load()
	if v == nil {
		t.Fatal("An error should be occured.")
	}
	if _, ok := v.(*net.OpError); !ok {
		t.Fatal("Unexpected error: %v", v)
	}
}
