package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"tool/dao"
	"tool/datetime"
	"tool/file"
	"tool/log"
	"tool/rd"
)

var (
	filename             = "data.txt"                                 // 模拟数据文件名
	rowCount             = 5000000                                    // 模拟数据行数
	redisHost            = "127.0.0.1"                                // Redis服务器地址
	redisPort            = 6379                                       // Redis服务端口
	redisPassword        = ""                                         // Redis密码
	redisDbIndex         = 0                                          // Redis数据库索引
	dbDriver             = "mysql"                                    // 数据库驱动
	dbUsername           = "root"                                     // 超级管理员账号
	dbPassword           = "root"                                     // 超级管理员密码
	dbProtocol           = "tcp"                                      // 协议
	dbHost               = "localhost"                                // 数据库地址
	dbPort               = 3306                                       // 数据库端口
	dbTempName           = "information_schema"                       // 默认数据库，需要通过次数据库构建dsn，建立链接后才能创建我们需要的数据库
	dbBaseName           = "db"                                       // 试验数据库
	tbBaseName           = "tb"                                       // 试验数据表
	dbBusinessUsername   = "business"                                 // 业务账号（实际业务中操作数据库使用的普通账号）
	dbBusinessPassword   = "business"                                 // 业务密码
	dbBusinessHost       = "%"                                        // 业务主机
	dbBusinessPrivileges = "INSERT,DELETE,UPDATE,SELECT"              // 业务权限（增删改查）
	dbBusinessScope      = dbBaseName + ".*"                          // 业务账号权限范围
	dbTransferUsername   = "transfer"                                 // 迁移账号（实际业务中操作数据库使用的普通账号）
	dbTransferPassword   = "transfer"                                 // 迁移密码
	dbTransferHost       = "%"                                        // 迁移主机
	dbTransferPrivileges = "INSERT,DELETE,UPDATE,SELECT,CREATE,SUPER" // 迁移权限
	dbTransferScope      = "*.*"                                      // 迁移账号权限范围
)

func main() {

	// 创建模拟数据集
	createData()
	// 设置分布式唯一ID初始值
	setRid()
	// 创建数据库
	createDatabase()

	log.Info("正在建立数据库链接：", dbBaseName)
	db := dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbUsername, dbPassword, dbProtocol, dbHost, dbPort, dbBaseName))
	log.Info("成功")

	// 创建数据表
	createTable(db)
	// 导入模拟数据集
	loadData(db)
	// 创建业务账号
	createUser(db, dbBusinessUsername, dbBusinessHost, dbBusinessPassword, dbBusinessPrivileges, dbBusinessScope)
	// 创建迁移账号
	createUser(db, dbTransferUsername, dbTransferHost, dbTransferPassword, dbTransferPrivileges, dbTransferScope)

	log.Info("正在关闭数据库链接：", dbBaseName)
	dao.Close(db)
	log.Info("成功")
	log.Info("环境准备完成")

}

// createData 创建模拟数据集
func createData() {

	if file.IsExists(filename) {
		log.Info("正在删除数据集：", filename)
		file.Unlink(filename)
		log.Info("成功")
	}

	log.Info("正在创建数据集：", filename, rowCount)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	_, _ = f.WriteString("create_time,rid,name\n")
	for i := 1; i <= rowCount; i++ {
		s := strconv.Itoa(i)
		if _, err := f.WriteString(strings.Join([]string{datetime.Get(), s, "name_" + s}, ",") + "\n"); err != nil {
			panic(err)
		}
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	log.Info("成功")

}

// setRid 设置路由主键起始值
func setRid() {

	log.Info("正在设置rid的值：", rowCount)
	r := rd.Connect(redisHost+":"+strconv.Itoa(redisPort), redisPassword, redisDbIndex)
	if err := r.Ping().Err(); err != nil {
		panic(err)
	}
	if err := r.Set("rid", rowCount, 0).Err(); err != nil {
		panic(err)
	}
	if err := r.Close(); err != nil {
		panic(err)
	}
	log.Info("成功")

}

// createDatabase 创建数据库
func createDatabase() {

	log.Info("正在建立数据库链接：", dbTempName)
	db := dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbUsername, dbPassword, dbProtocol, dbHost, dbPort, dbTempName))
	log.Info("成功")

	log.Info("正在删除数据库（如果存在的话）：", dbBaseName)
	dao.Execute(db, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbBaseName))
	log.Info("成功")

	log.Info("正在创建数据库：", dbBaseName)
	dao.Execute(db, fmt.Sprintf("CREATE DATABASE %s", dbBaseName))
	log.Info("成功")

	log.Info("正在关闭数据库链接：", dbTempName)
	dao.Close(db)
	log.Info("成功")

}

// createTable 创建数据表
func createTable(db *sql.DB) {

	log.Info("正在创建数据表：", tbBaseName)
	dao.Execute(db, fmt.Sprintf("CREATE TABLE `%s`\n(\n    `id`          bigint       NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n    `create_time` datetime     NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',\n    `update_time` datetime     NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',\n    `delete_time` datetime     NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '删除时间',\n    `status`      tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态：1未删除，0已删除',\n    `rid`         bigint unsigned NOT NULL DEFAULT '0' COMMENT '分布式唯一ID',\n    `name`        varchar(255) NOT NULL DEFAULT '' COMMENT '文本信息',\n    `count`       bigint unsigned NOT NULL DEFAULT '0' COMMENT '模拟请求',\n    PRIMARY KEY (`id`),\n    UNIQUE KEY `uk_rid` (`rid`) USING BTREE COMMENT '分布式唯一ID'\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;", tbBaseName))
	log.Info("成功")

}

// loadData 导入模拟数据集
func loadData(db *sql.DB) {

	log.Info("正在导入数据集：", filename, rowCount)
	dao.Execute(db, fmt.Sprintf("LOAD DATA INFILE '%s'\nINTO TABLE %s\nFIELDS\nTERMINATED BY ','\nIGNORE 1 LINES (create_time,rid,name)", filename, tbBaseName))
	log.Info("成功")

}

// createUser 创建数据库用户
func createUser(db *sql.DB, username string, host string, password string, privileges string, scope string) {

	if row := dao.QueryRowWithParameterList(db, "SELECT * FROM mysql.user WHERE User = ? AND Host = ?", []interface{}{username, host}); len(row) > 0 {
		log.Info("正在删除账号：", username, host)
		dao.Execute(db, fmt.Sprintf("DROP USER '%s'@'%s'", username, host))
		log.Info("成功")
	}

	log.Info("正在创建账号：", username, host, password)
	dao.Execute(db, fmt.Sprintf("CREATE USER '%s'@'%s' IDENTIFIED BY '%s'", username, host, password))
	log.Info("成功")

	log.Info("正在授予权限：", privileges)
	dao.Execute(db, fmt.Sprintf("GRANT %s ON %s TO '%s'@'%s'", privileges, scope, username, host))
	log.Info("成功")

	log.Info("正在刷新权限...")
	dao.Execute(db, "FLUSH PRIVILEGES")
	log.Info("成功")

}

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filename = dir + "/" + filename
}
