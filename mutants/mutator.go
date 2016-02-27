package mutants

import "go/ast"

type Mutation struct {
	OrgStmt string
	NewStmt string
	Reset   func()
}

type Mutator interface {
	Mutate(ast.Node) (Mutation, bool)
	Name() string
}
