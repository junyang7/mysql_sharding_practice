business    模拟业务逻辑
test        试验
tool        工具

----
流程
----
准备环境：数据库，数据表，数据集，账号，业务
执行业务（起2个服务），模拟增删改逻辑（要求存在并发）
启动分表脚本
修改配置为【读旧表，写旧表写信表，读写】，热重启业务（要求服务不停机，2个服务错开重启）
启动迁移脚本

---
测试1：数据库通过read_only控制读写
---
无序处理，迁移脚本已内嵌，默认60次内如果数据没有对齐则进行read_only操作


---
测试2：业务框架服务降级控制读写
---
第一次完成迁移后： 修改配置为【读旧表，写旧表写信表，只读】，热重启业务（要求服务不停机，2个服务错开重启）
迁移完毕恢复读写： 修改配置为【读旧表，写旧表写信表，读写】，热重启业务（要求服务不停机，2个服务错开重启）


----
迁移完成判断条件：
----
时差范围内，从旧表读取数据，判断是否需要处理，如果待处理数据条数是0，则初步判断数据已对齐
数据对齐后，开始执行校验程序，从旧表中，按照路由规则，逐个验证：旧表left join新表on所有字段（除了自增ID外），形成大表，判断旧表ID是null的个数是否是0。如果是0，则数据一致。
验证通过，迁移脚本退出。


----
脚本
----
```shell
cd business
go get github.com/go-redis/redis
go install github.com/go-redis/redis
go get github.com/go-sql-driver/mysql
go install github.com/go-sql-driver/mysql
cd ..
cd tool
go get github.com/go-redis/redis
go install github.com/go-redis/redis
go get github.com/go-sql-driver/mysql
go install github.com/go-sql-driver/mysql
cd ..
cd test
cd step0
go get github.com/go-redis/redis
go install github.com/go-redis/redis
go get github.com/go-sql-driver/mysql
go install github.com/go-sql-driver/mysql
cd ..
cd step1
go get github.com/go-redis/redis
go install github.com/go-redis/redis
go get github.com/go-sql-driver/mysql
go install github.com/go-sql-driver/mysql
cd ..
cd step2
go get github.com/go-redis/redis
go install github.com/go-redis/redis
go get github.com/go-sql-driver/mysql
go install github.com/go-sql-driver/mysql
cd ..
cd ..
```
