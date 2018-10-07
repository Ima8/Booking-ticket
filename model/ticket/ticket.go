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
		totalUncon = totalUncon + 1
	}
	// Get all key that didn't book
	if keys, _, err = clientRedis.Scan(0, "r_"+strconv.Itoa(round)+":*", 1000).Result(); err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(2)
	}
	for i := 0; i < len(keys); i++ {
		parts := strings.Split(keys[i], ":")
		remainTicket = append(remainTicket, (parts[1]))
		totalUncon = totalUncon + 1
	}
	return remainTicket, totalUncon

}

func BookTicket(seat string) bool {
	clientRedis, _ = redisConnector.GetConnection(0)
	//
	IsTicketAvailable(seat)
	return false
}

func GetCurrentRound() int {
	clientRedis, _ = redisConnector.GetConnection(0)
	round, _ := clientRedis.Get("current_round").Result()
	var currentRound int
	if len(round) >= 1 {
		currentRound, _ = strconv.Atoi(round)
	}
	return (currentRound)
}
func IsTicketAvailable(s string) bool {
	clientRedis, _ = redisConnector.GetConnection(0)
	currentRound := GetCurrentRound()
	if currentRound == 0 {
		return false
	}
	_, err := clientRedis.Get("r_1" + strconv.Itoa(currentRound) + ":" + s).Result()
	if err != nil {
		log.Printf("Seat %s, Round %d:Not found", s, currentRound)
		return false
	}
	return true
}
