package main

import (
	"log"
	"net/http"

	sw "github.com/RepinOleg/Banner_service/go"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
