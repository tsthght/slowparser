package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func NewDatabase(user, pwd, network, server string, port int, db string) (*sql.DB, error) {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",user, pwd, network, server, port, db)
	return sql.Open("mysql", conn)
}

func QueryIndexResult(db *sql.DB, sql, filter string) (string, error) {
	fmt.Printf("##sql %s\n", sql)
	rows, err := db.Query(sql)
	if err != nil {
		return "", err
	}
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return "", err
		}
		fmt.Printf("##row %s\n", s)
		fmt.Printf("##filter %s\n", filter)
		if strings.Contains(s, filter) {
			return s, nil
		}
	}
	return "", errors.New("no match rows")
}
