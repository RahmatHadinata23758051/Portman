package process

import (
	"testing"
)

func TestIsSystemProcess(t *testing.T) {
	tests := []struct {
		name     string
		isSystem bool
	}{
		{"postgres", true},
		{"POSTGRES", true},
		{"postgres.exe", true},
		{"POSTGRES.EXE", true},
		{"node", false},
		{"go", false},
		{"nginx", true},
		{"nginx.exe", true},
		{"nginx-prod", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSystemProcess(tt.name)
			if got != tt.isSystem {
				t.Errorf("IsSystemProcess(%q) = %v; want %v", tt.name, got, tt.isSystem)
			}
		})
	}
}
