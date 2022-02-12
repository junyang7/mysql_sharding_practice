package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tool/dao"
	"tool/datetime"
)

const (
	dbName = "db"
	tbName = "tb"
)

var db *sql.DB

func main() {

	fmt.Println("链接数据库：", "mysql")
	db = dao.Connect("mysql", "root:@tcp(127.0.0.1:3306)/mysql")

	fmt.Println("删除数据库：", dbName)
	dao.Execute(db, "DROP DATABASE IF EXISTS "+dbName)

	fmt.Println("创建数据库：", dbName)
	dao.Execute(db, "CREATE DATABASE "+dbName)

	fmt.Println("关闭数据库：", "mysql")
	_ = db.Close()

	fmt.Println("链接数据库：", dbName)
	db = dao.Connect("mysql", "root:@tcp(127.0.0.1:3306)/"+dbName)
	defer db.Close()

	fmt.Println("创建数据表：", tbName)
	dao.Execute(db, "CREATE TABLE `"+tbName+"` (\n  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `add_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',\n  `set_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',\n  `kid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '路由主键',\n  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '文本信息',\n  `count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '模拟请求',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;")

	wd, _ := os.Getwd()
	fName := wd + "/data.txt"
	os.Remove(fName)

	fL := 5000000
	fmt.Println("创建数据集：", fName, fL)
	f, err := os.OpenFile(fName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	_, _ = f.WriteString("id,add_time,set_time,kid,name,count\n")

	for i := 1; i <= fL; i++ {
		s := strconv.Itoa(i)
		d := datetime.Get()
		f.WriteString(strings.Join([]string{s, d, d, s, s, "0"}, ",") + "\n")
	}
	f.Close()

	fmt.Println("写入数据集...")
	dao.Execute(db, "LOAD DATA INFILE '"+fName+"'\nINTO TABLE tb\nFIELDS\nTERMINATED BY ','\nIGNORE 1 LINES")

}
