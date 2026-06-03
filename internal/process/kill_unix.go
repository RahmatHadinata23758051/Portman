//go:build !windows

package process

import (
	"os"
	"syscall"
)

func terminateGraceful(proc *os.Process) error {
	return proc.Signal(syscall.SIGTERM)
}
