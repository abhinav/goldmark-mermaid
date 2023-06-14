package mermaid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderMode_String(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.str, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.str, tt.mode.String())
		})
	}
}
