package main

import (
	"github.com/RepinOleg/Banner_service/internal/dbs"
	"github.com/RepinOleg/Banner_service/internal/handler"
	"github.com/RepinOleg/Banner_service/internal/memorycache"
	sw "github.com/RepinOleg/Banner_service/internal/router"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := dbs.LoadDBConfig()

	connect, err := dbs.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connect.Close()

	cache := memorycache.New(5*time.Minute, 10*time.Minute)
	repository := dbs.NewRepository(connect)
	handlers := handler.NewHandler(repository, cache)
	log.Printf("Server started")

	router := sw.NewRouter(handlers)

	log.Fatal(http.ListenAndServe(":8080", router))
}
