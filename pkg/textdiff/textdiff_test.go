package textdiff_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mstergianis/pacdiff/pkg/textdiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "./tests"

func TestMyer(t *testing.T) {
	tests, err := os.ReadDir(testDir)
	require.NoError(t, err, "ReadDir should not fail reading the tests directory")

	for _, test := range tests {
		require.True(t, test.IsDir(), "The file: %s, was found not to be a directory. Each test must be a directory", test.Name())
		t.Run(fmt.Sprintf("%s", test.Name()), func(t *testing.T) {
			tc, err := readTest(path.Join(testDir, test.Name()))
			require.NoError(t, err)

			actual, err := textdiff.Myer(tc.left, tc.lName, tc.right, tc.rName)
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
				tc.lName = path.Join(dir, name)
				lRaw, err := os.ReadFile(tc.lName)
				if err != nil {
					return nil, err
				}
				tc.left = string(lRaw)
			}
		case strings.HasPrefix(name, "right"):
			{
				tc.rName = path.Join(dir, name)
				rRaw, err := os.ReadFile(tc.rName)
				if err != nil {
					return nil, err
				}
				tc.right = string(rRaw)
			}
		case strings.HasPrefix(name, "expected"):
			{
				expectedRaw, err := os.ReadFile(path.Join(dir, name))
				if err != nil {
					return nil, err
				}
				tc.expected = string(expectedRaw)
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

	return tc, nil
}

func stringIsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

type testCase struct {
	left     string
	lName    string
	right    string
	rName    string
	expected string
	errorStr *string
}
