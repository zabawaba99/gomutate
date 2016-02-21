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
	case *ast.IfStmt:
		return nc.mutateBoundary(node.(*ast.IfStmt))
	default:
		// fmt.Printf("ParseFunc: %T\n", node)
	}

	return
}

func (nc *NegateConditionals) mutateBoundary(ifStmt *ast.IfStmt) (m Mutation, ok bool) {
	switch ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		bExpr := ifStmt.Cond.(*ast.BinaryExpr)
		oldOp := bExpr.Op

		bExpr.Op = negateConditionalsMapping[bExpr.Op]
		fmt.Println("Boundary Mutated")

		m = Mutation{
			OrgStmt: fmt.Sprintf("%s", oldOp),
			NewStmt: fmt.Sprintf("%s", bExpr.Op),
			Reset:   func() { bExpr.Op = oldOp },
		}
		// setup return
		return m, true
	default:
		fmt.Printf("mutateBoundary: %T\n", ifStmt.Cond)
	}

	return
}
