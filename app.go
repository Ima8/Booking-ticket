package main

import (
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

func RemainHeandler(w http.ResponseWriter, r *http.Request) {

}
func main() {
	fmt.Println("hello world")
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/remain", RemainHeandler).Methods("GET")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
