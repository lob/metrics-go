package errutil

import (
	"net"
	"os"
	"syscall"

	"github.com/pkg/errors"
)

// IsIgnorableErr returns true if the provided error is a EPIPE error
func IsIgnorableErr(err error) bool {
	e := errors.Cause(err)

	if netErr, ok := e.(*net.OpError); ok {
		if osErr, ok := netErr.Err.(*os.SyscallError); ok {
			return osErr.Err.Error() == syscall.EPIPE.Error() || osErr.Err.Error() == syscall.ECONNRESET.Error()
		}
	}

	return false
}
