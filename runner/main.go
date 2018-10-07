package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ima8/booking-ticket/model/ticket"
)

var baseAPI = "http://localhost:8000"
var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getRemainTicket() []string {
	data := new(ticket.TicketRemain)
	err := getJson(baseAPI+"/remaining", data)
	if err != nil {
		log.Fatalln(err)
	}
	return data.Seats
}
func bookTicket(seat string) bool {
	data := url.Values{}
	data.Set("seat", seat)
	_, err := http.Post(baseAPI+"/book", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(seat)
	}
	return true
}
func confirmTicket(seat string) bool {
	data := url.Values{}
	data.Set("seat", seat)
	_, err := http.Post(baseAPI+"/confirm", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(seat)
	}
	return true
}

func bookTicketAndConfirm(seat string, wg *sync.WaitGroup) {
	defer wg.Done()
	status := bookTicket(seat)
	if status == true {
		confirmTicket(seat)
	}
}

func startBookTicketAndConfirm() {
	var wg sync.WaitGroup
	seats := getRemainTicket()
	wg.Add(len(seats))
	for _, seat := range seats {
		go bookTicketAndConfirm(seat, &wg)
	}
	wg.Wait()
	startBookTicketAndConfirm()
}
func main() {
	forever := make(chan bool)
	startBookTicketAndConfirm()
	<-forever
}
