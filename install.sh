#!/usr/bin/env bash

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
