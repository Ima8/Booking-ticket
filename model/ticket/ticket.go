package ticket

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	redisConnector "github.com/ima8/booking-ticket/model/redis"
)

type TicketRemain struct {
	Seats                  []string `json:"seats"`
	UnconfimedTicketsCount int      `json:"unconfimedTicketsCount"`
}

var clientRedis *redis.Client

func GetRemainTicket(round int) ([]string, int) {
	//defer clientRedis.Close()
	var remainTicket []string
	var totalUncon int = 0

	clientRedis, _ = redisConnector.ConnectRedisServer(0)

	// var cursor uint64
	var err error
	var keys []string

	// Get all key that booked but didn't confirm yet
	if keys, _, err = clientRedis.Scan(0, "r_"+strconv.Itoa(round)+"_u:*", 1000).Result(); err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(2)
	}

	for i := 0; i < len(keys); i++ {
		parts := strings.Split(keys[i], ":")
		log.Println(parts[1])
		remainTicket = append(remainTicket, (parts[1]))
		totalUncon = totalUncon + 1
	}
	// Get all key that didn't book
	if keys, _, err = clientRedis.Scan(0, "r_"+strconv.Itoa(round)+":*", 1000).Result(); err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(2)
	}
	for i := 0; i < len(keys); i++ {
		parts := strings.Split(keys[i], ":")
		log.Println(parts[1])
		remainTicket = append(remainTicket, (parts[1]))

	}
	return remainTicket, totalUncon

}

func IsTicketFinish(round int) bool {
	var isFinish = true

	return isFinish
}
