package cmd

import (
	"github.com/RahmatHadinata23758051/Portman/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var filterCmd = &cobra.Command{
	Use:   "filter <query>",
	Short: "Open TUI with a pre-applied filter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return launchTUI(args[0])
	},
}

func launchTUI(initialFilter string) error {
	p := tea.NewProgram(tui.InitialModel(initialFilter), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
