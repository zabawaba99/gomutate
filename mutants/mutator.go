package mutants

import (
	"fmt"
	"go/ast"
	"os"
)

type Mutation struct {
	OrgStmt string
	NewStmt string
	Reset   func()
}

type Mutator interface {
	Mutate(ast.Node) (Mutation, bool)
	Name() string
}

func fLog(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}
