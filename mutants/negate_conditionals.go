package mutants

import (
	"fmt"
	"go/ast"
	"go/token"
)

var negateConditionalsMapping = map[token.Token]token.Token{
	//    ==         !=
	token.EQL: token.NEQ,
	//    !=         ==
	token.NEQ: token.EQL,

	//    <=         >
	token.LEQ: token.GTR,
	//    >=         <
	token.GEQ: token.LSS,

	//    <          >=
	token.LSS: token.GEQ,
	//    >          <=
	token.GTR: token.LEQ,
}

type NegateConditionals struct {
}

func (nc *NegateConditionals) Name() string {
	return "boundary_mutations"
}

func (nc *NegateConditionals) Mutate(node ast.Node) (m Mutation, ok bool) {
	switch node.(type) {
	case nil:
		// ignore
	case *ast.BinaryExpr:
		bExpr := node.(*ast.BinaryExpr)

		newOP, ok := negateConditionalsMapping[bExpr.Op]
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
