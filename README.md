阶段2：单库多表
====
将原有大表拆分成小表
拆分规则：
将大表拆分成n（n=32）个小表
中间变量：路由主键kid % (库数量 * n)
库：中间变量 / n
表：中间变量 % n
旧数据库：db
旧数据表：tb
数据库名保持不变，在数据库内部对大表进行拆分
tb => tb_n(n∈[0,31])
称tb为基表，_n通过路由规则计算得出，一般数据库配置只需要指定基表的配置即可

步骤
====
创建数据表:
tb_n(n∈[0,31])
修改配置：
基表：
启用路由：是
写：单表（指旧表），多表（指新表）
读：单表
依赖：
业务代码调用DAO时传入kid（提前做）
框架底层对基表的所有写操作（insert,update,delete,replace等）在原有写逻辑上增加按照新路由规则双写语句（将待写入的SQL复制一条，修改表名，一起执行）
结果：
insert写会双写
update,delete,replace写对于不在多表中的数据只会写入单表，对于同时存在单表和多表的数据会双写
后台进程：
对基表采用主键ID降序分页批量读取数据：
①如果数据在多表中不存在，则执行插入
②如果数据在多表中存在且数据update_time>多表数据的update_time，则执行覆盖
执行2遍
结果：
完成tb=>tb_n(n∈[0,31])数据迁移
单表和多表数据完全一致
如果不一致，则继续同步，继续验证，直到数据完全一致
修改配置：
基表：
启用路由：是
写：单表，多表
读：多表
结果：
初步完成迁移切换，数据在单表进行冗余，以备回滚
验证逻辑
此阶段理论上不应该出现问题（数据一致性在上一步骤得到验证，路由规则由框架提供）
修改配置：
基表：
启用路由：是
写：多表
读：多表
结果：
读写均操作多表，完成大表拆分任务
后续操作：
删除单表

实验
====
创建数据库：db
```sql
DROP DATABASE IF EXISTS `db`;
CREATE DATABASE `db`;
```
创建基表：tb
```sql
CREATE TABLE `tb` (
`id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
`add_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '创建时间',
`set_time` datetime NOT NULL DEFAULT '1970-01-01 00:00:00' COMMENT '修改时间',
`kid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '路由主键',
`name` varchar(255) NOT NULL DEFAULT '' COMMENT '文本信息',
`count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '模拟请求',
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```
创建模拟数据基：1000000
```sql
LOAD DATA INFILE '/Users/ziji/GolandProjects/mysql_sharding/step0/tb.csv'
INTO TABLE tb
FIELDS
TERMINATED BY ','
ENCLOSED BY '"'
IGNORE 1 LINES
```


大表拆分时迁移数据方案：

旧数据
----
1.脚本直接批量读取tb数据，然后逐行按照路由规则，写入新表。此方案效率太低，几乎不可用（500w数据，耗时约几十个小时，性能瓶颈主要在写）
2.脚本直接批量读取tb数据，然后逐行按照路由规则，写入txt文件，然后load到新表，速度可以。但是中间多了一层中转，提升了复杂度
3.脚本（可以完全不需要脚本，直接将SQL放在数据库上执行）直接采用insert ignore into ... select ... 语句，直接在MySQL内部批量筛选、计算、插入，效率成千上万倍提升，没有额外的复杂度
----
方案3，对于500w数据，耗时2分钟左右。当然，模拟数据比较简单，实际耗时会比实验多，但此方案速度最快，效率最高，需要在业务量低谷时段操作。
----
服务降级，只读不写，然后将所有旧表update_time大于等于迁移开始时间的数据，全部覆盖或插入新表中
启用双写，校验数据集
将读改为新表，验证逻辑
关闭写旧表




