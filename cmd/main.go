package main

import (
	"github.com/RepinOleg/Banner_service/internal/dbs"
	"github.com/RepinOleg/Banner_service/internal/handler"
	sw "github.com/RepinOleg/Banner_service/internal/router"
	"log"
	"net/http"
)

func main() {
	// Создаем конфигурацию для подключения к БД
	cfg := dbs.LoadDBConfig()
	//подключаемся
	connect, err := dbs.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer connect.Close()

	handlers := handler.NewHandler(connect)
	// создаем структуру с подключением к БД и будущим кэшом
	//app := handler.DBHandler{DB: db}
	log.Printf("Server started")
	// создаем роутер TODO подумать как сделать связь хэндлеров и БД
	router := sw.NewRouter(handlers)

	log.Fatal(http.ListenAndServe(":8080", router))
}
