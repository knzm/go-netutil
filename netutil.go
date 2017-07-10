package netutil

import (
	"io"
	"net"
	"os"
	"syscall"
)

func IsNetworkErrorFatal(err error) bool {
	if err == nil {
		// no error
		return false
	}

	if err == io.EOF {
		// connection closed
		return false
	}

	opErr, ok := err.(*net.OpError)
	if !ok {
		return true
	}

	if opErr.Timeout() || opErr.Temporary() {
		return false
	}

	syscallErr, ok := opErr.Err.(*os.SyscallError)
	if !ok {
		return true
	}

	errno := syscallErr.Err
	if errno == syscall.EPIPE || errno == syscall.ECONNRESET {
		// connection closed
		return false
	}

	return true
}
