package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/RahmatHadinata23758051/Portman/internal/output"
	"github.com/RahmatHadinata23758051/Portman/internal/ports"
	"github.com/spf13/cobra"
)

var (
	listJSON       bool
	allConnections bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List active local ports",
	RunE:  runList,
}

func init() {
	listCmd.Flags().BoolVar(&allConnections, "all", false, "Include established connections")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
}

func runList(cmd *cobra.Command, _ []string) error {
	entries, err := ports.Scan(ports.ScanOptions{IncludeEstablished: allConnections})
	if err != nil {
		if isPermissionError(err) {
			fmt.Fprintln(os.Stderr, "Permission denied while reading process information.")
			fmt.Fprintln(os.Stderr, "Try running portman with elevated privileges.")
			os.Exit(1)
		}
		return err
	}

	if listJSON {
		return output.PrintJSON(cmd.OutOrStdout(), entries)
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
