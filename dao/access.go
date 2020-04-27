package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func NewDatabase(user, pwd, network, server, port, db string) (*sql.DB, error) {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s",user, pwd, network, server, port, db)
	return sql.Open("mysql", conn)
}

func QueryIndexResult(db *sql.DB, sql, filter string) (string, error) {
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
		if strings.Contains(s, filter) {
			return s, nil
		}
	}
	return "", errors.New("no match")
}
