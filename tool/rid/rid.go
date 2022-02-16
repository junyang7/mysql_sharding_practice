package rid

import (
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"tool/rd"
)

var r *redis.Client

func Gen() int {
	rid, err := r.Incr("rid").Result()
	if err != nil {
		panic(err)
	}
	return int(rid)
}

func GetByMax() int {
	rid, err := r.Get("rid").Int()
	if err != nil {
		panic(err)
	}
	return rid
}

func GetByRandom() int {
	return rand.Intn(GetByMax())
}

func Init(redisHost string, redisPort int, redisPassword string, redisDbIndex int) {
	r = rd.Connect(redisHost+":"+strconv.Itoa(redisPort), redisPassword, redisDbIndex)
	if err := r.Ping().Err(); err != nil {
		panic(err)
	}
}
