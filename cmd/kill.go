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
	input := args[0]
	isPort, port := classifyArg(input)

	entries, scanErr := ports.Scan()
	if scanErr != nil {
		return scanErr
	}

	if isPort {
		entry, found := findEntryByPort(entries, port)
		if !found {
			fmt.Fprintf(os.Stderr, "No process found using port %d.\n", port)
			os.Exit(1)
		}

		warnSystemProcesses(cmd, []ports.PortEntry{entry})

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

	return killByName(cmd, entries, input)
}

func classifyArg(s string) (bool, uint32) {
	n, err := strconv.ParseUint(s, 10, 32)
	if err == nil && n > 0 && n <= 65535 {
		return true, uint32(n)
	}
	return false, 0
}

func killByName(cmd *cobra.Command, entries []ports.PortEntry, name string) error {
	matches := findEntriesByName(entries, name)
	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "No process found with name %s.\n", name)
		os.Exit(1)
	}

	warnSystemProcesses(cmd, matches)

	if !killYes {
		confirmed, err := promptConfirmMany(cmd, matches)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
			return nil
		}
	}

	for _, entry := range matches {
		if err := executeKill(cmd, entry, killForce); err != nil {
			return err
		}
	}
	return nil
}

func findEntryByPort(entries []ports.PortEntry, port uint32) (ports.PortEntry, bool) {
	for _, e := range entries {
		if e.Port == port {
			return e, true
		}
	}
	return ports.PortEntry{}, false
}

func findEntriesByName(entries []ports.PortEntry, name string) []ports.PortEntry {
	var matches []ports.PortEntry
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Process), strings.ToLower(name)) {
			matches = append(matches, e)
		}
	}
	return matches
}

func warnSystemProcesses(cmd *cobra.Command, entries []ports.PortEntry) {
	for _, entry := range entries {
		if process.IsSystemProcess(entry.Process) {
			fmt.Fprintf(cmd.OutOrStdout(),
				"Warning: %s (PID %d) may be a long-running system process.\n",
				entry.Process, entry.PID)
		}
	}
}

func promptConfirm(cmd *cobra.Command, entry ports.PortEntry) (bool, error) {
	fmt.Fprintf(cmd.OutOrStdout(),
		"Found %s PID %d using port %d.\n",
		entry.Process, entry.PID, entry.Port)
	return askYesNo(cmd, "Kill this process? [y/N] ")
}

func promptConfirmMany(cmd *cobra.Command, entries []ports.PortEntry) (bool, error) {
	fmt.Fprintln(cmd.OutOrStdout(), "Found the following processes:")
	for _, e := range entries {
		fmt.Fprintf(cmd.OutOrStdout(), "- %s (PID: %d, Port: %d)\n", e.Process, e.PID, e.Port)
	}
	return askYesNo(cmd, "Kill these processes? [y/N] ")
}

func askYesNo(cmd *cobra.Command, prompt string) (bool, error) {
	fmt.Fprint(cmd.OutOrStdout(), prompt)
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
