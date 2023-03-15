package main

import (
	"github.com/dobb2/go-musthave-devops-tpl/internal/handlers"
	"github.com/dobb2/go-musthave-devops-tpl/internal/storage/metrics/cache"
	"log"
	"net/http"
)

func main() {
	datastore := cache.Create()
	handler := handlers.New(datastore)
	// маршрутизация запросов обработчику
	http.HandleFunc("/update/", handler.Update)
	http.HandleFunc("/", handler.Other)
	// запуск сервера с адресом localhost, порт 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
