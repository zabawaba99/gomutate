package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

var negateConditionalsMapping = map[token.Token]token.Token{
	//    <          >
	token.LSS: token.GTR,
	//    >          <
	token.GTR: token.LSS,
}

type NegateConditionals struct {
}

func (nc *NegateConditionals) Name() string {
	return "boundary_mutations"
}

func (nc *NegateConditionals) Mutate(node ast.Node) (reset func(), ok bool) {
	switch node.(type) {
	case nil:
		// ignore
	case *ast.IfStmt:
		return nc.mutateBoundary(node.(*ast.IfStmt))
	default:
		fmt.Printf("ParseFunc: %T\n", node)
	}

	return
}

func (nc *NegateConditionals) mutateBoundary(ifStmt *ast.IfStmt) (reset func(), ok bool) {
	switch ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		bExpr := ifStmt.Cond.(*ast.BinaryExpr)
		oldOp := bExpr.Op

		bExpr.Op = negateConditionalsMapping[bExpr.Op]
		fmt.Println("Boundary Mutated")

		// setup return
		return func() { bExpr.Op = oldOp }, true
	default:
		fmt.Printf("mutateBoundary: %T\n", ifStmt)
	}

	return
}
