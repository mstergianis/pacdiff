package depthparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tests := []lexerTestCase{
		{
			name:  "only_number",
			input: "44",
			expected: []token{
				{pos: 0, typ: tokenNumber, value: "44"},
				{typ: tokenEOF, pos: 2},
			},
		},
		{
			name:  "basic",
			input: "4s",
			expected: []token{
				{typ: tokenNumber, pos: 0, value: "4"},
				{typ: tokenSpaceType, pos: 1, value: "s"},
				{pos: 2, typ: tokenEOF},
			},
		},
		{
			name:  "error",
			input: "xx",
			expected: []token{
				{typ: tokenError, pos: 0, value: "unsupported characters in input"},
				{pos: 0, typ: tokenEOF},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.input)
			tokens := []token{}
			for {
				token := l.lex()
				tokens = append(tokens, token)
				if token.typ == tokenEOF {
					break
				}
			}

			assert.Equal(t, test.expected, tokens, "tokens must match expeted")
		})
	}
}

type lexerTestCase struct {
	name     string
	input    string
	expected []token
}
