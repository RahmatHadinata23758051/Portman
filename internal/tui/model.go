package tui

import "github.com/RahmatHadinata23758051/Portman/internal/ports"

const refreshInterval = 2

// confirmState tracks whether an inline kill confirmation is waiting for input.
type confirmState int

const (
	confirmNone    confirmState = iota
	confirmPending confirmState = iota
)

// Model holds all TUI state. No rendering or event logic lives here.
type Model struct {
	allPorts   []ports.PortEntry
	visible    []ports.PortEntry
	cursor     int
	filter     string
	filterMode bool
	countdown  int
	width      int
	height     int
	confirm    confirmState
	loading    bool
	err        error
	showHome   bool
}
