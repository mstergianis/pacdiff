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
	leftPkg   *differPackage
	rightPkg  *differPackage
}

func NewDiffer(opts ...Option) *Differ {
	d := &Differ{}
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

func (d *Differ) TakeDiff() (diff.GroupedHunksSlice, error) {
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
	for _, v := range d.leftPkg.pkgs {
		leftPkg = v
	}

	var rightPkg *ast.Package
	for _, v := range d.rightPkg.pkgs {
		rightPkg = v
	}

	groupedHunks := diff.GroupedHunksSlice{}
	// compute the diff
	leftScope, err := NewScopeDispatcher(leftPkg, d.leftPkg.fs)
	if err != nil {
		return nil, err
	}
	rightScope, err := NewScopeDispatcher(rightPkg, d.rightPkg.fs)
	if err != nil {
		return nil, err
	}

	if leftScope.pkg.Name != rightScope.pkg.Name {
		leftDir := d.leftPath
		rightDir := d.rightPath

		groupedHunks.Add(leftDir, rightDir, diff.Hunk{
			LeftStart:  1,
			LeftEnd:    1,
			RightStart: 1,
			RightEnd:   1,
			Diffs: []diff.Diff{
				diff.Diff{
					Typ:     diff.Deletion,
					Content: fmtPkg(leftScope.pkg.Name),
				},
				diff.Diff{
					Typ:     diff.Insertion,
					Content: fmtPkg(rightScope.pkg.Name),
				},
			},
		})
	}

	for k, leftScopedObjectFilePair := range leftScope.Objects {
		leftScopedObject := leftScopedObjectFilePair.Object
		// log.Printf("%+v %T\n", leftScopedObject, leftScopedObject)
		// TODO maybe switch on (*Object).Kind
		switch leftT := leftScopedObject.Decl.(type) {
		case *ast.TypeSpec:
			{
				rightScopedObjectFilePair, presentInRightScope := rightScope.Objects[k]
				if !presentInRightScope {
					continue
				}
				rightScopedObject := rightScopedObjectFilePair.Object
				rightT := rightScopedObject.Decl.(*ast.TypeSpec)
				switch leftConcrete := leftT.Type.(type) {
				case *ast.StructType:
					{
						rightConcrete := rightT.Type.(*ast.StructType)

						// leftM := fieldListToMap(leftConcrete.Fields.List)
						rightM := fieldListToMap(rightConcrete.Fields.List)

						removedFields := []*ast.Field{}
						// addedFields := []*ast.Field{}

						for _, field := range leftConcrete.Fields.List {
							if len(field.Names) == 0 {
								continue
							}
							if _, inRight := rightM[field.Names[0].Name]; !inRight {
								removedFields = append(removedFields, field)
								leftFile := leftScope.FileSet.File(leftT.Pos())
								rightFile := rightScope.FileSet.File(rightT.Pos())
								groupedHunks.Add(leftFile.Name(), rightFile.Name(), diff.Hunk{
									LeftStart:  0,
									LeftEnd:    0,
									RightStart: 0,
									RightEnd:   0,
									Diffs:      []diff.Diff{},
								})
							}
						}
					}
				}
			}
		default:
			panic(fmt.Sprintf("unimplemented: %T", leftT))
		}
	}

	return groupedHunks, nil
}

func (d *Differ) validatePackage(path string) error {
	var pkgs *differPackage
	switch path {
	case d.leftPath:
		pkgs = d.leftPkg
	case d.rightPath:
		pkgs = d.rightPkg
	default:
		panic("*Differ.validatePackage: an unknown path was submitted")
	}
	if len(pkgs.pkgs) != 1 {
		return NonOneNumberOfPackages(path, len(pkgs.pkgs))
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
	for _, v := range []struct {
		pkg  **differPackage
		path string
	}{
		{pkg: &d.leftPkg, path: d.leftPath},
		{pkg: &d.rightPkg, path: d.rightPath},
	} {
		fs := token.NewFileSet()

		pkgs, err := parser.ParseDir(fs, v.path, nil, 0)
		if err != nil {
			return err
		}

		*v.pkg = &differPackage{
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

type ObjectFilePair struct {
	File   *ast.File
	Object *ast.Object
}

type ScopeDispatcher struct {
	pkg     *ast.Package
	Objects fileScope
	FileSet *token.FileSet
}

type fileScope map[string]*fileObject
type fileObject struct {
	Object *ast.Object
	File   *ast.File
}

// NewScopeDispatcher
func NewScopeDispatcher(pkg *ast.Package, fset *token.FileSet) (*ScopeDispatcher, error) {
	dispatcher := &ScopeDispatcher{
		pkg:     pkg,
		Objects: fileScope{},
		FileSet: fset,
	}

	for _, file := range dispatcher.pkg.Files {
		for k, obj := range file.Scope.Objects {
			if _, keyAlreadyPresentInPackage := dispatcher.Objects[k]; keyAlreadyPresentInPackage {
				return nil, fmt.Errorf("found duplicate symbol in package: %s", k)
			}
			dispatcher.Objects[k] = &fileObject{obj, file}
		}
	}

	return dispatcher, nil
}

func fieldListToMap(l []*ast.Field) map[string]*ast.Field {
	m := make(map[string]*ast.Field, len(l))
	for _, field := range l {
		m[field.Names[0].Name] = field
	}

	return m
}

func fieldMapKeyUnion(m1, m2 map[string]*ast.Field) map[string]empty {
	union := map[string]empty{}

	for k := range m1 {
		union[k] = empty{}
	}

	for k := range m2 {
		union[k] = empty{}
	}

	return union
}

type empty struct{}

func fmtPkg(pkg string) string {
	return fmt.Sprintf("package %s", pkg)
}
