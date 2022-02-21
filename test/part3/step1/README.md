阶段：大表拆分-准备
====

----
流程
----

按照配置创建数据表：tb_n(n∈[0,31])

----
脚本
----
```shell
go build main.go && ./main
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /Users/guoguo/GolandProjects/mysql_sharding_practice/tmp/step1 main.go

```
