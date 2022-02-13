package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"tool/dao"
	"tool/datetime"
)

var (
	dbBaseName = "db"
	tbBaseName = "tb"
	dbCount    = 1
	tbCount    = 32
	dbList     = map[string]*sql.DB{
		"db_0": dao.Connect("mysql", "root:aA$12345@tcp(127.0.0.1:3306)/"+dbBaseName),
	}
)

func main() {

	// 第1次迁移：插入
	dtS := "1970-01-01 00:00:00"
	dtE := datetime.Get()
	fmt.Println(dtS)
	fmt.Println(dtE)
	add(dtS, dtE)

	for i := 0; i < 100; i++ {
		dtS = dtE
		dtE = datetime.Get()
		fmt.Println(i)
		fmt.Println(dtS)
		fmt.Println(dtE)
		set(dtS, dtE)
	}

}

func add(dtS string, dtE string) {
	field := "add_time, set_time, kid, name, count"
	for i := 0; i < dbCount; i++ {
		for j := 0; j < tbCount; j++ {
			statement := fmt.Sprintf("INSERT IGNORE INTO %s (%s) SELECT %s FROM %s WHERE set_time >= ? AND set_time <= ? AND kid %% (%d * %d) %% %d = ?", tbBaseName+"_"+strconv.Itoa(j), field, field, tbBaseName, dbCount, tbCount, tbCount)
			parameter := []interface{}{dtS, dtE, j}
			fmt.Println(statement)
			fmt.Println(parameter)
			dao.ExecuteWithParameter(dbList["db_0"], statement, parameter)
		}
	}
	fmt.Println(datetime.Get()+": 完成", dtS, dtE)
}

func set(dtS string, dtE string) {
	num := 0
	statement := fmt.Sprintf("SELECT * FROM %s WHERE set_time >= ? AND set_time <= ?", tbBaseName)
	parameter := []interface{}{dtS, dtE}
	rowList := dao.QueryRowListWithParameterList(dbList["db_0"], statement, parameter)
	for _, rowTb := range rowList {
		kid, _ := strconv.Atoi(rowTb["kid"])
		m := kid % (dbCount * tbCount)
		//多库时候用到
		//dbIndex := m / tbCount
		tbIndex := m % tbCount
		tbName := "tb_" + strconv.Itoa(tbIndex)
		rowTbN := dao.QueryRowWithParameterList(dbList["db_0"], fmt.Sprintf("SELECT * FROM %s WHERE kid = ?", tbName), []interface{}{rowTb["kid"]})
		if len(rowTbN) > 0 {
			if rowTb["set_time"] > rowTbN["set_time"] {
				num++
				// set
				statement := fmt.Sprintf("UPDATE %s SET add_time = ?, set_time = ?, kid = ?, name = ?, count = ? WHERE kid = ?", tbName)
				parameter := []interface{}{rowTb["add_time"], rowTb["set_time"], rowTb["kid"], rowTb["name"], rowTb["count"], rowTb["kid"]}
				dao.ExecuteWithParameter(dbList["db_0"], statement, parameter)
				fmt.Println(datetime.Get(), "set", rowTb["kid"], statement, parameter)
			}
			continue
		}
		num++
		// add
		statement := fmt.Sprintf("INSERT INTO %s (add_time, set_time, kid, name, count) VALUES (?, ?, ?, ?, ?)", tbName)
		parameter := []interface{}{rowTb["add_time"], rowTb["set_time"], rowTb["kid"], rowTb["name"], rowTb["count"]}
		dao.ExecuteWithParameter(dbList["db_0"], statement, parameter)
		fmt.Println(datetime.Get(), "add", rowTb["kid"], statement, parameter)
	}
	if num == 0 {
		fmt.Println("同步完成")
		os.Exit(0)
	}
}
