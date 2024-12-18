package differ

import "fmt"

func PackageNotExist(path string) error {
	return fmt.Errorf("The package %q does not exist and cannot be parsed", path)
}
