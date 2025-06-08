package service

import (
	"ServiceApi/internal/photoservice/repository"
)

// Service представляет собой контейнер для всех сервисов
type Service struct {
	PhotoService *PhotoService
}

// NewService создает новый экземпляр контейнера сервисов
func NewService(repo *repository.Repository, imageProcessorConfig map[string]interface{}) *Service {
	imageProcessor := NewImageProcessor(imageProcessorConfig)
	return &Service{
		PhotoService: NewPhotoService(repo.PhotoRepo, imageProcessor),
	}
}
