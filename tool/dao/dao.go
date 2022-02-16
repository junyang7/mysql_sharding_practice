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
func Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		panic(err)
	}
}
func Execute(db *sql.DB, statement string) {
	if _, err := db.Exec(statement); err != nil {
		panic(err)
	}
}
func ExecuteWithParameter(db *sql.DB, statement string, parameterList []interface{}) {
	if _, err := db.Exec(statement, parameterList...); err != nil {
		panic(err)
	}
}
func QueryRowList(db *sql.DB, statement string) []map[string]string {

	rowList, err := db.Query(statement)
	if err != nil {
		panic(err)
	}

	fieldList, err := rowList.Columns()
	if err != nil {
		panic(err)
	}

	dest := make([]interface{}, len(fieldList))
	for i, _ := range dest {
		dest[i] = new(sql.RawBytes)
	}

	res := make([]map[string]string, 0)
	for rowList.Next() {
		err := rowList.Scan(dest...)
		if err != nil {
			panic(err)
		}
		row := make(map[string]string)
		for i, value := range dest {
			row[fieldList[i]] = string(*(value.(interface{}).(*sql.RawBytes)))
		}
		res = append(res, row)
	}

	_ = rowList.Close()
	return res

}
func QueryRowListWithParameterList(db *sql.DB, statement string, parameterList []interface{}) []map[string]string {

	rowList, err := db.Query(statement, parameterList...)
	if err != nil {
		panic(err)
	}

	fieldList, err := rowList.Columns()
	if err != nil {
		panic(err)
	}

	dest := make([]interface{}, len(fieldList))
	for i, _ := range dest {
		dest[i] = new(sql.RawBytes)
	}

	res := make([]map[string]string, 0)
	for rowList.Next() {
		err := rowList.Scan(dest...)
		if err != nil {
			panic(err)
		}
		row := make(map[string]string)
		for i, value := range dest {
			row[fieldList[i]] = string(*(value.(interface{}).(*sql.RawBytes)))
		}
		res = append(res, row)
	}

	_ = rowList.Close()
	return res

}
func QueryRow(db *sql.DB, statement string) map[string]string {
	res := QueryRowList(db, statement)
	if len(res) > 0 {
		return res[0]
	}
	return map[string]string{}
}
func QueryRowWithParameterList(db *sql.DB, statement string, parameterList []interface{}) map[string]string {
	res := QueryRowListWithParameterList(db, statement, parameterList)
	if len(res) > 0 {
		return res[0]
	}
	return map[string]string{}
}
