package ports

// PortEntry holds normalized information about a single local port and the
// process that owns it.
type PortEntry struct {
	Port     uint32 `json:"port"`
	Address  string `json:"address"`
	Protocol string `json:"protocol"`
	PID      int32  `json:"pid"`
	Process  string `json:"process"`
	State    string `json:"state"`
}
