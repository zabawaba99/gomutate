package mutants

import (
	"fmt"
	"go/ast"
	"go/token"

	log "github.com/Sirupsen/logrus"
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
	return "negate_conditionals"
}

func (nc *NegateConditionals) Mutate(node ast.Node) (m Mutation, ok bool) {
	switch node.(type) {
	case nil:
		// ignore
	case *ast.BinaryExpr:
		bExpr := node.(*ast.BinaryExpr)
		fields := log.Fields{"name": nc.Name(), "token": bExpr.Op}

		newOP, ok := negateConditionalsMapping[bExpr.Op]
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
		// setup return
		return m, true
	default:
		fields := log.Fields{"name": nc.Name(), "visited": fmt.Sprintf("%T", node)}
		log.WithFields(fields).Debug("found unknown node")
	}

	return
}
