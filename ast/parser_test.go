package ast

import (
	"fmt"
	"testing"

	"github.com/pingcap/parser/ast"
)

func TestPrintPrettyStmtNode(t *testing.T) {
	sqls := []string{
		`create table a (a int primary key, b int);`,
	}

	for _, sql := range sqls {
		PrintPrettyStmtNode(sql, "", "")
	}
}

type visitor struct {
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	fmt.Printf("other type: %T\n", in)
	return in, false
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	/*
		switch t :=in.(type) {
		case *ast.TableName:
			if t.Schema.String() != "" {
				fmt.Printf("table schema : %s\n", t.Schema.String())
			}
			fmt.Printf("table name : %s\n", t.Name)
		case *ast.ColumnName:
			v.columnNum ++
		case *driver.ValueExpr:
			fmt.Printf("#name: %d\n", t.GetInt64())
		case *ast.Limit:
			v.isLimit ++
		default:
			fmt.Printf("other type: %T\n", t)
		}
	*/
	return in, true
}

func TestTraverParser(t *testing.T) {
	sqls := []string{
		`create table a (a int primary key, b int)TABLET_MAX_SIZE = 268435456, TABLET_BLOCK_SIZE = 16384, REPLICA_NUM = 2, DEFAULT CHARSET = 'utf8';`,
	}
	stmts, _ := TiParse(string(sqls[0]), "", "")
	v := visitor{}
	for _, stmt := range stmts {
		stmt.Accept(&v)
	}
	//fmt.Printf("column num: %d\n", v.columnNum)
	//fmt.Printf("is limit :%d\n", v.isLimit)
}
