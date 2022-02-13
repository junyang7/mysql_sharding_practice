package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"tool/dao"
)

var (
	dbBaseName = "db"
	tbBaseName = "tb"
	tbCount    = 32
	dbList     = map[string]*sql.DB{
		"db_0": dao.Connect("mysql", "root:aA$12345@tcp(127.0.0.1:3306)/"+dbBaseName),
	}
)

func main() {

	for i := 0; i < tbCount; i++ {
		tbName := tbBaseName + "_" + strconv.Itoa(i)
		fmt.Println("创建数据表：", tbName)
		dao.Execute(dbList["db_0"], "DROP TABLE IF EXISTS `"+tbName+"`")
		dao.Execute(dbList["db_0"], "CREATE TABLE `"+tbName+"` (\n  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',\n  `add_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',\n  `set_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',\n  `kid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '路由主键',\n  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '文本信息',\n  `count` int unsigned NOT NULL DEFAULT '0' COMMENT '模拟请求',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `uk_kid` (`kid`) USING BTREE COMMENT '全局唯一ID'\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;")
	}

}
