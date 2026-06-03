//go:build windows

package process

import "os"

func terminateGraceful(proc *os.Process) error {
	return proc.Kill()
}
