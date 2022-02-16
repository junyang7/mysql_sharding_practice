module step1

go 1.17

require (
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
)

require tool v0.0.0

replace tool => ./../../tool
