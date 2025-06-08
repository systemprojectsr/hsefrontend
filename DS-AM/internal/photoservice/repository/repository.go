package repository

import (
	"ServiceApi/internal/photoservice/model"
	"context"
	"io"
)

// PhotoRepository определяет интерфейс для работы с хранилищем фотографий
type PhotoRepository interface {
	// SavePhoto сохраняет фотографию в хранилище и возвращает её метаданные
	SavePhoto(ctx context.Context, filename string, data io.Reader, metadata *model.PhotoMetadata) (*model.PhotoMetadata, error)

	// GetPhoto получает фотографию из хранилища по ID
	GetPhoto(ctx context.Context, id string, size string) (io.ReadCloser, *model.PhotoMetadata, error)

	// DeletePhoto удаляет фотографию из хранилища по ID
	DeletePhoto(ctx context.Context, id string) error

	// ListPhotos получает список фотографий с учетом фильтра
	ListPhotos(ctx context.Context, filter *model.PhotoQueryFilter) ([]*model.PhotoMetadata, error)

	// SaveProcessedPhoto сохраняет обработанную версию фотографии
	SaveProcessedPhoto(ctx context.Context, originalID string, size string, data io.Reader) (string, error)

	// GetPhotoURL возвращает URL для доступа к фотографии
	GetPhotoURL(ctx context.Context, id string, size string) (string, error)
}

// Repository содержит все репозитории
type Repository struct {
	PhotoRepo PhotoRepository
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(photoRepo PhotoRepository) *Repository {
	return &Repository{
		PhotoRepo: photoRepo,
	}
}
