package mutants

import (
	"fmt"
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionalsBoundary_Mutate(t *testing.T) {
	tests := []struct {
		oldToken token.Token
		newToken token.Token
		modified bool
	}{
		{ //    <          <=
			oldToken: token.LSS,
			newToken: token.LEQ,
			modified: true,
		},
		{ //    <=         <
			oldToken: token.LEQ,
			newToken: token.LSS,
			modified: true,
		},
		{ //    >          >=
			oldToken: token.GTR,
			newToken: token.GEQ,
			modified: true,
		},
		{ //    >=         >
			oldToken: token.GEQ,
			newToken: token.GTR,
			modified: true,
		},
		{ //    ||         ||
			oldToken: token.LOR,
			newToken: token.LOR,
			modified: false,
		},
	}

	for _, test := range tests {
		cb := &ConditionalsBoundary{}
		node := &ast.BinaryExpr{
			Op: test.oldToken,
		}

		m, ok := cb.Mutate(node)
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
