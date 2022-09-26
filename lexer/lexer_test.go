package lexer

import (
	"github.com/nanjingblue/maydb/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		input  string
		tokens []token.Token
		err    error
	}{
		{input: "CREATE TABLE users (id INT, name TEXT);",
			tokens: []token.Token{
				{
					Loc:   token.Location{Col: 0, Line: 0},
					Value: string(token.CreateKeyword),
					Kind:  token.KeywordKind,
				},
				{
					Loc:   token.Location{Col: 7, Line: 0},
					Value: string(token.TableKeyword),
					Kind:  token.KeywordKind,
				},
				{
					Loc:   token.Location{Col: 13, Line: 0},
					Value: "users",
					Kind:  token.IdentifierKind,
				},
				{
					Loc:   token.Location{Col: 19, Line: 0},
					Value: string(token.LeftParenSymbol),
					Kind:  token.SymbolKind,
				},
				{
					Loc:   token.Location{Col: 20, Line: 0},
					Value: "id",
					Kind:  token.IdentifierKind,
				},
				{
					Loc:   token.Location{Col: 23, Line: 0},
					Value: string(token.IntKeyword),
					Kind:  token.KeywordKind,
				},
				{
					Loc:   token.Location{Col: 26, Line: 0},
					Value: string(token.CommaSymbol),
					Kind:  token.SymbolKind,
				},
				{
					Loc:   token.Location{Col: 28, Line: 0},
					Value: "name",
					Kind:  token.IdentifierKind,
				},
				{
					Loc:   token.Location{Col: 33, Line: 0},
					Value: string(token.TextKeyword),
					Kind:  token.KeywordKind,
				},
				{
					Loc:   token.Location{Col: 37, Line: 0},
					Value: string(token.RightParenSymbol),
					Kind:  token.SymbolKind,
				},
				{
					Loc:   token.Location{Col: 38, Line: 0},
					Value: string(token.SemicolonSymbol),
					Kind:  token.SymbolKind,
				},
			},
			err: nil,
		},
		{
			input: "INSERT INTO users VALUES (1, 'Phil');",
			tokens: []token.Token{
				{
					Loc:   token.Location{Col: 0, Line: 0},
					Value: string(token.InsertKeyword),
					Kind:  token.KeywordKind,
				},
				{
					Loc:   token.Location{Col: 7, Line: 0},
					Value: string(token.IntoKeyword),
					Kind:  token.KeywordKind,
				},
				{
					Loc:   token.Location{Col: 12, Line: 0},
					Value: "users",
					Kind:  token.IdentifierKind,
				},
				{
					Loc:   token.Location{Col: 18, Line: 0},
					Value: string(token.ValuesKeyword),
					Kind:  token.KeywordKind,
				},
				{
					Loc:   token.Location{Col: 25, Line: 0},
					Value: string(token.LeftParenSymbol),
					Kind:  token.SymbolKind,
				},
				{
					Loc:   token.Location{Col: 26, Line: 0},
					Value: "1",
					Kind:  token.NumericKind,
				},
				{
					Loc:   token.Location{Col: 28, Line: 0},
					Value: string(token.CommaSymbol),
					Kind:  token.SymbolKind,
				},
				{
					Loc:   token.Location{Col: 30, Line: 0},
					Value: "Phil",
					Kind:  token.StringKind,
				},
				{
					Loc:   token.Location{Col: 36, Line: 0},
					Value: string(token.RightParenSymbol),
					Kind:  token.SymbolKind,
				},
				{
					Loc:   token.Location{Col: 37, Line: 0},
					Value: string(token.SemicolonSymbol),
					Kind:  token.SymbolKind,
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		tokens, err := Lex(test.input)
		assert.Equal(t, test.err, err, test.input)
		assert.Equal(t, len(test.tokens), len(tokens), test.input)

		for i, tok := range tokens {
			assert.Equal(t, &test.tokens[i], tok, test.input)
		}
	}
}
