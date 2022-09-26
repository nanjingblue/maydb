package token

type Location struct {
	Line uint
	Col  uint
}

type Keyword string

const (
	SelectKeyword Keyword = "select"
	FromKeyword   Keyword = "from"
	AsKeyword     Keyword = "as"
	TableKeyword  Keyword = "table"
	CreateKeyword Keyword = "create"
	InsertKeyword Keyword = "insert"
	IntoKeyword   Keyword = "into"
	ValuesKeyword Keyword = "values"
	IntKeyword    Keyword = "int"
	TextKeyword   Keyword = "text"
	WhereKeyword  Keyword = "where"
)

type Symbol string

const (
	SemicolonSymbol  Symbol = ";"
	AsteriskSymbol   Symbol = "*"
	CommaSymbol      Symbol = ","
	LeftParenSymbol  Symbol = "("
	RightParenSymbol Symbol = ")"
)

type TokenKind uint

const (
	KeywordKind TokenKind = iota
	SymbolKind
	IdentifierKind
	StringKind
	NumericKind
)

type Token struct {
	Value string
	Kind  TokenKind
	Loc   Location
}

func (t *Token) Equals(other *Token) bool {
	return t.Value == other.Value && t.Kind == other.Kind
}
