package main

import (
	"github.com/RepinOleg/Banner_service/internal/dbs"
	sw "github.com/RepinOleg/Banner_service/internal/router"
	"log"
	"net/http"
)

func main() {
	// Создаем конфигурацию для подключения к БД
	cfg := dbs.LoadDBConfig()
	//подключаемся
	db, err := dbs.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// создаем структуру с подключением к БД и будущим кэшом
	//app := handlers.DBHandler{DB: db}
	log.Printf("Server started")
	// создаем роутер TODO подумать как сделать связь хэндлеров и БД
	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
