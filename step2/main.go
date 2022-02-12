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
	dbCount    = 1
	tbCount    = 32
	dbList     = map[string]*sql.DB{
		"db_0": dao.Connect("mysql", "root:@tcp(127.0.0.1:3306)/"+dbBaseName),
	}
)

func main() {

	page := 1
	size := 10000

	for {

		fmt.Println(page)

		statement := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT %d, %d", tbBaseName, (page-1)*size, size)
		parameterList := []interface{}{}
		rowList := dao.QueryRowListWithParameterList(dbList["db_0"], statement, parameterList)

		if len(rowList) == 0 {
			break
		}

		page++

		for _, tbRow := range rowList {
			kid, err := strconv.ParseInt(tbRow["kid"], 10, 64)
			if err != nil {
				panic(err)
			}
			m := int(kid % (int64(dbCount) * int64(tbCount)))
			//多库时候用到
			//dbIndex := m / tbCount
			tbIndex := m % tbCount
			tbName := tbBaseName + "_" + strconv.Itoa(tbIndex)

			// 从多表中取数据
			tbNRow := dao.QueryRowWithParameterList(dbList["db_0"], "SELECT * FROM "+tbName+" WHERE kid = ?", []interface{}{tbRow["kid"]})
			if len(tbNRow) == 0 {
				dao.ExecuteWithParameter(dbList["db_0"], "INSERT INTO "+tbName+" (add_time, set_time, kid, name, count) VALUES (?, ?, ?, ?, ?)", []interface{}{tbRow["add_time"], tbRow["set_time"], tbRow["kid"], tbRow["name"], tbRow["count"]})
				continue
			}

			if tbRow["set_time"] > tbNRow["set_time"] {
				dao.ExecuteWithParameter(dbList["db_0"], "UPDATE "+tbName+" SET set_time = ?, count = ? WHERE kid = ?", []interface{}{tbRow["set_time"], tbRow["count"], tbRow["kid"]})
			}

		}
	}

}
