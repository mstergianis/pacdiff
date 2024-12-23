package differ

import "fmt"

func PackageNotExist(path string) error {
	return fmt.Errorf("The package %q does not exist and cannot be parsed", path)
}

func NonOneNumberOfPackages(path string, numPackages int) error {
	return fmt.Errorf("The provided path %q had a non-1 number of packages present: %d", path, numPackages)
}
