package output

import (
	"fmt"
	"io"

	"github.com/nata/portman/internal/ports"
)

const (
	colPort    = "PORT"
	colProcess = "PROCESS"
	colPID     = "PID"
	colState   = "STATE"
	colAddress = "ADDRESS"
	colProto   = "PROTOCOL"
)

// PrintTable writes a human-readable aligned table of port entries to w.
func PrintTable(w io.Writer, entries []ports.PortEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No active ports found.")
		return
	}

	widths := columnWidths(entries)
	printHeader(w, widths)
	printRows(w, entries, widths)
}

func columnWidths(entries []ports.PortEntry) [6]int {
	w := [6]int{
		len(colPort),
		len(colProcess),
		len(colPID),
		len(colState),
		len(colAddress),
		len(colProto),
	}
	for _, e := range entries {
		if n := len(fmt.Sprintf("%d", e.Port)); n > w[0] {
			w[0] = n
		}
		if n := len(e.Process); n > w[1] {
			w[1] = n
		}
		if n := len(fmt.Sprintf("%d", e.PID)); n > w[2] {
			w[2] = n
		}
		if n := len(e.State); n > w[3] {
			w[3] = n
		}
		if n := len(e.Address); n > w[4] {
			w[4] = n
		}
		if n := len(e.Protocol); n > w[5] {
			w[5] = n
		}
	}
	return w
}

func printHeader(w io.Writer, widths [6]int) {
	format := rowFormat(widths)
	fmt.Fprintf(w, format, colPort, colProcess, colPID, colState, colAddress, colProto)
}

func printRows(w io.Writer, entries []ports.PortEntry, widths [6]int) {
	format := rowFormat(widths)
	for _, e := range entries {
		fmt.Fprintf(w, format,
			fmt.Sprintf("%d", e.Port),
			e.Process,
			fmt.Sprintf("%d", e.PID),
			e.State,
			e.Address,
			e.Protocol,
		)
	}
}

func rowFormat(w [6]int) string {
	return fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%-%ds  %%-%ds  %%-%ds\n",
		w[0], w[1], w[2], w[3], w[4], w[5])
}
