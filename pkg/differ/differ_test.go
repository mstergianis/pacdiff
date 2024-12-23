package differ_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mstergianis/pacdiff/pkg/diff"
	differ "github.com/mstergianis/pacdiff/pkg/differ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "../../test"

func TestDiffer(t *testing.T) {
	tests, err := os.ReadDir(testDir)
	require.NoError(t, err, "ReadDir should not fail reading the tests directory")

	t.Run("differ with unset packages", func(t *testing.T) {
		_, err := differ.NewDiffer().Diff()
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

			actual, err := d.Diff()
			if tc.errorStr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, *tc.errorStr)
			}
			assert.Equal(t, tc.expected, actual)
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
				expectedRaw, err := os.ReadFile(path.Join(dir, name))
				if err != nil {
					return nil, err
				}
				tc.expected, err = diff.ParseDiff(expectedRaw)
				if err != nil {
					return nil, err
				}
			}
		case strings.HasPrefix(name, "error"):
			{
				errFile := path.Join(dir, name)
				errorRaw, err := os.ReadFile(errFile)
				if err == os.ErrNotExist {
					continue
				}
				if err != nil {
					return nil, err
				}
				errMap := map[string]any{}
				err = json.Unmarshal(errorRaw, &errMap)
				if err != nil {
					return nil, err
				}

				const errField = "error"
				rawErrString, errFieldIsPresent := errMap[errField]
				if !errFieldIsPresent {
					return nil, fmt.Errorf("could not find the field %q when parsing the file: %q", errField, errFile)
				}
				errString, ok := rawErrString.(string)
				if !ok {
					return nil, fmt.Errorf("could not cast %q to a string when parsing the file: %q", errField, errFile)
				}
				tc.errorStr = &errString
			}
		}
	}

	if stringIsEmpty(tc.left) || stringIsEmpty(tc.right) {
		return nil, fmt.Errorf("missing a field in the test case %q: left %q, right %q", dir, tc.left, tc.right)
	}

	return tc, nil
}

func stringIsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

type testCase struct {
	left     string
	right    string
	expected diff.Diff
	errorStr *string
}
