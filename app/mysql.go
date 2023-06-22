package app

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysql() *sql.DB {
	db, err := sql.Open("mysql", "root:pass@tcp(localhost:3306)/tikets")
	if err != nil {
		panic(err)
	}
	return db
}
