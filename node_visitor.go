package main

import (
	"go/ast"
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
	basedir := filepath.Join(mutationDir, v.mutation.Name(), mutation)
	if err := v.ast.write(basedir); err != nil {
		fLog("Could not create mutation file %s\n", err)
	}
	reset()
	v.ast.mtx.Unlock()

	cmd := exec.Command("go", "test", "."+separator+basedir+separator+"...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	if cmd.ProcessState.Success() {
		dLog("MUTANT IS ALIVE!")
	}

	return v
}
