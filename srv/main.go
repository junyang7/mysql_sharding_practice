package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"tool/dao"
	"tool/datetime"
	"tool/rd"
)

var (
	wg         sync.WaitGroup
	maxId      = 10000000
	dbBaseName = "db"
	tbBaseName = "tb"
	dbCount    = 1
	tbCount    = 32
	dbList     = map[string]*sql.DB{
		"db_0": dao.Connect("mysql", "root:aA$12345@tcp(127.0.0.1:3306)/"+dbBaseName),                 // db
		"db_1": dao.Connect("mysql", "root:aA$12345@tcp(127.0.0.1:3306)/"+dbBaseName+strconv.Itoa(1)), // db_1
	}
	rdList = map[string]*redis.Client{
		"rd_0": rd.Connect("127.0.0.1:6379", "", 0),
	}
	sleep = 5000000
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

// add 模拟新增
func add(c *conf) {
	for {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
		kid := genKid()
		tm := datetime.Get()
		parameter := []interface{}{tm, tm, kid, kid}
		statement := "INSERT INTO %s (add_time, set_time, kid, name) VALUES (?, ?, ?, ?)"
		write(c, statement, parameter, kid)
		fmt.Println(datetime.Get(), "add", parameter)
	}
	wg.Done()
}

// set 模拟修改
func set(c *conf) {
	for {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
		count := time.Now().Unix()
		tm := datetime.Get()
		kid := getRandomKid()
		parameter := []interface{}{count, tm, kid}
		statement := "UPDATE %s SET count = ?, set_time = ? WHERE kid = ?"
		write(c, statement, parameter, kid)
		fmt.Println(datetime.Get(), "set", parameter)
	}
	wg.Done()
}
func del(c *conf) {
	for {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
		kid := getRandomKid()
		parameter := []interface{}{kid}
		statement := "DELETE FROM %s WHERE kid = ?"
		write(c, statement, parameter, kid)
		fmt.Println(datetime.Get(), "del", parameter)
	}
	wg.Done()
}

func s(db *sql.DB, sql string, tbBaseName string, parameterList []interface{}) {
	statement := fmt.Sprintf(sql, tbBaseName)
	dao.ExecuteWithParameter(db, statement, parameterList)
	fmt.Println(tbBaseName, statement, parameterList)
}
func m(db *sql.DB, sql string, tbBaseName string, parameterList []interface{}, kid int, dbCount int, tbCount int) {
	m := kid % (dbCount * tbCount)
	//多库时候用到
	//dbIndex := m / tbCount
	tbIndex := m % tbCount
	tbName := tbBaseName + "_" + strconv.Itoa(tbIndex)
	statement := fmt.Sprintf(sql, tbName)
	dao.ExecuteWithParameter(db, statement, parameterList)
	fmt.Println(tbName, statement, parameterList)
}

// genKid 生成全局唯一KID
func genKid() int {
	kid, _ := rdList["rd_0"].Incr("kid").Result()
	return int(kid)
}

// getKid 获取当前最大KID的值
func getKid() int {
	kid, _ := rdList["rd_0"].Get("kid").Int()
	return kid
}

// getRandomKid 随机获取0-最大KID之间的一个KID，用来模拟修改和删除
func getRandomKid() int {
	return rand.Intn(getKid())
}

// write 处理数据库写操作
func write(c *conf, statement string, parameterList []interface{}, kid int) {
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
}
