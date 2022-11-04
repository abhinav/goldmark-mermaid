package mermaid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderMode_String(t *testing.T) {
	tests := []struct {
		mode RenderMode
		str  string
	}{
		{RenderModeAuto, "Auto"},
		{RenderModeClient, "Client"},
		{RenderModeServer, "Server"},
		{42, "RenderMode(42)"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			assert.Equal(t, tt.str, tt.mode.String())
		})
	}
}
