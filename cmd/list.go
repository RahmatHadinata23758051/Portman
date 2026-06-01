package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nata/portman/internal/output"
	"github.com/nata/portman/internal/ports"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List active local ports",
	RunE:  runList,
}

func runList(cmd *cobra.Command, _ []string) error {
	entries, err := ports.Scan()
	if err != nil {
		if isPermissionError(err) {
			fmt.Fprintln(os.Stderr, "Permission denied while reading process information.")
			fmt.Fprintln(os.Stderr, "Try running portman with elevated privileges.")
			os.Exit(1)
		}
		return err
	}

	output.PrintTable(cmd.OutOrStdout(), entries)
	return nil
}

func isPermissionError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "permission denied") ||
		strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "operation not permitted")
}
