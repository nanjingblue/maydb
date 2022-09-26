package backend

import (
	"errors"
	"github.com/nanjingblue/maydb/ast"
)

type ColumnType uint

const (
	TextType ColumnType = iota
	IntType
)

type Cell interface {
	AsText() string
	AsInt() int32
}

type Results struct {
	Columns []struct {
		Type ColumnType
		Name string
	}
	Rows [][]Cell
}

var (
	ErrTableDoesNotExist  = errors.New("table does not exist")
	ErrColumnDoesNotExist = errors.New("column does not exit")
	ErrInvalidDataType    = errors.New("invalid datatype")
	ErrMissingValues      = errors.New("missing value")
)

type Backend interface {
	CreateTable(*ast.CreateTableStatement) error
	Insert(*ast.InsertStatement) error
	Select(*ast.SelectStatement) (*Results, error)
}
