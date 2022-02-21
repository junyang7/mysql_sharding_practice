阶段：模拟单库单表环境
====

----
流程
----

创建模拟数据集：5000000行

设置分布式唯一ID初始值：5000000

创建数据库：db

创建数据表：tb

将模拟数据导入到db.tb中

创建业务账号：business，业务框架中操作数据库使用，这个账号可以有DBA提供，只需要表内数据的增删改查即可

创建迁移账号：transfer，迁移脚本使用，可以直接将脚本托管给DBA来运行，要求具有super(read_only=1时不受影响),create(建表),表内数据增删改查权限

----
准备
----

数据库服务器：mysql/8.0.28/127.0.0.1:3306

超级管理员账号密码：root/root，拥有create user/grant option/insert/delete/update/select/create/super权限

数据库名：db

数据表名：tb

数据行数：500w

导入导出路径：任意

----
参考
----

```shell
echo "
[mysqld]
secure_file_priv=
" > ~/.my.cnf


echo 'secure_file_priv=' >> /etc/my.cnf
```

```sql
SHOW VARIABLES LIKE "%secure_file_priv%";
```

```sql
CREATE TABLE `tb` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',
  `delete_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '删除时间',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态：1未删除，0已删除',
  `rid` bigint unsigned NOT NULL DEFAULT '0' COMMENT '分布式唯一ID',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '文本信息',
  `count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '模拟请求',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rid` (`rid`) USING BTREE COMMENT '分布式唯一ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
```

----
脚本
----
```shell
go build main.go && ./main
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /Users/guoguo/GolandProjects/mysql_sharding_practice/tmp/step0 main.go

```
