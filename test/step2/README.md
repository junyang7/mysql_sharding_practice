阶段：大表拆分-同步
====

----
初始迁移
----

程序启动后，记录当前时间A，设置起始时间B="1970-01-01 00:00:00"

从旧表中按路由规则筛选上述时段范围数据，采用insert ignore into...语句新表中

此部分迁移量能达到99%

----
差异迁移：不存在update，因为迁移有行锁，迁移数据一致性能得到保证，时间以外的数据只会add
----

初始迁移后，程序进入无限循环，每次都拿当前时间为结束时间，之前的结束时间位本次起始时间，筛选旧表修改时间在此范围内的数据

逐条在新表进行判断：存在且时间大于新表，则更新；不存在，则插入

每次循环，都记录本次循环处理的数据条数

如果待处理数据条数=0，说明不存在差异，已完成数据补齐

----
数据校验
----

按照路由规则，旧表筛选数据left join新表on所有字段（除了自增ID），判断旧表ID是NULL的个数是否是0

是0，则数据集一致，程序退出，否则不一致，如果不一致，则继续进行差异迁移

----
参考
----

```sql
SELECT 0 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 0) AS old, (SELECT COUNT(*) FROM tb_0) AS new UNION ALL
SELECT 1 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 1) AS old, (SELECT COUNT(*) FROM tb_1) AS new UNION ALL
SELECT 2 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 2) AS old, (SELECT COUNT(*) FROM tb_2) AS new UNION ALL
SELECT 3 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 3) AS old, (SELECT COUNT(*) FROM tb_3) AS new UNION ALL
SELECT 4 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 4) AS old, (SELECT COUNT(*) FROM tb_4) AS new UNION ALL
SELECT 5 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 5) AS old, (SELECT COUNT(*) FROM tb_5) AS new UNION ALL
SELECT 6 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 6) AS old, (SELECT COUNT(*) FROM tb_6) AS new UNION ALL
SELECT 7 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 7) AS old, (SELECT COUNT(*) FROM tb_7) AS new UNION ALL
SELECT 8 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 8) AS old, (SELECT COUNT(*) FROM tb_8) AS new UNION ALL
SELECT 9 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 9) AS old, (SELECT COUNT(*) FROM tb_9) AS new UNION ALL
SELECT 10 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 10) AS old, (SELECT COUNT(*) FROM tb_10) AS new UNION ALL
SELECT 11 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 11) AS old, (SELECT COUNT(*) FROM tb_11) AS new UNION ALL
SELECT 12 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 12) AS old, (SELECT COUNT(*) FROM tb_12) AS new UNION ALL
SELECT 13 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 13) AS old, (SELECT COUNT(*) FROM tb_13) AS new UNION ALL
SELECT 14 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 14) AS old, (SELECT COUNT(*) FROM tb_14) AS new UNION ALL
SELECT 15 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 15) AS old, (SELECT COUNT(*) FROM tb_15) AS new UNION ALL
SELECT 16 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 16) AS old, (SELECT COUNT(*) FROM tb_16) AS new UNION ALL
SELECT 17 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 17) AS old, (SELECT COUNT(*) FROM tb_17) AS new UNION ALL
SELECT 18 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 18) AS old, (SELECT COUNT(*) FROM tb_18) AS new UNION ALL
SELECT 19 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 19) AS old, (SELECT COUNT(*) FROM tb_19) AS new UNION ALL
SELECT 20 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 20) AS old, (SELECT COUNT(*) FROM tb_20) AS new UNION ALL
SELECT 21 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 21) AS old, (SELECT COUNT(*) FROM tb_21) AS new UNION ALL
SELECT 22 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 22) AS old, (SELECT COUNT(*) FROM tb_22) AS new UNION ALL
SELECT 23 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 23) AS old, (SELECT COUNT(*) FROM tb_23) AS new UNION ALL
SELECT 24 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 24) AS old, (SELECT COUNT(*) FROM tb_24) AS new UNION ALL
SELECT 25 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 25) AS old, (SELECT COUNT(*) FROM tb_25) AS new UNION ALL
SELECT 26 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 26) AS old, (SELECT COUNT(*) FROM tb_26) AS new UNION ALL
SELECT 27 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 27) AS old, (SELECT COUNT(*) FROM tb_27) AS new UNION ALL
SELECT 28 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 28) AS old, (SELECT COUNT(*) FROM tb_28) AS new UNION ALL
SELECT 29 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 29) AS old, (SELECT COUNT(*) FROM tb_29) AS new UNION ALL
SELECT 30 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 30) AS old, (SELECT COUNT(*) FROM tb_30) AS new UNION ALL
SELECT 31 AS `index`, (SELECT COUNT(*) FROM tb WHERE rid % (1 * 32) % 32 = 31) AS old, (SELECT COUNT(*) FROM tb_31) AS new
```


```sql
SELECT
    COUNT(*) AS `c`
FROM
    (
        SELECT
            `%s`.`id` AS `%s_id`,
            `%s`.`id` AS `%s_id`
        FROM
            `%s`
        LEFT JOIN
            `%s`
                ON
                    `%s`.`create_time` = `%s`.`create_time` AND
                    `%s`.`update_time` = `%s`.`update_time` AND
                    `%s`.`delete_time` = `%s`.`delete_time` AND
                    `%s`.`status` = `%s`.`status`  AND
                    `%s`.`rid` = `%s`.`rid` AND
                    `%s`.`name` = `%s`.`name` AND
                    `%s`.`count` = `%s`.`count`
        WHERE
            `%s`.`rid` %% (1 * 32) %% 32 = %d
    ) AS `t`
WHERE
    `t`.`%s_id` IS NULL
;
```
````text
tbBaseName,tbBaseName,
tbName,tbName,
tbBaseName,
tbName,
tbBaseName,tbName,
tbBaseName,tbName,
tbBaseName,tbName,
tbBaseName,tbName,
tbBaseName,tbName,
tbBaseName,tbName,
tbBaseName,tbName,
tbBaseName,j,
tbName,
````

----
脚本
----
```shell
go build main.go && ./main

```
