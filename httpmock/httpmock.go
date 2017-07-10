package httpmock

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/knzm/go-netutil"
	"github.com/knzm/go-netutil/netmock"
)

type Trace int

const (
	TraceOff = Trace(0 + iota)
	TraceOn
	TraceAndPrintCaller
)

type Server struct {
	NewConn func(wg *sync.WaitGroup) (net.Conn, error)
	Handler http.Handler
	Trace   Trace
	Timeout time.Duration
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	if timeout == time.Duration(0) {
		wg.Wait()
		return false
	}

	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.Wait()
	}()

	timer := time.After(timeout)
	select {
	case <-c:
		return false
	case <-timer:
		return true
	}
}

func RunOnce(s *Server) (lastError error) {
	type connAndError struct {
		conn net.Conn
		err  error
	}
	connCh := make(chan connAndError)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		conn, err := s.NewConn(&wg)
		connCh <- connAndError{conn, err}
	}()

	var listener net.Listener = &netmock.Listener{
		AcceptFunc: func() (net.Conn, error) {
			select {
			case ce := <-connCh:
				if ce.err != nil {
					return nil, ce.err
				}
				return ce.conn, nil
			}
		},
	}

	if s.Trace >= TraceOn {
		listener = &netutil.LoggingListener{
			Listener:    listener,
			PrintCaller: s.Trace == TraceAndPrintCaller,
		}
	}

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer wg.Done()
			s.Handler.ServeHTTP(w, req)
		}),
	}

	wg.Add(1)
	go func() {
		err := server.Serve(listener)
		if err != nil {
			lastError = err
		}
	}()

	if waitTimeout(&wg, s.Timeout) {
		lastError = errors.New("Timeout")
	} else {
		err := server.Close()
		if err != nil {
			lastError = err
		}
	}

	return
}
