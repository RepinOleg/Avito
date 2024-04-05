package main

import (
	"github.com/RepinOleg/Banner_service/internal/dbs"
	"github.com/RepinOleg/Banner_service/internal/handler"
	sw "github.com/RepinOleg/Banner_service/internal/router"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	cfg := dbs.LoadDBConfig()
	connect, err := dbs.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connect.Close()
	handlers := handler.NewHandler(connect)
	log.Printf("Server started")

	router := sw.NewRouter(handlers)

	log.Fatal(http.ListenAndServe(":8080", router))
}
