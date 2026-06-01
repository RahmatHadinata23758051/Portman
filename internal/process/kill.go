package process

import (
	"strings"

	gopsutil "github.com/shirou/gopsutil/v4/process"
)

// Killer terminates a process by PID.
// Implementations of this interface can be swapped in tests.
type Killer interface {
	Kill(pid int32, force bool) error
}

// SystemKiller implements Killer using the OS process API via gopsutil.
// On POSIX: Terminate sends SIGTERM, Kill sends SIGKILL.
// On Windows: both use TerminateProcess (no SIGTERM equivalent exists).
type SystemKiller struct{}

// Kill terminates the process with the given PID.
// When force is false, a graceful termination is attempted first.
// When force is true, the process is killed immediately.
func (SystemKiller) Kill(pid int32, force bool) error {
	p, err := gopsutil.NewProcess(pid)
	if err != nil {
		return err
	}
	if force {
		return p.Kill()
	}
	return p.Terminate()
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
	normalized := strings.ToLower(strings.TrimSuffix(name, ".exe"))
	_, found := systemProcessNames[normalized]
	return found
}
