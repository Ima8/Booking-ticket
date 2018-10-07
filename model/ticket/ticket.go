package ticket

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	if keys, _, err = clientRedis.Scan(0, "ru_"+strconv.Itoa(round)+":*", 1000).Result(); err != nil {
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
		//status, err := clientRedis.RenameNX("r_"+strconv.Itoa(currentRound)+":"+seat, "ru_"+strconv.Itoa(currentRound)+":"+seat).Result()
		_, err := clientRedis.Set("ru_"+strconv.Itoa(currentRound)+":"+seat, "", 10*time.Second).Result()
		if err != nil {
			return false
		} else {
			log.Println("Booked: " + seat)
			// Check is it still have ticket left if not init the new round
			return true
		}
	}
	return false

}

// isRoundAvailable mean all ticket already booked, may have some not confirm yet
func isRoundAvailable(round int) bool {
	// clientRedis, _ = redisConnector.GetConnection(0)
	// client
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
	isAvailable, err := clientRedis.Exists("r_" + strconv.Itoa(currentRound) + ":" + s).Result()
	isBooked, err2 := clientRedis.Exists("ru_" + strconv.Itoa(currentRound) + ":" + s).Result()
	if isAvailable == 1 || err != nil || isBooked == 0 || err2 != nil {
		log.Printf("Seat %s, Round %d:Not found", s, currentRound)
		return false
	}
	if isAvailable == 1 && isBooked == 0 {
		return true
	}
	return false
}
