package differ

import (
	"errors"
	"go/ast"
	"go/token"
	"os"
)

type Differ struct {
	leftPath  string
	rightPath string
	pkgs      map[string]*differPackage
}

func NewDiffer(opts ...Option) *Differ {
	d := &Differ{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func WithPackages(left, right string) func(*Differ) {
	return func(d *Differ) {
		d.leftPath = left
		d.rightPath = right
	}
}

func (d *Differ) Diff() error {
	if err := d.validatePackage(d.leftPath); err != nil {
		return err
	}
	return d.initPackages()
}

func (d *Differ) validatePackage(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return PackageNotExist(path)
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *Differ) initPackages() error {
	// for _, path := range d.paths {
	// 	fs := token.NewFileSet()

	// 	pkgs, err := parser.ParseDir(fs, path, nil, 0)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	d.pkgs[path] = &differPackage{
	// 		fs:   fs,
	// 		pkgs: pkgs,
	// 	}
	// }

	return nil
}

type Option func(*Differ)

type differPackage struct {
	fs   *token.FileSet
	pkgs map[string]*ast.Package
}
