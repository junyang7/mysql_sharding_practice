package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println(db)

	res, err := db.Exec("CREATE DATABASE db_0")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

}
