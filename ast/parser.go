package ast

import (
"fmt"

"github.com/kr/pretty"
"github.com/pingcap/parser"
"github.com/pingcap/parser/ast"
_ "github.com/pingcap/tidb/types/parser_driver"
)

func TiParse(sql, charset, collation string) ([]ast.StmtNode, error) {
	p := parser.New()
	stmt, warn, err := p.Parse(sql, charset, collation)
	for _, w := range warn {
		fmt.Println(w.Error())
	}
	return stmt, err
}

func PrintPrettyStmtNode(sql, charset, collation string) {
	tree, err := TiParse(sql, charset, collation)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = pretty.Println(tree)
	}
}

func TiDigestHash(sql string) string {
	return parser.DigestHash(sql)
}

func ParseType(sql, charset, collation string) (tp []string, err error) {
	//解析ast
	stmts, e := TiParse(sql, charset, collation)

	if e != nil {
		return nil, e
	}

	for _, stmt := range stmts {
		switch stmt.(type) {
		case *ast.CreateTableStmt:
			tp = append(tp, "CreateTable")
		case *ast.DropTableStmt:
			tp = append(tp, "DropTable")
		case *ast.TruncateTableStmt:
			tp = append(tp, "TruncateTable")
		case *ast.AlterTableStmt:
			tp = append(tp, "AlterTable")
		default:
			tp = append(tp, "OtherType")
		}
	}
	return
}
