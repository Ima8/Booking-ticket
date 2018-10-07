package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	clientRedis, _ = redisConnector.ConnectRedisServer(remainDB)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
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
	fmt.Println(p.Seat)
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// RemainHandler is API for get the remain ticket and number of unconfirm ticket
func RemainHandler(w http.ResponseWriter, r *http.Request) {
	remainTicket, totalUncon := ticket.GetRemainTicket(1)
	remainData := ticket.TicketRemain{
		Seats: remainTicket,
		UnconfimedTicketsCount: totalUncon,
	}
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(&remainData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}

func main() {
	fmt.Println("hello world")
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/remaining", RemainHandler).Methods("GET")
	r.HandleFunc("/book", BookHandler).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
