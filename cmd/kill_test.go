package cmd

import (
	"testing"
)

func TestClassifyArg(t *testing.T) {
	tests := []struct {
		input  string
		isPort bool
		port   uint32
	}{
		{"80", true, 80},
		{"3000", true, 3000},
		{"65535", true, 65535},
		{"0", false, 0},
		{"65536", false, 0},
		{"-1", false, 0},
		{"node", false, 0},
		{"8080a", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			isPort, port := classifyArg(tt.input)
			if isPort != tt.isPort || port != tt.port {
				t.Errorf("classifyArg(%q) = (%v, %d); want (%v, %d)", tt.input, isPort, port, tt.isPort, tt.port)
			}
		})
	}
}
