package diff_test

import (
	"os"
	"testing"

	"github.com/mstergianis/pacdiff/pkg/diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiff(t *testing.T) {
	t.Run("bad input", func(t *testing.T) {
		raw := []byte("}{")

		_, err := diff.ParseDiff(raw)
		assert.Error(t, err)
	})

	t.Run("can parse", func(t *testing.T) {
		raw, err := os.ReadFile("fixture")
		require.NoError(t, err)

		d, err := diff.ParseDiff(raw)
		assert.NoError(t, err)

		expected := diff.Diff{
			"type Todo": map[string]any{
				"fields": map[string]any{
					"ExtraField": map[string]any{
						"presentIn": "left",
						"type":      "int",
					},
				},
			},
		}

		require.EqualValues(t, expected, d)
	})
}
