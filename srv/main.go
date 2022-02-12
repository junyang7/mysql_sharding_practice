package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"tool/dao"
	"tool/datetime"
)

var (
	wg         sync.WaitGroup
	maxId      = 10000000
	dbBaseName = "db"
	tbBaseName = "tb"
	dbCount    = 1
	tbCount    = 32
	dbList     = map[string]*sql.DB{
		"db_0": dao.Connect("mysql", "root:@tcp(127.0.0.1:3306)/"+dbBaseName),                 // db
		"db_1": dao.Connect("mysql", "root:@tcp(127.0.0.1:3306)/"+dbBaseName+strconv.Itoa(1)), // db_1
	}
)

type conf struct {
	w int // 写，0单表（默认0），1多表，2单表多表
	r int // 读，0单表（默认0），1多表
}

func main() {

	///**
	//读单表，写单表
	//*/
	//do(&conf{0, 0})

	/**
	读单表，写单表多表
	*/
	do(&conf{2, 0})

}

func do(c *conf) {
	wg.Add(3)
	go add(c)
	go set(c)
	go del(c)
	wg.Wait()
}
func add(c *conf) {
	for {

		time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000000)))

		kid := time.Now().UnixMicro()
		tm := datetime.Get()
		parameterList := []interface{}{tm, tm, kid, kid}
		statement := "INSERT INTO %s (add_time, set_time, kid, name) VALUES (?, ?, ?, ?)"

		if c.w == 1 {
			// 写：多表
			m(dbList["db_0"], statement, tbBaseName, parameterList, kid, dbCount, tbCount)
		} else if c.w == 2 {
			// 写：单表+多表
			s(dbList["db_0"], statement, tbBaseName, parameterList)
			m(dbList["db_0"], statement, tbBaseName, parameterList, kid, dbCount, tbCount)
		} else {
			// 默认：写：单表
			s(dbList["db_0"], statement, tbBaseName, parameterList)
		}
		fmt.Println(datetime.Get(), "add", parameterList)

	}
	wg.Done()
}
func set(c *conf) {
	for {

		time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000000)))

		count := time.Now().Unix()
		tm := datetime.Get()
		kid := rand.Int63n(time.Now().UnixMicro())
		parameterList := []interface{}{count, tm, kid}
		statement := "UPDATE %s SET count = ?, set_time = ? WHERE id = ?"

		if c.w == 1 {
			// 写：多表
			m(dbList["db_0"], statement, tbBaseName, parameterList, kid, dbCount, tbCount)
		} else if c.w == 2 {
			// 写：单表+多表
			s(dbList["db_0"], statement, tbBaseName, parameterList)
			m(dbList["db_0"], statement, tbBaseName, parameterList, kid, dbCount, tbCount)
		} else {
			// 默认：写：单表
			s(dbList["db_0"], statement, tbBaseName, parameterList)
		}

		fmt.Println(datetime.Get(), "set", parameterList)

	}
	wg.Done()
}
func del(c *conf) {
	for {

		time.Sleep(time.Microsecond * time.Duration(rand.Intn(1000000)))

		kid := rand.Int63n(time.Now().UnixMicro())
		parameterList := []interface{}{kid}
		statement := "DELETE FROM %s WHERE id = ?"

		if c.w == 1 {
			// 写：多表
			m(dbList["db_0"], statement, tbBaseName, parameterList, kid, dbCount, tbCount)
		} else if c.w == 2 {
			// 写：单表+多表
			s(dbList["db_0"], statement, tbBaseName, parameterList)
			m(dbList["db_0"], statement, tbBaseName, parameterList, kid, dbCount, tbCount)
		} else {
			// 默认：写：单表
			s(dbList["db_0"], statement, tbBaseName, parameterList)
		}

		fmt.Println(datetime.Get(), "del", parameterList)

	}
	wg.Done()
}

func s(db *sql.DB, sql string, tbBaseName string, parameterList []interface{}) {
	dao.ExecuteWithParameter(db, fmt.Sprintf(sql, tbBaseName), parameterList)
}
func m(db *sql.DB, sql string, tbBaseName string, parameterList []interface{}, kid int64, dbCount int, tbCount int) {
	m := int(kid % (int64(dbCount) * int64(tbCount)))
	//多库时候用到
	//dbIndex := m / tbCount
	tbIndex := m % tbCount
	dao.ExecuteWithParameter(db, fmt.Sprintf(sql, tbBaseName+"_"+strconv.Itoa(tbIndex)), parameterList)
}
