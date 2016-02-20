package main

import (
	"fmt"
	"go/ast"
)

type BoundaryMutator struct {
}

func (bm *BoundaryMutator) Name() string {
	return "boundary_mutations"
}

func (bm *BoundaryMutator) Mutate(node ast.Node) (reset func(), ok bool) {
	switch node.(type) {
	case nil:
		// ignore
	case *ast.IfStmt:
		return bm.mutateBoundary(node.(*ast.IfStmt))
	default:
		fmt.Printf("ParseFunc: %T\n", node)
	}

	return
}

func (bm *BoundaryMutator) mutateBoundary(ifStmt *ast.IfStmt) (reset func(), ok bool) {
	switch ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		bExpr := ifStmt.Cond.(*ast.BinaryExpr)
		oldOp := bExpr.Op

		bExpr.Op = mutationMap[bExpr.Op]
		fmt.Println("Boundary Mutated")

		// setup return
		return func() { bExpr.Op = oldOp }, true
	default:
		fmt.Printf("mutateBoundary: %T\n", ifStmt)
	}

	return
}
