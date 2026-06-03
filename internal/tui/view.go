package tui

import (
	"fmt"
	"strings"

	"github.com/RahmatHadinata23758051/Portman/internal/ports"
	"github.com/charmbracelet/lipgloss"
)

var banner = " _____         _\n" +
	"|  __ \\       | |\n" +
	"| |__) |__  _ | |_  _ __ ___    __ _  _ __\n" +
	"|  ___// _ \\| || __|| '_ ` _ \\  / _` || '_ \\\n" +
	"| |   | (_) || || |_ | | | | | || (_| || | | |\n" +
	`|_|    \___/ |_| \__|_| |_| |_| \__,_||_| |_|`

// ── Palette ──────────────────────────────────────────────────────────────────

var (
	clrPrimary = lipgloss.Color("#7AA2F7") // soft blue
	clrGreen   = lipgloss.Color("#9ECE6A") // green
	clrAmber   = lipgloss.Color("#E0AF68") // amber / warning
	clrMuted   = lipgloss.Color("#565F89") // dim purple-gray
	clrText    = lipgloss.Color("#C0CAF5") // light lavender
	clrBorder  = lipgloss.Color("#414868") // border gray
	clrSelBg   = lipgloss.Color("#283457") // selected row background
	clrError   = lipgloss.Color("#F7768E") // red
)

// ── Styles ───────────────────────────────────────────────────────────────────

var (
	styleTitle = lipgloss.NewStyle().
			Foreground(clrText).
			Bold(true)

	styleVersion = lipgloss.NewStyle().
			Foreground(clrMuted)

	styleTimer = lipgloss.NewStyle().
			Foreground(clrPrimary)

	styleStatBox = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(clrBorder).
			Padding(0, 2)

	styleStatLabel = lipgloss.NewStyle().
			Foreground(clrMuted)

	styleStatValue = lipgloss.NewStyle().
			Foreground(clrText).
			Bold(true)

	styleTableHeader = lipgloss.NewStyle().
				Foreground(clrMuted).
				Bold(true)

	styleRow = lipgloss.NewStyle().
			Foreground(clrText)

	styleSelected = lipgloss.NewStyle().
			Background(clrSelBg).
			Foreground(clrText)

	styleListenState = lipgloss.NewStyle().
				Foreground(clrGreen)

	styleFooterKey = lipgloss.NewStyle().
			Foreground(clrPrimary).
			Bold(true)

	styleFooter = lipgloss.NewStyle().
			Foreground(clrMuted)

	styleFilterPrompt = lipgloss.NewStyle().
				Foreground(clrAmber).
				Bold(true)

	styleConfirm = lipgloss.NewStyle().
			Foreground(clrAmber).
			Bold(true)

	styleError = lipgloss.NewStyle().
			Foreground(clrError)

	styleEmpty = lipgloss.NewStyle().
			Foreground(clrMuted).
			Italic(true)
)

// ── View entry point ─────────────────────────────────────────────────────────

// View renders the complete TUI screen. No business logic runs here.
func (m Model) View() string {
	if m.showHome {
		return m.renderHome()
	}
	if m.loading {
		return "\n  " + styleEmpty.Render("scanning ports...") + "\n"
	}

	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderStats())
	b.WriteString("\n\n")
	b.WriteString(m.renderTable())
	b.WriteString("\n")
	b.WriteString(m.renderFooter())
	return b.String()
}

func (m Model) renderHome() string {
	styledBanner := styleTitle.Render(banner)

	var b strings.Builder
	b.WriteString(styledBanner)
	b.WriteString("\n\n")

	// Description
	b.WriteString("  " + styleConfirm.Render("Stop fighting EADDRINUSE.") + "\n\n")

	// Menu
	b.WriteString("  " + styleFooterKey.Render("[l]") + " " + styleRow.Render("List active ports") + "\n")
	b.WriteString("  " + styleFooterKey.Render("[/]") + " " + styleRow.Render("Filter ports on startup") + "\n")
	b.WriteString("  " + styleFooterKey.Render("[q]") + " " + styleRow.Render("Quit") + "\n\n")

	// Info
	b.WriteString("  " + styleFooter.Render("portman v0.1.0") + "\n")

	return b.String()
}

// ── Sections ─────────────────────────────────────────────────────────────────

func (m Model) renderHeader() string {
	styledBanner := styleTitle.Render(banner)

	right := styleTimer.Render(fmt.Sprintf("refreshing in %ds", m.countdown))
	version := styleVersion.Render("v0.1.0")

	width := m.width
	if width < 40 {
		width = 80
	}

	// version + timer di baris terakhir banner, rata kanan
	bottomLine := version + strings.Repeat(" ",
		max(1, width-lipgloss.Width(version)-lipgloss.Width(right)),
	) + right

	return styledBanner + "\n" + bottomLine
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) renderStats() string {
	total, listening, established := countStates(m.allPorts)

	box := func(label, value string) string {
		return styleStatBox.Render(
			styleStatLabel.Render(label) + "\n" +
				styleStatValue.Render(value),
		)
	}

	b1 := box("active ports", fmt.Sprintf("%d", total))
	b2 := box("listening", fmt.Sprintf("%d", listening))
	b3 := box("established", fmt.Sprintf("%d", established))

	return lipgloss.JoinHorizontal(lipgloss.Top, "  ", b1, "   ", b2, "   ", b3)
}

func (m Model) renderTable() string {
	if m.err != nil {
		return "  " + styleError.Render(fmt.Sprintf("error: %s", m.err)) + "\n"
	}

	if len(m.visible) == 0 {
		if m.filter != "" {
			return "  " + styleEmpty.Render(fmt.Sprintf("no results for %q", m.filter)) + "\n"
		}
		return "  " + styleEmpty.Render("No active ports found.") + "\n"
	}

	portW, procW, pidW, stateW := tableWidths(m.visible)
	format := fmt.Sprintf("  %%-%ds  %%-%ds  %%-%ds  %%s", portW, procW, pidW)

	var b strings.Builder

	b.WriteString(styleTableHeader.Render(
		fmt.Sprintf("  %-*s  %-*s  %-*s  %-*s",
			portW, "PORT", procW, "PROCESS", pidW, "PID", stateW, "STATE"),
	))
	b.WriteString("\n")

	for i, e := range m.visible {
		state := colorState(e.State, stateW)
		row := fmt.Sprintf(format,
			fmt.Sprintf("%d", e.Port),
			e.Process,
			fmt.Sprintf("%d", e.PID),
			state,
		)
		if i == m.cursor {
			b.WriteString(styleSelected.Render(
				fmt.Sprintf("  %-*s  %-*s  %-*d  %-*s",
					portW, fmt.Sprintf("%d", e.Port),
					procW, e.Process,
					pidW, e.PID,
					stateW, e.State,
				),
			))
		} else {
			b.WriteString(styleRow.Render(row))
		}
		b.WriteString("\n")
	}

	if m.confirm == confirmPending && m.cursor < len(m.visible) {
		entry := m.visible[m.cursor]
		b.WriteString("\n")
		b.WriteString("  " + styleConfirm.Render(
			fmt.Sprintf("Kill process %s PID %d using port %d? [y/f/N] ",
				entry.Process, entry.PID, entry.Port),
		))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderFooter() string {
	if m.filterMode {
		return "  " + styleFilterPrompt.Render("/") + " " +
			styleFooter.Render("filter: ") +
			styleFilterPrompt.Render(m.filter+"_")
	}

	type binding struct{ key, desc string }
	bindings := []binding{
		{"↑↓", "navigate"},
		{"k", "kill"},
		{"r", "refresh"},
		{"/", "filter"},
		{"q", "quit"},
	}

	var parts []string
	for _, b := range bindings {
		parts = append(parts,
			styleFooterKey.Render(b.key)+" "+styleFooter.Render(b.desc),
		)
	}
	sep := styleFooter.Render("    ")
	return "  " + strings.Join(parts, sep)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func countStates(entries []ports.PortEntry) (total, listening, established int) {
	total = len(entries)
	for _, e := range entries {
		switch e.State {
		case "LISTEN":
			listening++
		case "ESTABLISHED":
			established++
		}
	}
	return
}

func tableWidths(entries []ports.PortEntry) (portW, procW, pidW, stateW int) {
	portW, procW, pidW, stateW = 4, 7, 3, 5
	for _, e := range entries {
		if n := len(fmt.Sprintf("%d", e.Port)); n > portW {
			portW = n
		}
		if n := len(e.Process); n > procW {
			procW = n
		}
		if n := len(fmt.Sprintf("%d", e.PID)); n > pidW {
			pidW = n
		}
		if n := len(e.State); n > stateW {
			stateW = n
		}
	}
	return
}

func colorState(state string, width int) string {
	padded := fmt.Sprintf("%-*s", width, state)
	if state == "LISTEN" {
		return styleListenState.Render(padded)
	}
	return styleRow.Render(padded)
}
