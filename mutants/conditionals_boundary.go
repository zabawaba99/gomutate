package mutants

import (
	"fmt"
	"go/ast"
	"go/token"
)

var conditionalsBoundaryMapping = map[token.Token]token.Token{
	//    <          <=
	token.LSS: token.LEQ,
	//    <=         <
	token.LEQ: token.LSS,

	//    >          >=
	token.GTR: token.GEQ,
	//    >=         >
	token.GEQ: token.GTR,
}

type ConditionalsBoundary struct {
}

func (nc *ConditionalsBoundary) Name() string {
	return "conditionals_boundary"
}

func (nc *ConditionalsBoundary) Mutate(node ast.Node) (m Mutation, ok bool) {
	switch node.(type) {
	case nil:
		// ignore
	case *ast.BinaryExpr:
		bExpr := node.(*ast.BinaryExpr)

		newOP, ok := conditionalsBoundaryMapping[bExpr.Op]
		if !ok {
			// fmt.Println("Skipping...", bExpr.Op)
			return m, false
		}

		oldOp := bExpr.Op
		bExpr.Op = newOP
		// fmt.Println("Mutated...", oldOp)

		m = Mutation{
			OrgStmt: fmt.Sprintf("%s", oldOp),
			NewStmt: fmt.Sprintf("%s", bExpr.Op),
			Reset:   func() { bExpr.Op = oldOp },
		}
		// setup return
		return m, true
	default:
		// fmt.Printf("ParseFunc: %T\n", node)
	}

	return
}
