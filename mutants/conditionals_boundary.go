package mutants

import (
	"fmt"
	"go/ast"
	"go/token"

	log "github.com/Sirupsen/logrus"
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
		fields := log.Fields{"name": nc.Name(), "token": bExpr.Op}

		newOP, ok := conditionalsBoundaryMapping[bExpr.Op]
		if !ok {
			log.WithFields(fields).Debug("Skipping...")
			return m, false
		}

		oldOp := bExpr.Op
		bExpr.Op = newOP
		log.WithFields(fields).Debug("Mutated...")

		m = Mutation{
			OrgStmt: fmt.Sprintf("%s", oldOp),
			NewStmt: fmt.Sprintf("%s", bExpr.Op),
			Reset:   func() { bExpr.Op = oldOp },
		}
		return m, true
	default:
		fields := log.Fields{"name": nc.Name(), "visited": fmt.Sprintf("%T", node)}
		log.WithFields(fields).Debug("found unknown node")
	}

	return
}
