package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
	"tool/dao"
	"tool/datetime"
	"tool/file"
	"tool/in"
	"tool/rid"
)

var (
	wg                 sync.WaitGroup // 控制多个协程的退出
	db                 *sql.DB        // 数据库连接池操作句柄
	pid                int            // 当前进程PID
	pidFilename        = "pid.txt"    // 当前进程PID保存文件
	sleep              = 500000       //每个协程随机休眠最大微妙数
	redisHost          = "127.0.0.1"  // Redis服务器地址
	redisPort          = 6379         // Redis服务端口
	redisPassword      = ""           // Redis密码
	redisDbIndex       = 0            // Redis数据库索引
	tbBaseName         = "tb"         // 基表名称
	tbCount            = 32           // 表数量
	dbDriver           = "mysql"      // 数据库驱动
	dbHost             = "localhost"  // 数据库地址
	dbProtocol         = "tcp"        // 协议
	dbPort             = 3306         // 数据库端口
	dbBusinessUsername = "business"   // 业务账号（实际业务中操作数据库使用的普通账号）
	dbBusinessPassword = "business"   // 业务密码
	dbBaseName         = "db"         // 试验数据库
	dbCount            = 1            // 试验数据库数量
	confFilename       = "app.json"   // 项目配置文件路径
	conf               = &Conf{}      // 项目配置文件解析结果（项目读写控制）
)

type Conf struct {
	W          int `json:"w"`
	R          int `json:"r"`
	IsReadOnly int `json:"is_read_only"`
}

func init() {

	pid = os.Getpid()
	rid.Init(redisHost, redisPort, redisPassword, redisDbIndex)
	db = dao.Connect(dbDriver, fmt.Sprintf("%s:%s@%s(%s:%d)/%s", dbBusinessUsername, dbBusinessPassword, dbProtocol, dbHost, dbPort, dbBaseName))

}
func main() {

	parseCommand()
	savePid()
	loadConf()
	work()

}
func log(message ...interface{}) {
	fmt.Println(datetime.Get(), "[", pid, "]", message)
}
func parseCommand() {

	log("正在解析命令行参数...")

	if len(os.Args) == 1 || !in.StringList(os.Args[1], []string{"start", "restart", "stop"}) {
		fmt.Println("Usage:\n\tstart | restart | stop\n\tstart\t启动\n\trestart\t重启\n\tstop\t停止")
		os.Exit(0)
	}

	cmd := os.Args[1]
	log(cmd)

	switch cmd {
	case "stop":
		log("正在停止...")
		if err := syscall.Kill(file.ReadByInt(pidFilename), syscall.SIGKILL); err != nil {
			panic(err)
		}
		os.Exit(0)
	case "restart":
		log("正在重启...")
		if err := syscall.Kill(file.ReadByInt(pidFilename), syscall.SIGKILL); err != nil {
			panic(err)
		}
	case "start":
		log("正在启动...")
	default:
		log("命令未定义，操作已忽略...")
		os.Exit(0)
	}

}
func savePid() {

	log("正在保存pid...")
	file.SaveByInt(pidFilename, pid, os.ModePerm)
	log("成功")

}
func loadConf() {

	log("正在加载配置...")
	file.ReadByJson(confFilename, conf)
	log("成功")
	log(*conf)

}
func work() {

	wg.Add(3)

	go add()
	go set()
	go del()

	wg.Wait()

}
func add() {

	for {
		doAdd()
	}

}
func set() {

	for {
		doSet()
	}

}
func del() {

	for {
		doDel()
	}

}
func doAdd() {

	defer func() {
		if err := recover(); err != nil {
			log("捕获业务错误", err)
		}
	}()

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
	id := rid.Gen()
	tm := datetime.Get()
	statement := "INSERT INTO %s (create_time,rid,name) VALUES (?, ?, ?)"
	parameter := []interface{}{tm, id, "name_" + strconv.Itoa(id)}

	write("add", statement, parameter, id)

}
func doSet() {

	defer func() {
		if err := recover(); err != nil {
			log("捕获业务错误", err)
		}
	}()

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
	count := time.Now().Unix()
	tm := datetime.Get()
	id := rid.GetByRandom()
	statement := "UPDATE %s SET update_time = ?, count = ? WHERE rid = ?"
	parameter := []interface{}{tm, count, id}

	write("set", statement, parameter, id)

}
func doDel() {

	defer func() {
		if err := recover(); err != nil {
			log("捕获业务错误", err)
		}
	}()

	time.Sleep(time.Microsecond * time.Duration(rand.Intn(sleep)))
	id := rid.GetByRandom()
	tm := datetime.Get()
	statement := "UPDATE %s SET update_time = ?, delete_time = ?, status = ? WHERE rid = ?"
	parameter := []interface{}{tm, tm, 0, id}

	write("del", statement, parameter, id)

}
func write(action string, statement string, parameter []interface{}, id int) {

	if conf.IsReadOnly == 1 {
		panic([]interface{}{"服务降级项目只读", action, "old", statement, parameter})
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

	log(action, "old", statement, parameter)

}
func writeNew(action string, sql string, tbBaseName string, parameter []interface{}, id int, dbCount int, tbCount int) {

	m := id % (dbCount * tbCount)
	dbIndex := m / tbCount
	tbIndex := m % tbCount
	tbName := tbBaseName + "_" + strconv.Itoa(tbIndex)
	statement := fmt.Sprintf(sql, tbName)

	dao.ExecuteWithParameter(db, statement, parameter)
	log(action, "new", statement, parameter, id, m, dbIndex, tbIndex)

}
