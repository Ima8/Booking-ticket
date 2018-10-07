package initalTicket

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	redisConnector "github.com/ima8/booking-ticket/model/redis"
)

var startTime time.Time
var clientRedis *redis.Client

func init() {
	startTime = time.Now()
	clientRedis, _ = redisConnector.ConnectRedisServer(0)
}

func toCharStrArr(arr [26]string, i int) string {
	return arr[i-1]
}

func InitTicket(round int) {
	var arr = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	for i := 1; i < 6; i++ {
		for j := 1; j <= 20; j++ {
			var ticket = toCharStrArr(arr, i) + strconv.Itoa(j)
			clientRedis.Set("r_"+strconv.Itoa(round)+"_u:"+ticket, "", 10*time.Hour)
		}
	}

	for i := 6; i < 8; i++ {
		for j := 1; j <= 20; j++ {
			var ticket = toCharStrArr(arr, i) + strconv.Itoa(j)
			clientRedis.Set("r_"+strconv.Itoa(round)+"_u:"+ticket, "", 10*time.Hour)
		}
	}
}
