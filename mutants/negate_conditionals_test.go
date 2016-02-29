package mutants

import (
	"fmt"
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNegateConditionals_Mutate(t *testing.T) {
	tests := []struct {
		oldToken token.Token
		newToken token.Token
		modified bool
	}{
		{ //    ==         !=
			oldToken: token.EQL,
			newToken: token.NEQ,
			modified: true,
		},
		{ //    !=         ==
			oldToken: token.NEQ,
			newToken: token.EQL,
			modified: true,
		},
		{ //    <=          >
			oldToken: token.LEQ,
			newToken: token.GTR,
			modified: true,
		},
		{ //    >=          <
			oldToken: token.GEQ,
			newToken: token.LSS,
			modified: true,
		},
		{ //    <          >=
			oldToken: token.LSS,
			newToken: token.GEQ,
			modified: true,
		},
		{ //    >          <=
			oldToken: token.GTR,
			newToken: token.LEQ,
			modified: true,
		},
		{ //    ||         ||
			oldToken: token.LOR,
			newToken: token.LOR,
			modified: false,
		},
	}

	for _, test := range tests {
		nc := &NegateConditionals{}
		node := &ast.BinaryExpr{
			Op: test.oldToken,
		}

		m, ok := nc.Mutate(node)
		assert.Equal(t, test.modified, ok)
		assert.Equal(t, node.Op, test.newToken)

		if !test.modified {
			continue
		}

		assert.Equal(t, m.OrgStmt, fmt.Sprintf("%s", test.oldToken))
		assert.Equal(t, m.NewStmt, fmt.Sprintf("%s", test.newToken))

		m.Reset()
		assert.Equal(t, node.Op, test.oldToken)
	}
}
