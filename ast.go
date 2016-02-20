package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sync"
)

type Mutator interface {
	Mutate(ast.Node) (func(), bool)
	Name() string
}

type AST struct {
	mtx      sync.Mutex
	filename string
	fset     *token.FileSet
	pkgs     map[string]*ast.Package
	original *ast.File
}

func newAST(filename string) (*AST, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	// pkgs, err := parser.ParseDir(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	return &AST{
		filename: filename,
		fset:     fset,
		original: file,
	}, nil
}

func (a *AST) ApplyMutation(m Mutator) {
	mutation := new(ast.File)
	*mutation = *a.original

	visitor := newNodeVisitor(a, &BoundaryMutator{})
	ast.Walk(visitor, mutation)
}
