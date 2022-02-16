package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"
	"tool/dao"
	"tool/datetime"
	"tool/log"
)

var (
	db                       *sql.DB
	dbDriver                 = "mysql"     // 数据库驱动
	dbProtocol               = "tcp"       // 协议
	dbHost                   = "localhost" // 数据库地址
	dbPort                   = 3306        // 数据库端口
	dbTransferUsername       = "transfer"  // 迁移账号（实际业务中操作数据库使用的普通账号）
	dbTransferPassword       = "transfer"  // 迁移密码
	dbCount                  = 1           // 试验数据库数量
	dbBaseName               = "db"        // 试验数据库
	tbBaseName               = "tb"        // 试验数据表
	tbCount                  = 32          // 数据表数据量
	retryTimesBeforeReadOnly = 60          // 最大差异对比补齐次数（只读前）
)

func init() {
	db = dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbTransferUsername, dbTransferPassword, dbProtocol, dbHost, dbPort, dbBaseName))
}
func main() {

	dtS := "1970-01-01 00:00:00"
	dtE := datetime.Get()
	transferByInit(dtS, dtE)

	i := 0
	for {
		i++
		log.Info(i)
		if i > retryTimesBeforeReadOnly {
			dao.Execute(db, "SET GLOBAL read_only = 1;")
		}
		dtS = dtE
		dtE = datetime.Get()
		if transferByDiff(i, dtS, dtE) {
			if checkOk() {
				dao.Close(db)
				if i > retryTimesBeforeReadOnly {
					dao.Execute(db, "SET GLOBAL read_only = 0;")
				}
				log.Info("完成")
				os.Exit(0)
			}
		}
	}

}

func transferByInit(dtS string, dtE string) {

	field := "create_time,update_time,delete_time,status,rid,name,count"
	for i := 0; i < dbCount; i++ {
		for j := 0; j < tbCount; j++ {

			tbName := tbBaseName + "_" + strconv.Itoa(j)
			statement := fmt.Sprintf("INSERT IGNORE INTO %s (%s) SELECT %s FROM %s WHERE update_time >= ? AND update_time <= ? AND rid %% (%d * %d) %% %d = ?", tbName, field, field, tbBaseName, dbCount, tbCount, tbCount)
			parameter := []interface{}{dtS, dtE, j}

			dao.ExecuteWithParameter(db, statement, parameter)
			log.Info(tbName, statement, parameter)

		}
	}

}
func transferByDiff(i int, dtS string, dtE string) bool {

	task := 0
	statement := fmt.Sprintf("SELECT * FROM %s WHERE update_time >= ? AND update_time <= ?", tbBaseName)
	parameter := []interface{}{dtS, dtE}
	rowList := dao.QueryRowListWithParameterList(db, statement, parameter)

	for _, rowTb := range rowList {

		rid, _ := strconv.Atoi(rowTb["rid"])
		m := rid % (dbCount * tbCount)
		dbIndex := m / tbCount
		tbIndex := m % tbCount
		tbName := "tb_" + strconv.Itoa(tbIndex)

		rowTbN := dao.QueryRowWithParameterList(db, fmt.Sprintf("SELECT * FROM %s WHERE rid = ?", tbName), []interface{}{rowTb["rid"]})
		if len(rowTbN) > 0 {
			if rowTb["update_time"] > rowTbN["update_time"] {

				statement := fmt.Sprintf("UPDATE %s SET create_time = ?, update_time = ?, delete_time = ?, status = ?, rid = ?, name = ?, count = ? WHERE rid = ?", tbName)
				parameter := []interface{}{rowTb["create_time"], rowTb["update_time"], rowTb["delete_time"], rowTb["status"], rowTb["rid"], rowTb["name"], rowTb["count"], rowTb["rid"]}

				dao.ExecuteWithParameter(db, statement, parameter)
				log.Info(i, "set", statement, parameter, rid, m, dbIndex, tbIndex)

				task++

			}
			continue
		}

		statement := fmt.Sprintf("INSERT INTO %s (create_time,update_time,delete_time,status,rid,name,count) VALUES (?, ?, ?, ?, ?, ?, ?)", tbName)
		parameter := []interface{}{rowTb["create_time"], rowTb["update_time"], rowTb["delete_time"], rowTb["status"], rowTb["rid"], rowTb["name"], rowTb["count"]}

		dao.ExecuteWithParameter(db, statement, parameter)
		log.Info(i, "add", statement, parameter, rid, m, dbIndex, tbIndex)

		task++

	}

	return task == 0

}
func checkOk() bool {

	for i := 0; i < dbCount; i++ {
		for j := 0; j < tbCount; j++ {

			time.Sleep(time.Second)
			tbName := tbBaseName + "_" + strconv.Itoa(j)
			statement := fmt.Sprintf(
				"SELECT\n    COUNT(*) AS `c`\nFROM\n    (\n        SELECT\n            `%s`.`id` AS `%s_id`,\n            `%s`.`id` AS `%s_id`\n        FROM\n            `%s`\n        LEFT JOIN\n            `%s`\n                ON\n                    `%s`.`create_time` = `%s`.`create_time` AND\n                    `%s`.`update_time` = `%s`.`update_time` AND\n                    `%s`.`delete_time` = `%s`.`delete_time` AND\n                    `%s`.`status` = `%s`.`status`  AND\n                    `%s`.`rid` = `%s`.`rid` AND\n                    `%s`.`name` = `%s`.`name` AND\n                    `%s`.`count` = `%s`.`count`\n        WHERE\n            `%s`.`rid` %% (1 * 32) %% 32 = %d\n    ) AS `t`\nWHERE\n    `t`.`%s_id` IS NULL\n;",
				tbBaseName, tbBaseName,
				tbName, tbName,
				tbBaseName,
				tbName,
				tbBaseName, tbName,
				tbBaseName, tbName,
				tbBaseName, tbName,
				tbBaseName, tbName,
				tbBaseName, tbName,
				tbBaseName, tbName,
				tbBaseName, tbName,
				tbBaseName, j,
				tbName,
			)

			row := dao.QueryRow(db, statement)
			log.Info(i, j, row["c"])

			c, err := strconv.Atoi(row["c"])
			if err != nil {
				log.Info(err)
				return false
			}

			if c > 0 {
				return false
			}

		}
	}

	return true

}
