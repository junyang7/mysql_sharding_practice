package main

import (
	"database/sql"
	"fmt"
	"mysql_sharding/tool/dao"
	"mysql_sharding/tool/datetime"
)

const (
	dbName = "db"
	tbName = "tb"
)

var db *sql.DB

func main() {

	fmt.Println("链接数据库：", dbName)
	db = dao.Connect("mysql", "root:@tcp(127.0.0.1:3306)/"+dbName)
	defer db.Close()

	fmt.Println("创建数据表：", tbName)
	dao.Execute(db, "CREATE TABLE `"+tbName+"` (\n  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `add_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',\n  `set_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',\n  `kid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '路由主键',\n  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '文本信息',\n  `count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '模拟请求',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;")

	for i := 1; i <= 1000000; i++ {
		parameterList := []interface{}{i, i, i, 0, datetime.Get(), datetime.Get()}
		fmt.Println("写入模拟数据：", parameterList)
		dao.ExecuteWithParameter(db, "INSERT INTO "+tbName+" (id, kid, name, count, add_time, set_time) VALUES (?, ?, ?, ?, ?, ?)", parameterList)
	}

	fmt.Println("初始化完成")

}
