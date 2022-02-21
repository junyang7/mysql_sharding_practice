package main

import (
	"fmt"
	"strconv"
	"tool/dao"
	"tool/log"
)

var (
	dbDriver           = "mysql"        // 数据库驱动
	dbProtocol         = "tcp"          // 协议
	dbHost             = "172.16.10.30" // 数据库地址
	dbPort             = 3306           // 数据库端口
	dbBaseName         = "db"           // 试验数据库
	tbBaseName         = "tb"           // 试验数据表
	dbTransferUsername = "transfer"     // 迁移账号（实际业务中操作数据库使用的普通账号）
	dbTransferPassword = "aA!12345"     // 迁移密码
	tbCount            = 32
)

func main() {

	db := dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbTransferUsername, dbTransferPassword, dbProtocol, dbHost, dbPort, dbBaseName))
	statement := "CREATE TABLE `%s`\n(\n    `id`          bigint       NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n    `create_time` datetime     NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',\n    `update_time` datetime     NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',\n    `delete_time` datetime     NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '删除时间',\n    `status`      tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态：1未删除，0已删除',\n    `rid`         bigint unsigned NOT NULL DEFAULT '0' COMMENT '分布式唯一ID',\n    `name`        varchar(255) NOT NULL DEFAULT '' COMMENT '文本信息',\n    `count`       bigint unsigned NOT NULL DEFAULT '0' COMMENT '模拟请求',\n    PRIMARY KEY (`id`),\n    UNIQUE KEY `uk_rid` (`rid`) USING BTREE COMMENT '分布式唯一ID'\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;"
	for i := 0; i < tbCount; i++ {
		tbName := tbBaseName + "_" + strconv.Itoa(i)
		log.Info("正在创建数据表：", tbName)
		dao.Execute(db, fmt.Sprintf(statement, tbName))
		log.Info("成功")

	}
	dao.Close(db)

}
