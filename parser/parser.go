package parser

import (
	"errors"
	"fmt"
	"github.com/nanjingblue/maydb/ast"
	"github.com/nanjingblue/maydb/lexer"
	"github.com/nanjingblue/maydb/token"
)

func tokenFromKeyword(k token.Keyword) token.Token {
	return token.Token{
		Kind:  token.KeywordKind,
		Value: string(k),
	}
}

func tokenFromSymbol(s token.Symbol) token.Token {
	return token.Token{
		Kind:  token.SymbolKind,
		Value: string(s),
	}
}

func expectToken(tokens []*token.Token, cursor uint, t token.Token) bool {
	if cursor >= uint(len(tokens)) {
		return false
	}
	return t.Equals(tokens[cursor])
}

func helpMessage(tokens []*token.Token, cursor uint, msg string) {
	var c *token.Token
	if cursor < uint(len(tokens)) {
		c = tokens[cursor]
	} else {
		c = tokens[cursor-1]
	}

	fmt.Printf("[%d,%d]: %s, got: %s\n", c.Loc.Line, c.Loc.Col, msg, c.Value)
}

func Parse(source string) (*ast.Ast, error) {
	tokens, err := lexer.Lex(source)
	if err != nil {
		return nil, err
	}

	a := ast.Ast{}
	cursor := uint(0)
	for cursor < uint(len(tokens)) {
		stmt, newCursor, ok := parseStatement(tokens, cursor, tokenFromSymbol(token.SemicolonSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected statement")
			return nil, errors.New("failed to parse, expected statement")
		}
		cursor = newCursor

		a.Statements = append(a.Statements, stmt)

		atLeastOneSemicolon := false
		for expectToken(tokens, cursor, tokenFromSymbol(token.SemicolonSymbol)) {
			cursor++
			atLeastOneSemicolon = true
		}

		if !atLeastOneSemicolon {
			helpMessage(tokens, cursor, "Expected semi-colon delimiter between statements")
			return nil, errors.New("missing semi-colon between statements")
		}
	}

	return &a, nil
}

func parseStatement(tokens []*token.Token, initialCursor uint, delimiter token.Token) (*ast.Statement, uint, bool) {
	cursor := initialCursor

	// Look for a SELECT statement
	semicolonToken := tokenFromSymbol(token.SemicolonSymbol)
	slct, newCursor, ok := parseSelectStatement(tokens, cursor, semicolonToken)
	if ok {
		return &ast.Statement{
			Kind:            ast.SelectKind,
			SelectStatement: slct,
		}, newCursor, true
	}

	// Look for a INSERT statement
	inst, newCursor, ok := parseInsertStatement(tokens, cursor, semicolonToken)
	if ok {
		return &ast.Statement{
			Kind:            ast.InsertKind,
			InsertStatement: inst,
		}, newCursor, true
	}

	// Look for a CREATE statement
	crtTbl, newCursor, ok := parseCreateTableStatement(tokens, cursor, semicolonToken)
	if ok {
		return &ast.Statement{
			Kind:                 ast.CreateTableKind,
			CreateTableStatement: crtTbl,
		}, newCursor, true
	}

	return nil, initialCursor, false
}

func parseSelectStatement(tokens []*token.Token, initialCursor uint, delimiter token.Token) (*ast.SelectStatement, uint, bool) {
	cursor := initialCursor
	if !expectToken(tokens, cursor, tokenFromKeyword(token.SelectKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	slct := ast.SelectStatement{}

	exps, newCursor, ok := parseExpressions(tokens, cursor, []token.Token{tokenFromKeyword(token.FromKeyword), delimiter})
	if !ok {
		return nil, initialCursor, false
	}

	slct.Item = *exps
	cursor = newCursor

	if expectToken(tokens, cursor, tokenFromKeyword(token.FromKeyword)) {
		cursor++

		from, newCursor, ok := parseToken(tokens, cursor, token.IdentifierKind)
		if !ok {
			helpMessage(tokens, cursor, "Expected FROM token")
			return nil, initialCursor, false
		}

		slct.From = *from
		cursor = newCursor
	}

	return &slct, cursor, true
}

func parseToken(tokens []*token.Token, initialCursor uint, kind token.TokenKind) (*token.Token, uint, bool) {
	cursor := initialCursor

	if cursor >= uint(len(tokens)) {
		return nil, initialCursor, false
	}

	current := tokens[cursor]
	if current.Kind == kind {
		return current, cursor + 1, true
	}

	return nil, initialCursor, false
}

func parseExpressions(tokens []*token.Token, initialCursor uint, delimiters []token.Token) (*[]*ast.Expression, uint, bool) {
	cursor := initialCursor

	var exps []*ast.Expression
outer:
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		// Look for delimiter
		current := tokens[cursor]
		for _, delimiter := range delimiters {
			if delimiter.Equals(current) {
				break outer
			}
		}

		// Look for comma
		if len(exps) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(token.CommaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}

			cursor++
		}

		// Look for expression
		exp, newCursor, ok := parseExpression(tokens, cursor, tokenFromSymbol(token.CommaSymbol))
		if !ok {
			helpMessage(tokens, cursor, "Expected expression")
			return nil, initialCursor, false
		}
		cursor = newCursor

		exps = append(exps, exp)
	}

	return &exps, cursor, true
}

func parseExpression(tokens []*token.Token, initialCursor uint, _ token.Token) (*ast.Expression, uint, bool) {
	cursor := initialCursor

	kinds := []token.TokenKind{token.IdentifierKind, token.NumericKind, token.StringKind}
	for _, kind := range kinds {
		t, newCursor, ok := parseToken(tokens, cursor, kind)
		if ok {
			return &ast.Expression{
				Literal: t,
				Kind:    ast.LiteralKind,
			}, newCursor, true
		}
	}

	return nil, initialCursor, false
}

func parseInsertStatement(tokens []*token.Token, initialCursor uint, delimiter token.Token) (*ast.InsertStatement, uint, bool) {
	cursor := initialCursor

	// Look for INSERT
	if !expectToken(tokens, cursor, tokenFromKeyword(token.InsertKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	// Look for INTO
	if !expectToken(tokens, cursor, tokenFromKeyword(token.IntoKeyword)) {
		helpMessage(tokens, cursor, "Expected into")
		return nil, initialCursor, false
	}
	cursor++

	// Look for table name
	table, newCursor, ok := parseToken(tokens, cursor, token.IdentifierKind)
	if !ok {
		helpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for VALUES
	if !expectToken(tokens, cursor, tokenFromKeyword(token.ValuesKeyword)) {
		helpMessage(tokens, cursor, "Expected VALUES")
		return nil, initialCursor, false
	}
	cursor++

	// Look for left paren
	if !expectToken(tokens, cursor, tokenFromSymbol(token.LeftParenSymbol)) {
		helpMessage(tokens, cursor, "Expected left paren")
		return nil, initialCursor, false
	}
	cursor++

	// Look for expression list
	values, newCursor, ok := parseExpressions(tokens, cursor, []token.Token{tokenFromSymbol(token.RightParenSymbol)})
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	// Look for right paren
	if !expectToken(tokens, cursor, tokenFromSymbol(token.RightParenSymbol)) {
		helpMessage(tokens, cursor, "Expected right paren")
		return nil, initialCursor, false
	}
	cursor++

	return &ast.InsertStatement{
		Table:  *table,
		Values: values,
	}, cursor, true
}

func parseCreateTableStatement(tokens []*token.Token, initialCursor uint, delimiter token.Token) (*ast.CreateTableStatement, uint, bool) {
	cursor := initialCursor

	if !expectToken(tokens, cursor, tokenFromKeyword(token.CreateKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	if !expectToken(tokens, cursor, tokenFromKeyword(token.TableKeyword)) {
		return nil, initialCursor, false
	}
	cursor++

	name, newCursor, ok := parseToken(tokens, cursor, token.IdentifierKind)
	if !ok {
		helpMessage(tokens, cursor, "Expected table name")
		return nil, initialCursor, false
	}
	cursor = newCursor
	if !expectToken(tokens, cursor, tokenFromSymbol(token.LeftParenSymbol)) {
		helpMessage(tokens, cursor, "Expected left parenthesis")
		return nil, initialCursor, false
	}
	cursor++

	cols, newCursor, ok := parseColumnDefinitions(tokens, cursor, tokenFromSymbol(token.RightParenSymbol))
	if !ok {
		return nil, initialCursor, false
	}
	cursor = newCursor

	if !expectToken(tokens, cursor, tokenFromSymbol(token.RightParenSymbol)) {
		helpMessage(tokens, cursor, "Expected right parenthesis")
		return nil, initialCursor, false
	}
	cursor++
	return &ast.CreateTableStatement{
		Name: *name,
		Cols: cols,
	}, cursor, true
}

func parseColumnDefinitions(tokens []*token.Token, initialCursor uint, delimiter token.Token) (*[]*ast.ColumnDefinition, uint, bool) {
	cursor := initialCursor

	var cds []*ast.ColumnDefinition
	for {
		if cursor >= uint(len(tokens)) {
			return nil, initialCursor, false
		}

		current := tokens[cursor]
		if delimiter.Equals(current) {
			break
		}
		if len(cds) > 0 {
			if !expectToken(tokens, cursor, tokenFromSymbol(token.CommaSymbol)) {
				helpMessage(tokens, cursor, "Expected comma")
				return nil, initialCursor, false
			}
			cursor++
		}
		id, newCursor, ok := parseToken(tokens, cursor, token.IdentifierKind)
		if !ok {
			helpMessage(tokens, cursor, "Expected column name")
			return nil, initialCursor, false
		}
		cursor = newCursor

		ty, newCursor, ok := parseToken(tokens, cursor, token.KeywordKind)
		if !ok {
			helpMessage(tokens, cursor, "Expected column type")
			return nil, initialCursor, false
		}
		cursor = newCursor

		cds = append(cds, &ast.ColumnDefinition{
			Name:     *id,
			Datatype: *ty,
		})
	}
	return &cds, cursor, true
}
