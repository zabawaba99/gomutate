package main

import (
	"go/ast"
	"go/token"
	"path/filepath"
)

type nodeVistor struct {
	mutation Mutator
	file     *token.File
	ast      *AST
}

func newNodeVisitor(a *AST, file *token.File, m Mutator) *nodeVistor {
	return &nodeVistor{ast: a, file: file, mutation: m}
}

func (v *nodeVistor) Visit(n ast.Node) ast.Visitor {
	v.ast.mtx.Lock()
	reset, ok := v.mutation.Mutate(n)
	if !ok {
		v.ast.mtx.Unlock()
		return v
	}

	mutation := randomStr(10)
	basedir := filepath.Join(mutationDir, v.mutation.Name(), mutation)
	if err := v.ast.write(basedir); err != nil {
		fLog("Could not create mutation file %s\n", err)
	}
	reset()
	v.ast.mtx.Unlock()

	md := MutantData{
		Filename:   v.file.Name(),
		LineNumber: v.file.Line(n.Pos()),
	}
	md.save(basedir)

	return v
}
