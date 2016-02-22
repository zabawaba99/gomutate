package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zabawaba99/gomutate/mutants"
)

type AST struct {
	mtx   sync.Mutex
	fset  *token.FileSet
	pkgs  map[string]*ast.Package
	files map[string]*token.File
}

func newAST(filename string) (*AST, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	a := &AST{
		fset:  fset,
		pkgs:  pkgs,
		files: map[string]*token.File{},
	}
	fset.Iterate(func(f *token.File) bool {
		a.files[trimWD(f.Name())] = f
		return true
	})

	dLog("Parsed - %s - %#v", filename, a.pkgs)

	return a, nil
}

func (a *AST) forEachFile(fn func(string, *ast.File) error) error {
	for _, pkg := range a.pkgs {
		for fname, ast := range pkg.Files {
			fname = trimWD(fname)
			if err := fn(fname, ast); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AST) ApplyMutation(m mutants.Mutator) {
	a.forEachFile(func(name string, file *ast.File) error {
		if strings.HasSuffix(name, "_test.go") {
			return nil
		}

		visitor := newNodeVisitor(a, a.files[name], m)
		ast.Walk(visitor, file)
		return nil
	})
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
