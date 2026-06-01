package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill <port|process>",
	Short: "Kill a process by port number or process name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintln(cmd.OutOrStdout(), "kill: not yet implemented")
		return nil
	},
}
