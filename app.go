package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ima8/Booking-ticket/initalTicket"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	redisConnector "github.com/ima8/booking-ticket/model/redis"
	"github.com/ima8/booking-ticket/model/ticket"
	"github.com/spf13/viper"
)

var startTime time.Time
var clientRedis *redis.Client
var remainDB int

type ResponseTicket struct {
	Success           bool   `json:"success"`
	Seat              string `json:"seat"`
	ReserveExiredTime string `json:"reservedExpiredTime"`
}

type ResponseTicketWithoutReserveExiredTime struct {
	Success bool   `json:"success"`
	Seat    string `json:"seat"`
}
type ResponseWithAllTicket struct {
	Success bool     `json:"success"`
	Seats   []string `json:"seat"`
}

func uptime() time.Duration {
	return time.Since(startTime)
}

func loadConf() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("\n Fatal error config file: %s ", err))
	}
	remainDB = viper.GetInt("remain_DB")
}

func init() {
	startTime = time.Now()
	loadConf()
	clientRedis, _ = redisConnector.GetConnection(remainDB)
}

// HomeHandler is a health check api
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Seat string
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
	}
	p := new(requestData)
	decoder := schema.NewDecoder()
	err = decoder.Decode(p, r.Form)
	if err != nil {
		fmt.Println("Error decoding")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing parameter")
		return
	}
	isBook := ticket.ConfirmTicket(p.Seat)
	w.Header().Set("Content-Type", "application/json")
	if isBook == true {
		response := ResponseTicketWithoutReserveExiredTime{
			Success: true,
			Seat:    p.Seat,
		}
		b, err := json.Marshal(&response)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, string(b))

	} else {
		response := ResponseTicketWithoutReserveExiredTime{
			Success: false,
			Seat:    p.Seat,
		}
		b, err := json.Marshal(&response)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, string(b))

	}
}

// BookHandler is API For booking the ticket
func BookHandler(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Seat string
	}
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
	}
	p := new(requestData)
	decoder := schema.NewDecoder()
	err = decoder.Decode(p, r.Form)
	if err != nil {
		fmt.Println("Error decoding")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing parameter")
		return
	}
	isBook := ticket.BookTicket(p.Seat)
	w.Header().Set("Content-Type", "application/json")
	if isBook == true {
		t := time.Now()
		t = t.Add(10 * time.Second)
		response := ResponseTicket{
			Success:           true,
			Seat:              p.Seat,
			ReserveExiredTime: t.Format(time.RFC850),
		}
		b, err := json.Marshal(&response)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, string(b))

	} else {
		response := ResponseTicket{
			Success:           false,
			Seat:              p.Seat,
			ReserveExiredTime: "",
		}
		b, err := json.Marshal(&response)
		if err != nil {
			fmt.Println(err)
			return
		}
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, string(b))

	}
}

// RemainHandler is API for get the remain ticket and number of unconfirm ticket
func RemainHandler(w http.ResponseWriter, r *http.Request) {
	currentRound := ticket.GetCurrentRound()
	if currentRound == 0 {
		fmt.Fprintf(w, "")
		return
	}

	remainTicket, totalUncon := ticket.GetRemainTicket(currentRound)
	remainData := ticket.TicketRemain{
		UnconfimedTicketsCount: totalUncon,
		Round: currentRound,
		Seats: remainTicket,
	}
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(&remainData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}

// InitHandler is API for inital ticket for first round
func InitHandler(w http.ResponseWriter, r *http.Request) {
	initalTicket.InitTicket(1)
	fmt.Fprintf(w, "DONE")
}

// AllTicketHandler is API for get all ticket
func AllTicketHandler(w http.ResponseWriter, r *http.Request) {
	tickets := ticket.GetAllTicket()
	response := ResponseWithAllTicket{
		Success: true,
		Seats:   tickets,
	}
	b, err := json.Marshal(&response)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(b))

}
func main() {
	fmt.Println("Server starting at :8000")
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/remaining", RemainHandler).Methods("GET")
	r.HandleFunc("/book", BookHandler).Methods("POST")
	r.HandleFunc("/confirm", ConfirmHandler).Methods("POST")
	r.HandleFunc("/init", InitHandler).Methods("GET")
	r.HandleFunc("/all", AllTicketHandler).Methods("GET")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
