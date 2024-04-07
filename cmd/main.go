package main

import (
	"github.com/RepinOleg/Banner_service/internal/handler"
	"github.com/RepinOleg/Banner_service/internal/memorycache"
	"github.com/RepinOleg/Banner_service/internal/repository"
	sw "github.com/RepinOleg/Banner_service/internal/router"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := repository.LoadDBConfig()

	connect, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connect.Close()

	cache := memorycache.New(5*time.Minute, 10*time.Minute)
	repo := repository.NewRepository(connect)
	handlers := handler.NewHandler(repo, cache)
	log.Printf("Server started")

	router := sw.NewRouter(handlers)

	log.Fatal(http.ListenAndServe(":8080", router))
}
