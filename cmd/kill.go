package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nata/portman/internal/ports"
	"github.com/nata/portman/internal/process"
	"github.com/spf13/cobra"
)

var (
	killYes   bool
	killForce bool
)

var killCmd = &cobra.Command{
	Use:   "kill <port|process>",
	Short: "Kill a process by port number or process name",
	Args:  cobra.ExactArgs(1),
	RunE:  runKill,
}

func init() {
	killCmd.Flags().BoolVar(&killYes, "yes", false, "Skip confirmation prompt")
	killCmd.Flags().BoolVar(&killForce, "force", false, "Force kill with SIGKILL instead of SIGTERM")
}

func runKill(cmd *cobra.Command, args []string) error {
	port, err := parsePort(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	entries, scanErr := ports.Scan()
	if scanErr != nil {
		return scanErr
	}

	entry, found := findEntryByPort(entries, port)
	if !found {
		fmt.Fprintf(os.Stderr, "No process found using port %d.\n", port)
		os.Exit(1)
	}

	if process.IsSystemProcess(entry.Process) {
		fmt.Fprintf(cmd.OutOrStdout(),
			"Warning: %s may be a long-running system process.\n"+
				"Use --force or confirm carefully if you really want to stop it.\n",
			entry.Process)
	}

	if !killYes {
		confirmed, promptErr := promptConfirm(cmd, entry)
		if promptErr != nil {
			return promptErr
		}
		if !confirmed {
			fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
			return nil
		}
	}

	return executeKill(cmd, entry, killForce)
}

func parsePort(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil || n == 0 || n > 65535 {
		return 0, fmt.Errorf("Invalid port: %s\nPort must be between 1 and 65535.", s)
	}
	return uint32(n), nil
}

func findEntryByPort(entries []ports.PortEntry, port uint32) (ports.PortEntry, bool) {
	for _, e := range entries {
		if e.Port == port {
			return e, true
		}
	}
	return ports.PortEntry{}, false
}

func promptConfirm(cmd *cobra.Command, entry ports.PortEntry) (bool, error) {
	fmt.Fprintf(cmd.OutOrStdout(),
		"Found %s PID %d using port %d.\n",
		entry.Process, entry.PID, entry.Port)
	fmt.Fprint(cmd.OutOrStdout(), "Kill this process? [y/N] ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, err
		}
		return false, nil
	}

	answer := strings.TrimSpace(scanner.Text())
	return answer == "y" || answer == "Y", nil
}

func executeKill(cmd *cobra.Command, entry ports.PortEntry, force bool) error {
	k := process.SystemKiller{}
	if err := k.Kill(entry.PID, force); err != nil {
		if isPermissionError(err) {
			fmt.Fprintln(os.Stderr, "Permission denied. Try running with elevated privileges.")
			os.Exit(1)
		}
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Killed %s (PID %d).\n", entry.Process, entry.PID)
	return nil
}
