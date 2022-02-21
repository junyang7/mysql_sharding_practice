阶段：大表拆分-清理
====

----
流程
----

删除旧表

----
脚本
----
```shell
go build main.go && ./main
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /Users/guoguo/GolandProjects/mysql_sharding_practice/tmp/step3 main.go

```
