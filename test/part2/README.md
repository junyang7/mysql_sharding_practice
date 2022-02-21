# 阶段：单库单表-1主n从

----
数据库名（已有数据）：db
主库1（读写）：172.16.10.30
从库2（只读）：172.16.10.27
从库3（只读）：172.16.10.28

## 安装数据库

```shell
setenforce 0
systemctl stop firewalld
systemctl disable firewalld
curl -O https://repo.mysql.com//mysql80-community-release-el7-5.noarch.rpm
yum -y install mysql80-community-release-el7-5.noarch.rpm
yum install -y --downloadonly --downloaddir=. mysql-community-server
rpm -Uvh --force --nodeps *.rpm
systemctl enable mysqld
systemctl start mysqld
grep 'temporary password' /var/log/mysqld.log
mysql -uroot -p

```

## 初始化配置

```sql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'aA!12345';
USE mysql;
UPDATE user SET host = '%' WHERE user = 'root';
FLUSH PRIVILEGES;
SHOW VARIABLES LIKE '%uuid%';
SHOW VARIABLES LIKE '%server_id%';
SHOW VARIABLES LIKE '%log_bin%';
exit;

```

## master配置

```shell
echo '
read_only=off
server-id=1
binlog-do-db=db
' >> /etc/my.cnf
systemctl restart mysqld
```

## slaver配置

```shell
echo '
read_only=on
server-id=2
replicate-do-db=db
report-host=172.16.10.27
' >> /etc/my.cnf
systemctl restart mysqld

echo '
read_only=on
server-id=3
replicate-do-db=db
report-host=172.16.10.28
' >> /etc/my.cnf
systemctl restart mysqld
```

## 实施过程

### 业务链接主库，正常模拟读写

### master机器操作

```sql
mysql -uroot -p

flush tables with read lock;
show master status;

```

```shell
mysqldump -uroot -p -q db > db.sql
```

```sql
mysql -uroot -p

unlock tables;

```

```shell
scp -r db.sql root@172.16.10.27:/tmp/db.sql
```
```shell
scp -r db.sql root@172.16.10.28:/tmp/db.sql
```

### slaver机器操作

```sql
mysql -uroot -p

CREATE DATABASE IF NOT EXISTS db DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;
exit;
```
```shell
mysql -uroot -p db < /tmp/db.sql
```
```sql
mysql -uroot -p

reset slave;
change master to master_host='172.16.10.30', master_user='root',master_password='aA!12345', master_log_file='binlog.000004', master_log_pos=4592195;
start slave;
show slave status \G;

```


## 验证

```shell

export http_proxy="http://172.16.9.31:4780"
export https_proxy="http://172.16.9.31:4780"

yum -y install https://mirrors.tuna.tsinghua.edu.cn/percona/yum/percona-release-1.0-17.noarch.rpm
yum install percona-toolkit -y


# 一次只能校验一个主库和一个从库，多个从库无法同时校验
pt-table-checksum --nocheck-replication-filters --no-check-binlog-format --replicate=db.checksums h=127.0.0.1,u=root,p='aA!12345',P=3306 --databases=db



mysql8不被支持
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'aA!12345';
FLUSH PRIVILEGES;
```

从库修复
```shell
pt-table-sync --sync-to-master h=172.16.10.27,u=root,p='aA!12345',P=3306 --print
```

