package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbName := "db_0"

	var (
		db  *sql.DB
		err error
	)

	fmt.Println("链接数据库：", "mysql")
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/mysql")
	if err != nil {
		panic(err)
	}
	fmt.Println("检查数据库是否已经存在，如果存在则删除：", dbName)
	if _, err := db.Exec("DROP DATABASE IF EXISTS " + dbName); err != nil {
		panic(err)
	}
	fmt.Println("创建数据库：", dbName)
	if _, err := db.Exec("CREATE DATABASE " + dbName); err != nil {
		panic(err)
	}
	fmt.Println("关闭数据库：", "mysql")
	_ = db.Close()

	fmt.Println("链接数据库：", dbName)
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tbName := "tb"
	fmt.Println("创建数据表：", tbName)
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
 id bigint(20) NOT NULL AUTO_INCREMENT,
 name varchar(255) NOT NULL DEFAULT '',
 PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`, tbName)); err != nil {
		panic(err)
	}

	for i := 1; i <= 100000; i++ {
		fmt.Println("写入模拟数据：", i)
		if _, err := db.Exec(fmt.Sprintf(`
INSERT INTO %s (id, name) VALUES (?, ?);
`, tbName), i, i); err != nil {
			panic(err)
		}
	}

	//	for i := 0; i < 4; i++ {
	//		tbName := "tb_" + strconv.Itoa(i)
	//		fmt.Println("创建数据表：", tbName)
	//		if _, err := db.Exec(fmt.Sprintf(`
	//CREATE TABLE %s (
	// id bigint(20) NOT NULL AUTO_INCREMENT,
	// name varchar(255) NOT NULL DEFAULT '',
	// PRIMARY KEY (id)
	//) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	//`, tbName)); err != nil {
	//			panic(err)
	//		}
	//	}

}
