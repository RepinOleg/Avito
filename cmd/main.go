package main

import (
	"log"
	"net/http"

	"github.com/RepinOleg/Banner_service/internal/handler"
	"github.com/RepinOleg/Banner_service/internal/repository"
	"github.com/RepinOleg/Banner_service/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	cfg := repository.LoadDBConfig()

	connect, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer connect.Close()

	repos := repository.NewRepository(connect)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	log.Printf("Server started")

	router := handler.NewRouter(handlers)

	if err = http.ListenAndServe(viper.GetString("port"), router); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
