package parser

import (
	"fmt"
	"github.com/nanjingblue/maydb/ast"
	"github.com/nanjingblue/maydb/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		source string
		ast    *ast.Ast
	}{
		{
			source: "CREATE TABLE users (id INT, name TEXT);",
			ast: &ast.Ast{
				Statements: []*ast.Statement{
					{
						Kind: ast.CreateTableKind,
						CreateTableStatement: &ast.CreateTableStatement{
							Name: token.Token{
								Loc:   token.Location{Col: 13, Line: 0},
								Kind:  token.IdentifierKind,
								Value: "users",
							},
							Cols: &[]*ast.ColumnDefinition{
								{
									Name: token.Token{
										Loc:   token.Location{Col: 20, Line: 0},
										Kind:  token.IdentifierKind,
										Value: "id",
									},
									Datatype: token.Token{
										Loc:   token.Location{Col: 23, Line: 0},
										Kind:  token.KeywordKind,
										Value: "int",
									},
								},
								{
									Name: token.Token{
										Loc:   token.Location{Col: 28, Line: 0},
										Kind:  token.IdentifierKind,
										Value: "name",
									},
									Datatype: token.Token{
										Loc:   token.Location{Col: 33, Line: 0},
										Kind:  token.KeywordKind,
										Value: "text",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		fmt.Println("(Parser) Testing: ", test.source)
		asts, err := Parse(test.source)
		assert.Nil(t, err, test.source)
		assert.Equal(t, test.ast, asts, test.source)
	}
}
