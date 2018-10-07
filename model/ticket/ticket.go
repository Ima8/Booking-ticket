package ticket

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Ima8/Booking-ticket/initalTicket"
	"github.com/go-redis/redis"
	redisConnector "github.com/ima8/booking-ticket/model/redis"
)

type TicketRemain struct {
	Seats                  []string `json:"seats"`
	UnconfimedTicketsCount int      `json:"unconfimedTicketsCount"`
	Round                  int      `json:"round"`
}

var clientRedis *redis.Client

// GetRemainTicket is a function for get Remain Ticket of current round
func GetRemainTicket(round int) ([]string, int) {
	//defer clientRedis.Close()
	var remainTicket []string
	var totalUncon int = 0

	clientRedis, _ = redisConnector.GetConnection(0)

	// var cursor uint64
	var err error
	var keys []string

	var bookedTicket []string
	// Get all key that booked but didn't confirm yet
	if keys, _, err = clientRedis.Scan(0, "ru_"+strconv.Itoa(round)+":*", 1000).Result(); err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(2)
	}

	for i := 0; i < len(keys); i++ {
		parts := strings.Split(keys[i], ":")
		bookedTicket = append(bookedTicket, (parts[1]))
		totalUncon = totalUncon + 1
	}
	// Get all ticket that didn't confirm
	if keys, _, err = clientRedis.Scan(0, "r_"+strconv.Itoa(round)+":*", 1000).Result(); err != nil {
		fmt.Println("ERROR: %s", err)
		os.Exit(2)
	}
	for i := 0; i < len(keys); i++ {
		parts := strings.Split(keys[i], ":")
		isBooked := false
		for i := range bookedTicket {
			if bookedTicket[i] == parts[1] {
				isBooked = true
				break
			}
		}
		if isBooked == false {
			remainTicket = append(remainTicket, (parts[1]))
			totalUncon = totalUncon + 1
		}
	}
	return remainTicket, totalUncon

}

func ConfirmTicket(seat string) bool {
	isConfirmTicketSuccess := false
	clientRedis, _ = redisConnector.GetConnection(0)
	currentRound := GetCurrentRound()
	if currentRound == 0 {
		return false
	}
	isBooked, err := clientRedis.Exists("ru_" + strconv.Itoa(currentRound) + ":" + seat).Result()
	if isBooked == 1 && err == nil {
		status, err := clientRedis.RenameNX("r_"+strconv.Itoa(currentRound)+":"+seat, "b_"+strconv.Itoa(currentRound)+":"+seat).Result()
		if status == true && err == nil {
			clientRedis.Del("ru_" + strconv.Itoa(currentRound) + ":" + seat)
			log.Println(isRoundFull(currentRound))
			// Check is it still have ticket left if not init the new round
			if isRoundFull(currentRound) {
				initalTicket.InitTicket(currentRound + 1)
			}
			isConfirmTicketSuccess = true
		} else {
			log.Printf("Seat %s, Round %d:Someone Confirm before", seat, currentRound)
			isConfirmTicketSuccess = false
		}
	} else {
		log.Printf("Seat %s, Round %d:Didn't book yet", seat, currentRound)
		isConfirmTicketSuccess = false
	}

	return isConfirmTicketSuccess
}

// BookTicket is a function for booking the ticket of current round
func BookTicket(seat string) bool {
	clientRedis, _ = redisConnector.GetConnection(0)
	currentRound := GetCurrentRound()
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
			return true
		}
	}
	return false

}

// isRoundFull mean all ticket already booked and confirm
func isRoundFull(round int) bool {
	_, unconfirmTicket := GetRemainTicket(round)
	if unconfirmTicket > 0 {
		return false
	}
	return true
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

func GetAllTicket() []string {
	var tickets []string
	clientRedis, _ = redisConnector.GetConnection(0)
	currentRound := GetCurrentRound()
	if currentRound == 0 {
		return []string{}
	}
	for i := 1; i <= currentRound; i++ {
		// var cursor uint64
		var err error
		var keys []string
		strRound := strconv.Itoa(i)
		if keys, _, err = clientRedis.Scan(0, "b_"+strRound+":*", 1000).Result(); err != nil {
			fmt.Println("ERROR: %s", err)
			os.Exit(2)
		}

		for i := 0; i < len(keys); i++ {
			parts := strings.Split(keys[i], ":")
			tickets = append(tickets, strRound+"-"+parts[1])
		}
	}

	return tickets
}
func isTicketAvailable(currentRound int, s string) bool {
	clientRedis, _ = redisConnector.GetConnection(0)
	isAvailable, err := clientRedis.Exists("r_" + strconv.Itoa(currentRound) + ":" + s).Result()
	isBooked, err2 := clientRedis.Exists("ru_" + strconv.Itoa(currentRound) + ":" + s).Result()
	// If ticket available and already booked
	if isAvailable == 1 || err != nil {
		if isBooked == 1 || err2 != nil {
			log.Printf("Seat %s, Round %d:Already Book", s, currentRound)
			return false
		} else {
			return true
		}

	} else {
		log.Printf("Seat %s, Round %d:Not found", s, currentRound)
		return false
	}
}
