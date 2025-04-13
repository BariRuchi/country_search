package main

import (
	"CountrySearch/countryhandler"
	"log"
	"net/http"
)

func main() {
	handler := countryhandler.New()
	http.HandleFunc("/api/countries/search", handler.CountryHandler)

	log.Println("Listening on port :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
