package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"sync"
)

type Mutator interface {
	Mutate(ast.Node) (func(), bool)
	Name() string
}

type AST struct {
	mtx  sync.Mutex
	fset *token.FileSet
	pkgs map[string]*ast.Package
}

func newAST(filename string) (*AST, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	return &AST{
		fset: fset,
		pkgs: pkgs,
	}, nil
}

func (a *AST) ApplyMutation(m Mutator) {
	// mutation := new(ast.File)
	// *mutation = *a.original

	visitor := newNodeVisitor(a, &BoundaryMutator{})
	a.forEachFile(func(name string, file *ast.File) error {
		ast.Walk(visitor, file)
		return nil
	})
}

func (a *AST) forEachFile(fn func(string, *ast.File) error) error {
	for _, pkg := range a.pkgs {
		for fname, ast := range pkg.Files {
			if err := fn(fname, ast); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AST) write(basepath string) error {
	return a.forEachFile(func(name string, ast *ast.File) error {
		filename := filepath.Join(basepath, name)
		if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
			return err
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		err = printer.Fprint(file, a.fset, ast)
		if err != nil {
			return err
		}

		return nil
	})
}
