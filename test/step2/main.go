package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"tool/dao"
	"tool/datetime"
	"tool/log"
)

var (
	db                       *sql.DB
	dbDriver                 = "mysql"               // 数据库驱动
	dbProtocol               = "tcp"                 // 协议
	dbHost                   = "localhost"           // 数据库地址
	dbPort                   = 3306                  // 数据库端口
	dbTransferUsername       = "transfer"            // 迁移账号（实际业务中操作数据库使用的普通账号）
	dbTransferPassword       = "transfer"            // 迁移密码
	dbCount                  = 1                     // 试验数据库数量
	dbBaseName               = "db"                  // 试验数据库
	tbBaseName               = "tb"                  // 试验数据表
	tbCount                  = 32                    // 数据表数据量
	retryTimesBeforeReadOnly = 60                    // 最大差异对比补齐次数（只读前）
	dbDefaultDatetime        = "1970-01-01 00:00:00" // 数据库默认时间
	sleep                    = 100000                // 每个耗时SQL间隔
)

func init() {
	db = dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbTransferUsername, dbTransferPassword, dbProtocol, dbHost, dbPort, dbBaseName))
}
func main() {

	dtS := dbDefaultDatetime
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
				if i > retryTimesBeforeReadOnly {
					dao.Execute(db, "SET GLOBAL read_only = 0;")
				}
				log.Info("完成")
				dao.Close(db)
				os.Exit(0)
			}
		}

		time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))

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
			time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))

		}
	}

}

func transferByDiff(i int, dtS string, dtE string) bool {

	task := 0
	field := "create_time,update_time,delete_time,status,rid,name,count"
	statement := fmt.Sprintf("SELECT rid,update_time FROM %s WHERE update_time >= ? AND update_time <= ?", tbBaseName)
	parameter := []interface{}{dtS, dtE}
	rowList := dao.QueryRowListWithParameterList(db, statement, parameter)

	for _, rowTb := range rowList {

		time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
		rid, _ := strconv.Atoi(rowTb["rid"])
		m := rid % (dbCount * tbCount)
		dbIndex := m / tbCount
		tbIndex := m % tbCount
		tbName := "tb_" + strconv.Itoa(tbIndex)

		rowTbN := dao.QueryRowWithParameterList(db, fmt.Sprintf("SELECT * FROM %s WHERE rid = ?", tbName), []interface{}{rowTb["rid"]})
		if len(rowTbN) > 0 {
			continue
		}

		statement := fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s WHERE rid = ?", tbName, field, field, tbBaseName)
		parameter := []interface{}{rowTb["rid"]}

		dao.ExecuteWithParameter(db, statement, parameter)
		log.Info(i, "add", statement, parameter, rid, m, dbIndex, tbIndex)

		task++

	}

	return task == 0

}
func checkOk() bool {

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))

	for i := 0; i < dbCount; i++ {
		for j := 0; j < tbCount; j++ {

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

			time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))

		}
	}

	return true

}
