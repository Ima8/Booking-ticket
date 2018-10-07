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

// GetRemainTicket is a function for get Remain Ticket of current round
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

// BookTicket is a function for booking the ticket of current round
func BookTicket(seat string) bool {
	clientRedis, _ = redisConnector.GetConnection(0)
	currentRound := getCurrentRound()
	if currentRound == 0 {
		return false
	}
	canBookTicket := isTicketAvailable(currentRound, seat)
	if canBookTicket == true {
		status, err := clientRedis.RenameNX("r_"+strconv.Itoa(currentRound)+":"+seat, "r_"+strconv.Itoa(currentRound)+"_u:"+seat).Result()
		if status == false || err != nil {
			return false
		} else {
			log.Println("Booked: " + seat)
			return true
		}
	}
	return false

}

func getCurrentRound() int {
	clientRedis, _ = redisConnector.GetConnection(0)
	round, _ := clientRedis.Get("current_round").Result()
	var currentRound int
	if len(round) >= 1 {
		currentRound, _ = strconv.Atoi(round)
	}
	return (currentRound)
}

func isTicketAvailable(currentRound int, s string) bool {
	clientRedis, _ = redisConnector.GetConnection(0)
	status, err := clientRedis.Exists("r_" + strconv.Itoa(currentRound) + ":" + s).Result()
	if status == 0 || err != nil {
		log.Printf("Seat %s, Round %d:Not found", s, currentRound)
		return false
	}
	if status == 1 {
		return true
	}
	return false
}
