package diff_test

import (
	_ "embed"
	"testing"

	"github.com/mstergianis/pacdiff/pkg/diff"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
)

//go:embed expectedHunk1.diff
var expectedHunk1 string

//go:embed expectedHunk2.diff
var expectedHunk2 string

func TestHunk(t *testing.T) {
	testCases := []testCase[diff.Hunk, string]{
		{
			input: diff.Hunk{
				LeftStart:  1,
				LeftEnd:    2,
				RightStart: 1,
				RightEnd:   4,
				Diffs: []diff.Diff{
					{Typ: diff.Equality, Content: "hello"},
					{Typ: diff.Deletion, Content: "world"},
					{Typ: diff.Insertion, Content: "globe"},
					{Typ: diff.Insertion, Content: "carrier"},
					{Typ: diff.Insertion, Content: "atlas"},
				},
			},
			expected: expectedHunk1,
			name:     "test each DiffTyp and different start and end lines",
		},
		{
			input: diff.Hunk{
				LeftStart:  1,
				LeftEnd:    1,
				RightStart: 1,
				RightEnd:   1,
				Diffs: []diff.Diff{
					{Typ: diff.Deletion, Content: "hello"},
					{Typ: diff.Insertion, Content: "world"},
				},
			},
			expected: expectedHunk2,
			name:     "same start and end line",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, tc.expected, tc.input.String())
		})
	}
}

func TestDiff(t *testing.T) {
	testCases := []testCase[diff.Diff, string]{
		{
			input:    diff.Diff{Typ: diff.Insertion, Content: "hello world"},
			expected: "+hello world",
			name:     "insertion",
		},
		{
			input:    diff.Diff{Typ: diff.Deletion, Content: "hello globe"},
			expected: "-hello globe",
			name:     "deletion and different word",
		},
		{
			input:    diff.Diff{Typ: diff.Equality, Content: "hello world"},
			expected: " hello world",
			name:     "equality",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, tc.expected, tc.input.String(), "adds the content to the DiffTyp string")
		})
	}

}

func TestDiffTyp(t *testing.T) {
	testCases := []testCase[diff.DiffTyp, string]{
		{input: diff.Insertion, expected: "+", name: "insertion"},
		{input: diff.Deletion, expected: "-", name: "deletion"},
		{input: diff.Equality, expected: " ", name: "equality"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, tc.expected, tc.input.String())
		})
	}

	t.Run("panic case", func(t *testing.T) {
		var dt diff.DiffTyp = -1
		defer func() {
			if r := recover(); r != nil {
				switch r := r.(type) {
				case string:
					Equal[string](t, "error: encountered an unknown diff.DiffTyp -1", r)
				default:
					t.Fatal("error: TestDiffTyp panic case, did not receive a string")
				}
			}
		}()
		_ = dt.String()
	})
}

type testCase[I, E any] struct {
	input    I
	expected E
	name     string
}

func Equal[T any](t assert.TestingT, expected, actual T, msgAndArgs ...interface{}) bool {
	return assert.Equal(t, expected, actual, msgAndArgs...)
}
