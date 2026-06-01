package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/RahmatHadinata23758051/Portman/internal/ports"
	"github.com/RahmatHadinata23758051/Portman/internal/process"
	tea "github.com/charmbracelet/bubbletea"
)

// Message types used by the Bubble Tea runtime.
type tickMsg struct{}

type scanResultMsg struct {
	entries []ports.PortEntry
	err     error
}

type killResultMsg struct {
	err error
}

// InitialModel returns a ready-to-run Model. Pass a non-empty initialFilter
// to pre-populate the filter (portman filter <query>).
func InitialModel(initialFilter string) Model {
	return Model{
		countdown: refreshInterval,
		filter:    initialFilter,
		loading:   initialFilter != "",
		showHome:  initialFilter == "",
	}
}

// Init starts the first tick and an immediate port scan.
func (m Model) Init() tea.Cmd {
	if m.showHome {
		return nil
	}
	return tea.Batch(doTick(), doScan())
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func doScan() tea.Cmd {
	return func() tea.Msg {
		entries, err := ports.Scan(ports.ScanOptions{})
		return scanResultMsg{entries: entries, err: err}
	}
}

func doKill(entry ports.PortEntry) tea.Cmd {
	return func() tea.Msg {
		k := process.SystemKiller{}
		err := k.Kill(entry.PID, false)
		return killResultMsg{err: err}
	}
}

// Update is the single event dispatcher — no rendering, no direct I/O.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		return m.onTick()
	case scanResultMsg:
		return m.onScan(msg), nil
	case killResultMsg:
		return m, doScan()
	case tea.KeyMsg:
		return m.onKey(msg)
	}
	return m, nil
}

func (m Model) onTick() (Model, tea.Cmd) {
	m.countdown--
	if m.countdown <= 0 {
		m.countdown = refreshInterval
		return m, tea.Batch(doTick(), doScan())
	}
	return m, doTick()
}

func (m Model) onScan(msg scanResultMsg) Model {
	m.loading = false
	if msg.err != nil {
		m.err = msg.err
		return m
	}
	m.err = nil
	m.allPorts = msg.entries
	m.visible = applyFilter(m.allPorts, m.filter)
	if m.cursor >= len(m.visible) && len(m.visible) > 0 {
		m.cursor = len(m.visible) - 1
	}
	return m
}

func (m Model) onKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.confirm == confirmPending {
		return m.onConfirmKey(msg)
	}
	if m.filterMode {
		return m.onFilterKey(msg)
	}
	return m.onNormalKey(msg)
}

func (m Model) onNormalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showHome {
		switch msg.String() {
		case "l", "enter", " ":
			m.showHome = false
			m.loading = true
			m.countdown = refreshInterval
			return m, tea.Batch(doTick(), doScan())
		case "/":
			m.showHome = false
			m.filterMode = true
			m.loading = true
			m.countdown = refreshInterval
			return m, tea.Batch(doTick(), doScan())
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m, tea.Quit
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < len(m.visible)-1 {
			m.cursor++
		}
	case "k":
		if len(m.visible) > 0 {
			m.confirm = confirmPending
		}
	case "r":
		m.countdown = refreshInterval
		return m, doScan()
	case "/":
		m.filterMode = true
	}
	return m, nil
}

func (m Model) onFilterKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.filterMode = false
		m.filter = ""
		m.visible = applyFilter(m.allPorts, m.filter)
		m.cursor = 0
	case "enter":
		m.filterMode = false
	case "backspace":
		if len(m.filter) > 0 {
			m.filter = m.filter[:len(m.filter)-1]
			m.visible = applyFilter(m.allPorts, m.filter)
			m.cursor = 0
		}
	default:
		if len(msg.Runes) == 1 {
			m.filter += string(msg.Runes)
			m.visible = applyFilter(m.allPorts, m.filter)
			m.cursor = 0
		}
	}
	return m, nil
}

func (m Model) onConfirmKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y", "enter":
		if len(m.visible) > 0 && m.cursor < len(m.visible) {
			entry := m.visible[m.cursor]
			m.confirm = confirmNone
			return m, doKill(entry)
		}
		m.confirm = confirmNone
	default:
		m.confirm = confirmNone
	}
	return m, nil
}

// applyFilter returns the subset of entries matching filter against process
// name or port number. Returns the full slice when filter is empty.
func applyFilter(entries []ports.PortEntry, filter string) []ports.PortEntry {
	if filter == "" {
		return entries
	}
	lower := strings.ToLower(filter)
	var result []ports.PortEntry
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Process), lower) ||
			strings.Contains(fmt.Sprintf("%d", e.Port), lower) {
			result = append(result, e)
		}
	}
	return result
}
