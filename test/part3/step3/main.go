package main

import (
	"fmt"
	"tool/dao"
)

var (
	dbDriver           = "mysql"     // 数据库驱动
	dbProtocol         = "tcp"       // 协议
	dbHost             = "localhost" // 数据库地址
	dbPort             = 3306        // 数据库端口
	dbBaseName         = "db"        // 试验数据库
	tbBaseName         = "tb"        // 试验数据表
	dbTransferUsername = "transfer"  // 迁移账号（实际业务中操作数据库使用的普通账号）
	dbTransferPassword = "transfer"  // 迁移密码
)

func main() {

	db := dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbTransferUsername, dbTransferPassword, dbProtocol, dbHost, dbPort, dbBaseName))
	dao.Execute(db, fmt.Sprintf("TRUNCATE TABLE %s;", tbBaseName))
	dao.Execute(db, fmt.Sprintf("DROP TABLE %s;", tbBaseName))
	dao.Close(db)

}
