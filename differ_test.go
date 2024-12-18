package differ_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	differ "github.com/mstergianis/pacdiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "tests"

func TestDiffer(t *testing.T) {
	tests, err := os.ReadDir(testDir)
	require.NoError(t, err, "ReadDir should not fail reading the tests directory")

	t.Run("differ with unset packages", func(t *testing.T) {
		err := differ.NewDiffer().Diff()
		assert.Error(t, err, "differ must fail when no packages have been provided")
		assert.Equal(t, err, differ.PackageNotExist(""))
	})

	for _, test := range tests {
		require.True(t, test.IsDir(), "The file: %s, was found not to be a directory. Each test must be a directory", test.Name())
		t.Run(fmt.Sprintf("%s", test.Name()), func(t *testing.T) {
			tc, err := readTest(path.Join(testDir, test.Name()))
			require.NoError(t, err)

			d := differ.NewDiffer(differ.WithPackages(
				tc.left,
				tc.right,
			))

			err = d.Diff()
			require.NoError(t, err)
		})
	}
}

func readTest(dir string) (*testCase, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	tc := &testCase{}
	for _, file := range files {
		name := file.Name()
		switch {
		case strings.HasPrefix(name, "left"):
			{
				tc.left = path.Join(dir, name)
			}
		case strings.HasPrefix(name, "right"):
			{
				tc.right = path.Join(dir, name)
			}
		case strings.HasPrefix(name, "expected"):
			{
				tc.expected = path.Join(dir, name)
			}
		}
	}

	if tc.left == "" || tc.right == "" || tc.expected == "" {
		return nil, fmt.Errorf("missing a field in the test case %q: left %q, right %q, expected %q", dir, tc.left, tc.right, tc.expected)
	}

	return tc, nil
}

type testCase struct {
	left     string
	right    string
	expected string
}
