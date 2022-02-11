package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect(driverName string, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}
func Execute(db *sql.DB, sql string) {
	if _, err := db.Exec(sql); err != nil {
		panic(err)
	}
}
func ExecuteWithParameter(db *sql.DB, sql string, parameterList []interface{}) {
	if _, err := db.Exec(sql, parameterList...); err != nil {
		panic(err)
	}
}
