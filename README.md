# mysql_sharding_practice

MySQL-分库分表实践

----

随着业务的增长，单库单表已经不能满足数据服务的需求，我们需要及时的进行分库分表操作来提升数据库服务的高可用。

## 目录结构

```text
.
├── README.md   说明文件
├── business    模拟业务
├── doc         文档
├── install.sh  安装脚本
├── test        具体测试实施方案
└── tool        常用工具
```

## 实践思路

常见演进过程：单库单表-》单库多表-》多库多表

### 拆表

将单库中的大表按照一定规则拆分成多个小表

### 分库

将一个库按照一定规则拆分成多个库，每个库的表数量和表结构完全一致，通常每个库独占一个实例
