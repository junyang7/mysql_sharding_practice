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


