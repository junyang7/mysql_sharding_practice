# 业务模拟

主要模拟数据库的增伤改写操作，模拟过程尽可能的增加并发，比如同时部署多个服务

----

## 目录结构

```text
.
├── README.md   说明文件
├── app.json    项目配置文件
├── go.mod
└── main.go     模拟代码
```

## 配置文件

```json5
{
  // 业务写控制：0写旧表，1写新表，2写旧表写新表
  "w": 0,
  // 业务读控制：0读旧表，1读新表
  "r": 0
}
```

## 命令参考

```shell
# 启动
go build main.go && nohup ./main start >> log.txt 2>&1 &

# 重启
go build main.go && nohup ./main restart >> log.txt 2>&1 &

# 停止
go build main.go && nohup ./main stop >> log.txt 2>&1 &

# 日志
tail -f log.txt
```
