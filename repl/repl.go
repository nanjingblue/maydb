package repl

import (
	"bufio"
	"fmt"
	"github.com/nanjingblue/maydb/ast"
	"github.com/nanjingblue/maydb/backend"
	"github.com/nanjingblue/maydb/parser"
	"io"
	"strings"
)

func Start(in io.Reader, out io.Writer) {
	mb := backend.NewMemoryBacked()

	reader := bufio.NewReader(in)
	fmt.Println("Welcome to gosql.")
	for {
		fmt.Print("# ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		text = text[:len(text)-1]

		asts, err := parser.Parse(text)
		if err != nil {
			panic(err)
		}
		for _, stmt := range asts.Statements {
			switch stmt.Kind {
			case ast.CreateTableKind:
				err = mb.CreateTable(asts.Statements[0].CreateTableStatement)
				if err != nil {
					panic(err)
				}
				fmt.Println("ok")
			case ast.InsertKind:
				err = mb.Insert(stmt.InsertStatement)
				if err != nil {
					panic(err)
				}
				fmt.Println("ok")
			case ast.SelectKind:
				results, err := mb.Select(stmt.SelectStatement)
				if err != nil {
					panic(err)
				}
				for _, col := range results.Columns {
					io.WriteString(out, fmt.Sprintf("| %s", col.Name))
				}
				io.WriteString(out, "|")

				for i := 0; i < 20; i++ {
					io.WriteString(out, "=")
				}
				io.WriteString(out, "\n")

				for _, result := range results.Rows {
					io.WriteString(out, "|")

					for i, cell := range result {
						typ := results.Columns[i].Type
						s := ""
						switch typ {
						case backend.IntType:
							s = fmt.Sprintf("%d", cell.AsInt())
						case backend.TextType:
							s = cell.AsText()
						}
						io.WriteString(out, fmt.Sprintf("%s |", s))
					}
					fmt.Println()
				}
				fmt.Println("ok")
			}
		}
	}
}
