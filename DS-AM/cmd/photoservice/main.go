package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ServiceApi/internal/photoservice/handler"
	"ServiceApi/internal/photoservice/repository"
	"ServiceApi/internal/photoservice/service"
	"ServiceApi/pkg/photostore"
)

// Config содержит конфигурацию для фотосервиса
type Config struct {
	Server struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	} `json:"server"`
	Storage struct {
		Type   string                 `json:"type"`
		Config map[string]interface{} `json:"config"`
	} `json:"storage"`
	ImageProcessing map[string]interface{} `json:"image_processing"`
}

func main() {
	// Загружаем конфигурацию
	config, err := loadConfig("configs/photoservice.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализируем хранилище фотографий
	var photoRepo repository.PhotoRepository
	if config.Storage.Type == "seaweedfs" {
		// Получаем конфигурацию для SeaweedFS
		seaweedConfig := photostore.SeaweedFSConfig{
			MasterURL:       config.Storage.Config["master_url"].(string),
			VolumeServerURL: config.Storage.Config["volume_server_url"].(string),
			FilerURL:        config.Storage.Config["filer_url"].(string),
			Replication:     config.Storage.Config["replication"].(string),
		}
		photoRepo = photostore.NewSeaweedFSStore(seaweedConfig)
	} else {
		log.Fatalf("Unsupported storage type: %s", config.Storage.Type)
	}

	// Инициализируем репозиторий, сервис и обработчик
	repo := repository.NewRepository(photoRepo)
	svc := service.NewService(repo, config.ImageProcessing)
	h := handler.NewHandler(svc)

	// Настраиваем маршруты
	router := h.SetupRoutes()

	// Запускаем сервер
	serverAddr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Starting photo service on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Настраиваем корректное завершение при получении сигнала
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// loadConfig загружает конфигурацию из JSON файла
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %v", err)
	}

	return &config, nil
}
