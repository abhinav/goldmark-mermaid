package mermaid

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtender_rendererAuto(t *testing.T) {
	t.Parallel()

	t.Run("client", func(t *testing.T) {
		t.Parallel()

		ext := Extender{
			execLookPath: func(string) (string, error) {
				return "", errors.New("great sadness")
			},
		}

		mode, r := ext.renderer()
		assert.Equal(t, RenderModeClient, mode)
		assert.IsType(t, new(ClientRenderer), r)
	})

	t.Run("server", func(t *testing.T) {
		t.Parallel()

		ext := Extender{
			execLookPath: func(string) (string, error) {
				return "/path/to/mmdc", nil
			},
		}

		mode, r := ext.renderer()
		assert.Equal(t, RenderModeServer, mode)
		assert.IsType(t, new(ServerRenderer), r)
	})

	t.Run("unknown mode", func(t *testing.T) {
		t.Parallel()

		ext := Extender{RenderMode: 42}
		assert.Panics(t, func() {
			ext.renderer()
		})
	})
}
