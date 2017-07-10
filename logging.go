package netutil

import (
	"fmt"
	"log"
	"net"
	"runtime"
)

func printCaller(offset int) {
	batchSize := 16
	for {
		callers := make([]uintptr, batchSize)
		n := runtime.Callers(offset, callers)
		if n == 0 {
			break
		}
		for _, pc := range callers[:n] {
			f := runtime.FuncForPC(pc)
			name := f.Name()
			file, line := f.FileLine(pc)
			fmt.Printf("|Caller| %s at %s:%d\n", name, file, line)
		}
		if n < batchSize {
			break
		}
		offset += batchSize
	}
}

func formatBytesWithMaxSize(b []byte, maxSize int) string {
	if len(b) < maxSize {
		return fmt.Sprintf("%q", b)
	}

	return fmt.Sprintf("%q...", b[:maxSize])
}

type LoggingConn struct {
	net.Conn
	PrintCaller bool
}

func (c *LoggingConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	var readBytes []byte
	if err == nil {
		readBytes = b[:n]
	}
	log.Printf("Conn.Read(%s) = (%d, %#v)\n", formatBytesWithMaxSize(readBytes, 1024), n, err)
	if c.PrintCaller {
		printCaller(2)
	}
	return
}

func (c *LoggingConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	var writtenBytes []byte
	if err == nil {
		writtenBytes = b[:n]
	}
	log.Printf("Conn.Write(%s) = (%d, %#v)\n", formatBytesWithMaxSize(writtenBytes, 1024), n, err)
	if c.PrintCaller {
		printCaller(2)
	}
	return
}

func (c *LoggingConn) Close() (err error) {
	err = c.Conn.Close()
	log.Printf("Conn.Close() = %#v\n", err)
	if c.PrintCaller {
		printCaller(2)
	}
	return
}

type LoggingListener struct {
	net.Listener
	PrintCaller bool
}

func (l *LoggingListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &LoggingConn{Conn: conn, PrintCaller: l.PrintCaller}, nil
}
