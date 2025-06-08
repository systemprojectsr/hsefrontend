package main

import (
	"ServiceApi/internal/handler"
	"ServiceApi/internal/repository"
	"ServiceApi/internal/service"
	"ServiceApi/pkg/database"
	"log"
	"net/http"
)

func main() {
	// Инициализация базы данных
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация репозитория, сервиса и хендлера
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	// Настройка маршрутов
	http.HandleFunc("/search", h.SearchServices)

	// Запуск сервера
	log.Println("Starting server on :8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
