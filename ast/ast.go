package ast

import "github.com/nanjingblue/maydb/token"

type Ast struct {
	Statements []*Statement
}

type AstKind uint

const (
	SelectKind AstKind = iota
	CreateTableKind
	InsertKind
)

type ExpressionKind uint

const (
	LiteralKind ExpressionKind = iota
)

type Expression struct {
	Literal *token.Token
	Kind    ExpressionKind
}

type Statement struct {
	SelectStatement      *SelectStatement
	CreateTableStatement *CreateTableStatement
	InsertStatement      *InsertStatement
	Kind                 AstKind
}

type InsertStatement struct {
	Table  token.Token
	Values *[]*Expression
}

type ColumnDefinition struct {
	Name     token.Token
	Datatype token.Token
}

type CreateTableStatement struct {
	Name token.Token
	Cols *[]*ColumnDefinition
}

type SelectStatement struct {
	Item []*Expression
	From token.Token
}
