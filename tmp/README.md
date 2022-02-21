```shell

# 初始化环境
chmod +x ./step0
./step0

# 启动业务
echo '
{
  "w": 0,
  "r": 0,
  "is_read_only": 0
}
' > app.json
chmod +x ./business
nohup ./business start >> log.txt 2>&1 &

# 建表
chmod +x ./step1
./step1

# 双写
echo '
{
  "w": 2,
  "r": 0,
  "is_read_only": 0
}
' > app.json
nohup ./business restart >> log.txt 2>&1 &

# 同步
chmod +x ./step2
./step2

# 单写
echo '
{
  "w": 1,
  "r": 1,
  "is_read_only": 0
}
' > app.json
nohup ./business restart >> log.txt 2>&1 &
```

```shell
tail -f log.txt

```
