package differ

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/mstergianis/pacdiff/pkg/diff"
)

type Differ struct {
	leftPath  string
	rightPath string
	pkgs      map[string]*differPackage
}

func NewDiffer(opts ...Option) *Differ {
	d := &Differ{pkgs: map[string]*differPackage{}}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func WithPackages(left, right string) Option {
	return func(d *Differ) {
		d.leftPath = left
		d.rightPath = right
	}
}

func (d *Differ) Diff() (diff.Diff, error) {
	if err := d.validatePackagePath(d.leftPath); err != nil {
		return nil, err
	}

	if err := d.initPackages(); err != nil {
		return nil, err
	}

	if err := d.validatePackage(d.leftPath); err != nil {
		return nil, err
	}
	if err := d.validatePackage(d.rightPath); err != nil {
		return nil, err
	}

	var leftPkg *ast.Package
	for _, v := range d.pkgs[d.leftPath].pkgs {
		leftPkg = v
	}

	var rightPkg *ast.Package
	for _, v := range d.pkgs[d.rightPath].pkgs {
		rightPkg = v
	}

	result := diff.Diff{}
	// compute the diff
	collectScopes(&leftPkg.Scope, leftPkg.Files)
	collectScopes(&rightPkg.Scope, rightPkg.Files)

	for k, leftV := range leftPkg.Scope.Objects {
		// check the diff
		// rightV := rightPkg.Scope.Objects[k]
		switch leftT := leftV.Decl.(type) {
		case *ast.TypeSpec:
			{
				rightScopedObject, presentInRightScope := rightPkg.Scope.Objects[k]
				if !presentInRightScope {
					continue
				}
				rightT := rightScopedObject.Decl.(*ast.TypeSpec)
				switch leftConcrete := leftT.Type.(type) {
				case *ast.StructType:
					{
						rightConcrete := rightT.Type.(*ast.StructType)
						rightM := fieldListToMap(rightConcrete.Fields.List)

						for _, field := range leftConcrete.Fields.List {
							if len(field.Names) == 0 {
								continue
							}
							if _, inRight := rightM[field.Names[0].Name]; !inRight {
								result["type "+leftT.Name.Name] = map[string]any{
									"fields": map[string]any{
										field.Names[0].Name: map[string]any{
											"presentIn": "left",
											"type":      field.Type.(*ast.Ident).Name,
										},
									},
								}
							}
						}
					}
				}
			}
		default:
			panic(fmt.Sprintf("unimplemented: %T", leftT))
		}
	}

	return result, nil
}

func (d *Differ) validatePackage(path string) error {
	if len(d.pkgs[path].pkgs) != 1 {
		return NonOneNumberOfPackages(path, len(d.pkgs[path].pkgs))
	}

	return nil
}

func (d *Differ) validatePackagePath(path string) error {
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
	for _, path := range []string{d.leftPath, d.rightPath} {
		fs := token.NewFileSet()

		pkgs, err := parser.ParseDir(fs, path, nil, 0)
		if err != nil {
			return err
		}

		d.pkgs[path] = &differPackage{
			fs:   fs,
			pkgs: pkgs,
		}
	}

	return nil
}

type Option func(*Differ)

type differPackage struct {
	fs   *token.FileSet
	pkgs map[string]*ast.Package
}

// collectScopes
//
// destructively modifies [packageScope] by collecting all symbols across all [files]
func collectScopes(packageScope **ast.Scope, files map[string]*ast.File) error {
	if *packageScope == nil {
		*packageScope = ast.NewScope(nil)
	}

	for _, file := range files {
		if file == nil {
			continue
		}

		for k, obj := range file.Scope.Objects {
			if _, keyAlreadyPresentInPackage := (*packageScope).Objects[k]; keyAlreadyPresentInPackage {
				return fmt.Errorf("found duplicate symbol in package: %s", k)
			}

			(*packageScope).Objects[k] = obj
		}
		file.Scope.Outer = *packageScope
	}
	return nil
}

func fieldListToMap(l []*ast.Field) map[string]*ast.Field {
	m := make(map[string]*ast.Field, len(l))
	for _, field := range l {
		m[field.Names[0].Name] = field
	}

	return m
}
