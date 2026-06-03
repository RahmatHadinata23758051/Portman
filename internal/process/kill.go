package process

import (
	"fmt"
	"os"
	"strings"
)

// Killer terminates a process by PID.
// Implementations of this interface can be swapped in tests.
type Killer interface {
	Kill(pid int32, force bool) error
}

// SystemKiller implements Killer using the OS process API.
type SystemKiller struct{}

// Kill terminates the process with the given PID.
// When force is false, a graceful termination is attempted first.
// When force is true, the process is killed immediately.
func (SystemKiller) Kill(pid int32, force bool) error {
	if pid <= 0 {
		return fmt.Errorf("invalid PID %d: cannot terminate system or unknown process", pid)
	}
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}
	if force {
		return proc.Kill()
	}
	return terminateGraceful(proc)
}

// systemProcessNames is the set of process names that warrant a warning
// before the user kills them.
var systemProcessNames = map[string]struct{}{
	"systemd":      {},
	"launchd":      {},
	"svchost":      {},
	"postgres":     {},
	"mysql":        {},
	"redis-server": {},
	"docker":       {},
	"dockerd":      {},
	"containerd":   {},
	"sshd":         {},
	"nginx":        {},
	"apache2":      {},
}

// IsSystemProcess reports whether name matches a known long-running system
// process. The check is case-insensitive and ignores a trailing .exe suffix.
func IsSystemProcess(name string) bool {
	normalized := strings.TrimSuffix(strings.ToLower(name), ".exe")
	_, found := systemProcessNames[normalized]
	return found
}
