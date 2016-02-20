package main

import (
	"go/ast"
	"go/printer"
	"os"
	"os/exec"
	"path/filepath"
)

type nodeVistor struct {
	mutation Mutator
	ast      *AST
}

func newNodeVisitor(a *AST, m Mutator) *nodeVistor {
	return &nodeVistor{ast: a, mutation: m}
}

func (v *nodeVistor) Visit(n ast.Node) ast.Visitor {
	v.ast.mtx.Lock()
	reset, ok := v.mutation.Mutate(n)
	if !ok {
		v.ast.mtx.Unlock()
		return v
	}

	mutation := randomStr(10)
	filename := filepath.Join(mutationDir, v.mutation.Name(), mutation, v.ast.filename)
	if err := v.writeFile(filename); err != nil {
		fLog("Could not create mutation file %s\n", err)
	}
	reset()
	v.ast.mtx.Unlock()

	cmd := exec.Command("go", "test", "./"+filepath.Dir(filename)+"/...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fLog("Could not run tests %s\n", err)
	}

	return v
}

func (v *nodeVistor) writeFile(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = printer.Fprint(file, v.ast.fset, v.ast.original)
	if err != nil {
		return err
	}

	return nil
}
