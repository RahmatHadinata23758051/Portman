package ports

import (
	"sort"
	"strings"

	psnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// Scan returns all active TCP connections on the local machine.
// LISTEN entries are sorted before ESTABLISHED entries.
// If a process name cannot be resolved, the entry uses "unknown".
func Scan() ([]PortEntry, error) {
	conns, err := psnet.Connections("tcp")
	if err != nil {
		return nil, err
	}

	entries := make([]PortEntry, 0, len(conns))
	for _, c := range conns {
		if c.Laddr.Port == 0 {
			continue
		}
		entry := buildEntry(c)
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return statePriority(entries[i].State) < statePriority(entries[j].State)
	})

	return entries, nil
}

func buildEntry(c psnet.ConnectionStat) PortEntry {
	return PortEntry{
		Port:     c.Laddr.Port,
		Address:  normalizeAddr(c.Laddr.IP),
		Protocol: "tcp",
		PID:      c.Pid,
		Process:  resolveProcessName(c.Pid),
		State:    normalizeState(c.Status),
	}
}

func resolveProcessName(pid int32) string {
	if pid == 0 {
		return "unknown"
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		return "unknown"
	}
	name, err := p.Name()
	if err != nil || name == "" {
		return "unknown"
	}
	return name
}

func normalizeAddr(ip string) string {
	if ip == "" || ip == "::" || ip == "0.0.0.0" {
		return "0.0.0.0"
	}
	return ip
}

func normalizeState(s string) string {
	switch strings.ToUpper(s) {
	case "LISTEN":
		return "LISTEN"
	case "ESTABLISHED":
		return "ESTABLISHED"
	case "TIME_WAIT":
		return "TIME_WAIT"
	case "CLOSE_WAIT":
		return "CLOSE_WAIT"
	default:
		if s == "" {
			return "UNKNOWN"
		}
		return strings.ToUpper(s)
	}
}

func statePriority(state string) int {
	switch state {
	case "LISTEN":
		return 0
	case "ESTABLISHED":
		return 1
	default:
		return 2
	}
}
