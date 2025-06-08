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
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	http.HandleFunc("/search", h.SearchServices)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
