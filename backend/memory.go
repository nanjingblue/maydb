package backend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/nanjingblue/maydb/ast"
	"github.com/nanjingblue/maydb/token"
	"strconv"
)

type MemoryCell []byte

func (mc MemoryCell) AsInt() int32 {
	var i int32
	err := binary.Read(bytes.NewBuffer(mc), binary.BigEndian, &i)
	if err != nil {
		panic(err)
	}
	return i
}

func (mc MemoryCell) AsText() string {
	return string(mc)
}

type Table struct {
	Columns     []string
	ColumnTypes []ColumnType
	Rows        [][]MemoryCell
}

type MemoryBackend struct {
	Tables map[string]*Table
}

func NewMemoryBacked() *MemoryBackend {
	return &MemoryBackend{
		Tables: map[string]*Table{},
	}
}

// CreateTable 创建表
func (mb *MemoryBackend) CreateTable(crt *ast.CreateTableStatement) error {
	t := Table{}
	mb.Tables[crt.Name.Value] = &t
	if crt.Cols == nil {
		return nil
	}
	for _, col := range *crt.Cols {
		t.Columns = append(t.Columns, col.Name.Value)

		var dt ColumnType
		switch col.Datatype.Value {
		case "int":
			dt = IntType
		case "text":
			dt = TextType
		default:
			return ErrInvalidDataType
		}
		t.ColumnTypes = append(t.ColumnTypes, dt)
	}
	return nil
}

func (mb *MemoryBackend) Insert(inst *ast.InsertStatement) error {
	table, ok := mb.Tables[inst.Table.Value]
	if !ok {
		return ErrTableDoesNotExist
	}
	if inst.Values == nil {
		return nil
	}

	var row []MemoryCell
	if len(*inst.Values) != len(table.Columns) {
		return ErrMissingValues
	}
	for _, value := range *inst.Values {
		if value.Kind != ast.LiteralKind {
			fmt.Println("Skipping non-literal.")
			continue
		}
		row = append(row, mb.TokenToCell(value.Literal))
	}
	table.Rows = append(table.Rows, row)
	return nil
}

func (mb *MemoryBackend) TokenToCell(t *token.Token) MemoryCell {
	if t.Kind == token.NumericKind {
		buf := new(bytes.Buffer)
		i, err := strconv.Atoi(t.Value)
		if err != nil {
			panic(err)
		}
		err = binary.Write(buf, binary.BigEndian, int32(i))
		if err != nil {
			panic(err)
		}
		return MemoryCell(buf.Bytes())
	}
	if t.Kind == token.StringKind {
		return MemoryCell(t.Value)
	}
	return nil
}

func (mb *MemoryBackend) Select(slct *ast.SelectStatement) (*Results, error) {
	table, ok := mb.Tables[slct.From.Value]
	if !ok {
		return nil, ErrTableDoesNotExist
	}
	results := [][]Cell{}
	var columns []struct {
		Type ColumnType
		Name string
	}
	for i, row := range table.Rows {
		var result []Cell
		isFirstRow := i == 0
		for _, exp := range slct.Item {
			if exp.Kind != ast.LiteralKind {
				fmt.Println("skipping non-literal expression.")
				continue
			}
			lit := exp.Literal
			if lit.Kind == token.IdentifierKind {
				found := false
				for i, tableCol := range table.Columns {
					if tableCol == lit.Value {
						if isFirstRow {
							columns = append(columns, struct {
								Type ColumnType
								Name string
							}{Type: table.ColumnTypes[i], Name: lit.Value})
						}
						result = append(result, row[i])
						found = true
						break
					}
				}
				if !found {
					return nil, ErrColumnDoesNotExist
				}
				continue
			}
			return nil, ErrColumnDoesNotExist
		}
		results = append(results, result)
	}
	return &Results{
		Columns: columns,
		Rows:    results,
	}, nil
}
