package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello world")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func RemainHandler(w http.ResponseWriter, r *http.Request) {
	type ticketRemain struct {
		Remain int `json:"remain"`
	}
	remainData := ticketRemain{Remain: 1}
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(&remainData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}

func BookHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println("hello world")
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/remain", RemainHandler).Methods("GET")
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
