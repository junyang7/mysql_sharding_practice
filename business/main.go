package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
	"tool/dao"
	"tool/datetime"
	"tool/log"
	"tool/rid"
)

var (
	wg                   sync.WaitGroup
	db                   *sql.DB
	pid                  int
	pidFilename          = "pid.txt"
	sleep                = 500000
	redisHost            = "127.0.0.1" // Redis服务器地址
	redisPort            = 6379        // Redis服务端口
	redisPassword        = ""          // Redis密码
	redisDbIndex         = 0
	tbBaseName           = "tb"
	tbCount              = 32
	dbDriver             = "mysql"                       // 数据库驱动
	dbHost               = "localhost"                   // 数据库地址
	dbProtocol           = "tcp"                         // 协议
	dbPort               = 3306                          // 数据库端口
	dbBusinessUsername   = "business"                    // 业务账号（实际业务中操作数据库使用的普通账号）
	dbBusinessPassword   = "business"                    // 业务密码
	dbBusinessHost       = "%"                           // 业务主机
	dbBusinessPrivileges = "INSERT,DELETE,UPDATE,SELECT" // 业务权限（增删改查）
	dbBaseName           = "db"                          // 试验数据库
	dbCount              = 1                             // 试验数据库数量
	confFilename         = "app.json"                    // 项目配置文件路径
	conf                 = &Conf{}                       // 项目配置文件解析结果（项目读写控制）
)

type Conf struct {
	W          int `json:"w"`
	R          int `json:"r"`
	IsReadOnly int `json:"is_read_only"`
}

func init() {

	rid.Init(redisHost, redisPort, redisPassword, redisDbIndex)
	db = dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbBusinessUsername, dbBusinessPassword, dbProtocol, dbHost, dbPort, dbBaseName))

}
func main() {

	// 预处理
	prepare()
	// 工作
	do(conf)

}
func prepare() {

	if len(os.Args) == 1 {
		fmt.Println(`
Usage:
	start | restart | stop
	start	启动
	restart	重启
	stop	停止
`)
		os.Exit(0)
	}

	pid = os.Getpid()
	cmd := os.Args[1]
	log.Info("[", pid, "]", cmd)

	switch cmd {
	case "stop":
		log.Info("[", pid, "]", "停止...")
		b, _ := ioutil.ReadFile("pid.txt")
		pid, _ := strconv.Atoi(string(b))
		_ = syscall.Kill(pid, syscall.SIGKILL)
		os.Exit(0)
	case "restart":
		log.Info("[", pid, "]", "重启...")
		b, _ := ioutil.ReadFile("pid.txt")
		pid, _ := strconv.Atoi(string(b))
		_ = syscall.Kill(pid, syscall.SIGKILL)
	case "start":
		log.Info("[", pid, "]", "启动...")
	default:
		log.Info("[", pid, "]", "命令未定义，操作已忽略")
		os.Exit(0)
	}
	f, err := os.OpenFile("pid.txt", os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Info("[", pid, "]", "无法创建pid.txt文件")
		os.Exit(0)
	}
	_, _ = f.WriteString(strconv.Itoa(pid))
	_ = f.Close()

	log.Info("[", pid, "]", "加载配置文件：", confFilename)
	loadConf()
	log.Info("[", pid, "]", *conf)

}
func loadConf() {
	n, err := ioutil.ReadFile(confFilename)
	if err != nil {
		panic(err)
	}
	_ = json.Unmarshal(n, conf)
}
func do(conf *Conf) {
	wg.Add(3)
	go add(conf)
	go set(conf)
	go del(conf)
	wg.Wait()
}
func add(conf *Conf) {

	for {
		doAdd(conf)
	}

	wg.Done()

}
func doAdd(conf *Conf) {

	defer func() {
		if err := recover(); err != nil {
			log.Info("错误", err)
		}
	}()

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
	id := rid.Gen()
	tm := datetime.Get()
	statement := "INSERT INTO %s (create_time,rid,name) VALUES (?, ?, ?)"
	parameter := []interface{}{tm, id, "name_" + strconv.Itoa(id)}

	write("add", conf, statement, parameter, id)

}
func set(conf *Conf) {

	for {
		doSet(conf)
	}

	wg.Done()

}
func doSet(conf *Conf) {

	defer func() {
		if err := recover(); err != nil {
			log.Info("[", pid, "]", "捕获错误", err)
		}
	}()

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
	count := time.Now().Unix()
	tm := datetime.Get()
	id := rid.GetByRandom()
	statement := "UPDATE %s SET update_time = ?, count = ? WHERE rid = ?"
	parameter := []interface{}{tm, count, id}

	write("set", conf, statement, parameter, id)

}
func del(conf *Conf) {

	for {
		doDel(conf)
	}

	wg.Done()

}
func doDel(conf *Conf) {

	defer func() {
		if err := recover(); err != nil {
			log.Info("[", pid, "]", "捕获错误", err)
		}
	}()

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
	id := rid.GetByRandom()
	tm := datetime.Get()
	statement := "UPDATE %s SET update_time = ?, delete_time = ?, status = ? WHERE rid = ?"
	parameter := []interface{}{tm, tm, 0, id}

	write("del", conf, statement, parameter, id)

}
func write(action string, conf *Conf, statement string, parameter []interface{}, id int) {

	if conf.IsReadOnly == 1 {
		log.Info("[", pid, "]", "服务降级项目只读", action, "old", statement, parameter)
		return
	}
	if conf.W == 1 {
		// 1写新表
		writeNew(action, statement, tbBaseName, parameter, id, dbCount, tbCount)
	} else if conf.W == 2 {
		// 2写旧表写新表
		writeOld(action, statement, tbBaseName, parameter)
		writeNew(action, statement, tbBaseName, parameter, id, dbCount, tbCount)
	} else {
		// 默认写旧表
		writeOld(action, statement, tbBaseName, parameter)
	}

}
func writeOld(action string, sql string, tbBaseName string, parameter []interface{}) {

	statement := fmt.Sprintf(sql, tbBaseName)
	dao.ExecuteWithParameter(db, statement, parameter)

	log.Info("[", pid, "]", action, "old", statement, parameter)

}
func writeNew(action string, sql string, tbBaseName string, parameter []interface{}, id int, dbCount int, tbCount int) {

	m := id % (dbCount * tbCount)
	dbIndex := m / tbCount
	tbIndex := m % tbCount
	tbName := tbBaseName + "_" + strconv.Itoa(tbIndex)
	statement := fmt.Sprintf(sql, tbName)

	dao.ExecuteWithParameter(db, statement, parameter)
	log.Info("[", pid, "]", action, "new", statement, parameter, id, m, dbIndex, tbIndex)

}
