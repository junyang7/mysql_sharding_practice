package main

import (
	"database/sql"
	"fmt"
)

// 第一次分表
// 数据库数量不变，将一张达标拆分成4张小表
// 规则：
// 		中间变量 = KEY % (库数量 * 每个库的表数量)
//		库 = 中间变量 / 每个库的表数量
//      表 = 中间变量 % 每个库的表数量

func main() {
	dbName := "db_0"

	var (
		db  *sql.DB
		err error
	)
	fmt.Println("链接数据库：", dbName)
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("第一次分表")
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
