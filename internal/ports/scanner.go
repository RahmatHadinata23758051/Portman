package ports

import (
	"fmt"
	"sort"
	"strings"

	psnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// ScanOptions controls what the scanner returns.
type ScanOptions struct {
	IncludeEstablished bool
}

// Scan returns active TCP LISTEN ports by default.
// Pass ScanOptions{IncludeEstablished: true} to also include outbound connections.
func Scan(opts ScanOptions) ([]PortEntry, error) {
	conns, err := psnet.Connections("tcp")
	if err != nil {
		return nil, err
	}

	// key: "pid:port:state" — deduplicates same port on multiple interfaces
	seen := make(map[string]bool)
	entries := make([]PortEntry, 0)

	for _, c := range conns {
		if c.Laddr.Port == 0 {
			continue
		}

		state := normalizeState(c.Status)

		if state == "ESTABLISHED" && !opts.IncludeEstablished {
			continue
		}

		if state != "LISTEN" && state != "ESTABLISHED" {
			continue
		}

		// drop ESTABLISHED on loopback — internal IPC, not useful to show
		if state == "ESTABLISHED" && isLoopback(c.Raddr.IP) {
			continue
		}

		key := fmt.Sprintf("%d:%d:%s", c.Pid, c.Laddr.Port, state)
		if seen[key] {
			continue
		}
		seen[key] = true

		entries = append(entries, buildEntry(c, state))
	}

	sort.Slice(entries, func(i, j int) bool {
		pi, pj := statePriority(entries[i].State), statePriority(entries[j].State)
		if pi != pj {
			return pi < pj
		}
		return entries[i].Port < entries[j].Port
	})

	return entries, nil
}

func isLoopback(ip string) bool {
	return strings.HasPrefix(ip, "127.") || ip == "::1"
}

func buildEntry(c psnet.ConnectionStat, state string) PortEntry {
	return PortEntry{
		Port:     c.Laddr.Port,
		Address:  normalizeAddr(c.Laddr.IP),
		Protocol: "tcp",
		PID:      c.Pid,
		Process:  resolveProcessName(c.Pid),
		State:    state,
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
	default:
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
