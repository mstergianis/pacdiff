package depthparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []parserTestCase{
		{
			name:     "passes through lexer error",
			input:    "xx",
			expected: "",
			err:      "lexer error: unsupported characters in input",
		},
		{
			name:     "raises error when number is not first",
			input:    "s",
			expected: "",
			err:      "parse error: expected a number but got spaceType with value \"s\" at position 0",
		},
		{
			name:     "raises error when no spaceType is given",
			input:    "2",
			expected: "",
			err:      "parse error: expected a spaceType but got EOF with value \"\" at position 1",
		},
		{
			name:     "happy space path",
			input:    "2s",
			expected: "  ",
			err:      "",
		},
		{
			name:     "happy tab path",
			input:    "2t",
			expected: "		",
			err:      "",
		},
		{
			name:     "happy dot path",
			input:    "2d",
			expected: "••",
			err:      "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := Parse(test.input)
			if test.err != "" {
				assert.EqualError(t, err, test.err, "Parse err output did not match expected")
			}
			assert.Equal(t, output, test.expected, "Parse string output did not match expected")
		})
	}
}

type parserTestCase struct {
	name     string
	input    string
	expected string
	err      string
}
