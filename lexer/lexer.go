package lexer

import (
	"fmt"
	"github.com/nanjingblue/maydb/token"
	"strings"
)

type Cursor struct {
	pointer uint
	loc     token.Location
}

type lexer func(string, Cursor) (*token.Token, Cursor, bool)

func Lex(source string) ([]*token.Token, error) {
	var tokens []*token.Token
	cur := Cursor{}

lex:
	for cur.pointer < uint(len(source)) {
		lexers := []lexer{lexKeyword, lexSymbol, lexString, lexNumeric, lexIdentifier}
		for _, l := range lexers {
			if tok, newCursor, ok := l(source, cur); ok {
				cur = newCursor

				if tok != nil {
					tokens = append(tokens, tok)
				}
				continue lex
			}
		}
		hint := ""
		if len(tokens) > 0 {
			hint = " after " + tokens[len(tokens)-1].Value + ": '" + string(source[cur.pointer]) + "' "
		}
		for _, t := range tokens {
			fmt.Println(t)
		}
		return nil, fmt.Errorf("unable to lex Token %s, at %d:%d", hint, cur.loc.Line, cur.loc.Col)
	}
	return tokens, nil
}

func lexNumeric(source string, ic Cursor) (*token.Token, Cursor, bool) {
	cur := ic

	periodFound := false
	expMarkerFound := false

	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]
		cur.loc.Col++

		isDigit := c >= '0' && c <= '9'
		isPeriod := c == '.'
		isExpMarker := c == 'e'

		// Must start with a digit or period
		if cur.pointer == ic.pointer {
			if !isDigit && !isPeriod {
				return nil, ic, false
			}

			periodFound = isPeriod
			continue
		}

		if isPeriod {
			if periodFound {
				return nil, ic, false
			}

			periodFound = true
			continue
		}

		if isExpMarker {
			if expMarkerFound {
				return nil, ic, false
			}

			// No periods allowed after expMarker
			periodFound = true
			expMarkerFound = true

			// expMarker must be followed by digits
			if cur.pointer == uint(len(source)-1) {
				return nil, ic, false
			}

			cNext := source[cur.pointer+1]
			if cNext == '-' || cNext == '+' {
				cur.pointer++
				cur.loc.Col++
			}

			continue
		}

		if !isDigit {
			break
		}
	}

	// No characters accumulated
	if cur.pointer == ic.pointer {
		return nil, ic, false
	}

	return &token.Token{
		Value: source[ic.pointer:cur.pointer],
		Loc:   ic.loc,
		Kind:  token.NumericKind,
	}, cur, true
}

func lexCharacterDelimited(source string, ic Cursor, delimiter byte) (*token.Token, Cursor, bool) {
	cur := ic

	if len(source[cur.pointer:]) == 0 {
		return nil, ic, false
	}

	if source[cur.pointer] != delimiter {
		return nil, ic, false
	}

	cur.loc.Col++
	cur.pointer++

	var value []byte
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c := source[cur.pointer]

		if c == delimiter {
			// SQL escapes are via double characters, not backslash.
			if cur.pointer+1 >= uint(len(source)) || source[cur.pointer+1] != delimiter {
				cur.pointer++
				cur.loc.Col++
				return &token.Token{
					Value: string(value),
					Loc:   ic.loc,
					Kind:  token.StringKind,
				}, cur, true
			}
			value = append(value, delimiter)
			cur.pointer++
			cur.loc.Col++
		}

		value = append(value, c)
		cur.loc.Col++
	}

	return nil, ic, false
}

func lexString(source string, ic Cursor) (*token.Token, Cursor, bool) {
	return lexCharacterDelimited(source, ic, '\'')
}

func lexSymbol(source string, ic Cursor) (*token.Token, Cursor, bool) {
	c := source[ic.pointer]
	cur := ic
	// Will get overwritten later if not an ignored syntax
	cur.pointer++
	cur.loc.Col++

	switch c {
	// Syntax that should be thrown away
	case '\n':
		cur.loc.Line++
		cur.loc.Col = 0
		fallthrough
	case '\t':
		fallthrough
	case ' ':
		return nil, cur, true
	}

	// Syntax that should be kept
	symbols := []token.Symbol{
		token.CommaSymbol,
		token.LeftParenSymbol,
		token.RightParenSymbol,
		token.SemicolonSymbol,
		token.AsteriskSymbol,
	}

	var options []string
	for _, s := range symbols {
		options = append(options, string(s))
	}

	// Use `ic`, not `cur`
	match := longestMatch(source, ic, options)
	// Unknown character
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.Col = ic.loc.Col + uint(len(match))

	return &token.Token{
		Value: match,
		Loc:   ic.loc,
		Kind:  token.SymbolKind,
	}, cur, true
}

func lexKeyword(source string, ic Cursor) (*token.Token, Cursor, bool) {
	cur := ic
	keywords := []token.Keyword{
		token.SelectKeyword,
		token.InsertKeyword,
		token.ValuesKeyword,
		token.TableKeyword,
		token.CreateKeyword,
		token.WhereKeyword,
		token.FromKeyword,
		token.IntoKeyword,
		token.IntKeyword,
		token.TextKeyword,
	}

	var options []string
	for _, k := range keywords {
		options = append(options, string(k))
	}

	match := longestMatch(source, ic, options)
	if match == "" {
		return nil, ic, false
	}

	cur.pointer = ic.pointer + uint(len(match))
	cur.loc.Col = ic.loc.Col + uint(len(match))

	return &token.Token{
		Value: match,
		Kind:  token.KeywordKind,
		Loc:   ic.loc,
	}, cur, true
}

// longestMatch iterates through a source string starting at the given
// Cursor to find the longest matching substring among the provided
// options
func longestMatch(source string, ic Cursor, options []string) string {
	var value []byte
	var skipList []int
	var match string

	cur := ic

	for cur.pointer < uint(len(source)) {

		value = append(value, strings.ToLower(string(source[cur.pointer]))...)
		cur.pointer++

	match:
		for i, option := range options {
			for _, skip := range skipList {
				if i == skip {
					continue match
				}
			}

			// Deal with cases like INT vs INTO
			if option == string(value) {
				skipList = append(skipList, i)
				if len(option) > len(match) {
					match = option
				}

				continue
			}

			sharesPrefix := string(value) == option[:cur.pointer-ic.pointer]
			tooLong := len(value) > len(option)
			if tooLong || !sharesPrefix {
				skipList = append(skipList, i)
			}
		}

		if len(skipList) == len(options) {
			break
		}
	}

	return match
}

func lexIdentifier(source string, ic Cursor) (*token.Token, Cursor, bool) {
	// Handle separately if is a double-quoted identifier
	if token, newCursor, ok := lexCharacterDelimited(source, ic, '"'); ok {
		return token, newCursor, true
	}

	cur := ic

	c := source[cur.pointer]
	// Other characters count too, big ignoring non-ascii for now
	isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
	if !isAlphabetical {
		return nil, ic, false
	}
	cur.pointer++
	cur.loc.Col++

	value := []byte{c}
	for ; cur.pointer < uint(len(source)); cur.pointer++ {
		c = source[cur.pointer]

		// Other characters count too, big ignoring non-ascii for now
		isAlphabetical := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isNumeric := c >= '0' && c <= '9'
		if isAlphabetical || isNumeric || c == '$' || c == '_' {
			value = append(value, c)
			cur.loc.Col++
			continue
		}

		break
	}

	if len(value) == 0 {
		return nil, ic, false
	}

	return &token.Token{
		// Unquoted dentifiers are case-insensitive
		Value: strings.ToLower(string(value)),
		Loc:   ic.loc,
		Kind:  token.IdentifierKind,
	}, cur, true
}
