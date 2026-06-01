package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nata/portman/internal/ports"
)

// PrintJSON marshals entries as a JSON array and writes it to w.
// Only JSON is written — no surrounding text, no trailing newline beyond
// what json.Marshal produces. Errors are returned to the caller.
func PrintJSON(w io.Writer, entries []ports.PortEntry) error {
	if entries == nil {
		entries = []ports.PortEntry{}
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	_, err = fmt.Fprintf(w, "%s\n", data)
	return err
}
