package mutants

import (
	"fmt"
	"go/ast"
	"os"
)

type Mutator interface {
	Mutate(ast.Node) (func(), bool)
	Name() string
}

func fLog(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}
