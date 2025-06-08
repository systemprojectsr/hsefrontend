package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"ServiceApi/internal/photoservice/model"
	"ServiceApi/internal/photoservice/repository"
)

// PhotoService представляет собой сервис для работы с фотографиями
type PhotoService struct {
	photoRepo      repository.PhotoRepository
	imageProcessor *ImageProcessor
}

// NewPhotoService создает новый экземпляр сервиса для работы с фотографиями
func NewPhotoService(photoRepo repository.PhotoRepository, imageProcessor *ImageProcessor) *PhotoService {
	return &PhotoService{
		photoRepo:      photoRepo,
		imageProcessor: imageProcessor,
	}
}

// UploadPhoto загружает фотографию, обрабатывает ее и сохраняет в хранилище
func (s *PhotoService) UploadPhoto(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, request *model.UploadPhotoRequest) (*model.PhotoMetadata, error) {
	// Создаем метаданные фотографии
	metadata := &model.PhotoMetadata{
		FileName:     fileHeader.Filename,
		FileSize:     fileHeader.Size,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		UploadedAt:   time.Now(),
		UserID:       request.UserID,
		CompanyID:    request.CompanyID,
		TaskID:       request.TaskID,
		IsTaskResult: request.IsTaskResult,
	}

	// Обрабатываем изображение и создаем разные размеры
	processedImages, err := s.imageProcessor.ProcessImage(ctx, file, fileHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %v", err)
	}

	// Сохраняем оригинальную фотографию
	originalImage := processedImages["original"]
	savedMetadata, err := s.photoRepo.SavePhoto(ctx, fileHeader.Filename, originalImage, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to save original photo: %v", err)
	}

	// Сохраняем обработанные версии
	thumbnailID, err := s.photoRepo.SaveProcessedPhoto(ctx, savedMetadata.ID, "thumbnail", processedImages["thumbnail"])
	if err != nil {
		return nil, fmt.Errorf("failed to save thumbnail: %v", err)
	}

	mediumID, err := s.photoRepo.SaveProcessedPhoto(ctx, savedMetadata.ID, "medium", processedImages["medium"])
	if err != nil {
		return nil, fmt.Errorf("failed to save medium photo: %v", err)
	}

	largeID, err := s.photoRepo.SaveProcessedPhoto(ctx, savedMetadata.ID, "large", processedImages["large"])
	if err != nil {
		return nil, fmt.Errorf("failed to save large photo: %v", err)
	}

	// Формируем URL для доступа к разным размерам
	thumbnailURL, _ := s.photoRepo.GetPhotoURL(ctx, thumbnailID, "thumbnail")
	mediumURL, _ := s.photoRepo.GetPhotoURL(ctx, mediumID, "medium")
	largeURL, _ := s.photoRepo.GetPhotoURL(ctx, largeID, "large")

	// Обновляем метаданные с URL-ами
	savedMetadata.URLs.Thumbnail = thumbnailURL
	savedMetadata.URLs.Medium = mediumURL
	savedMetadata.URLs.Large = largeURL

	return savedMetadata, nil
}

// GetPhoto получает фотографию из хранилища по ID
func (s *PhotoService) GetPhoto(ctx context.Context, id string, size string) (io.ReadCloser, *model.PhotoMetadata, error) {
	return s.photoRepo.GetPhoto(ctx, id, size)
}

// DeletePhoto удаляет фотографию из хранилища по ID
func (s *PhotoService) DeletePhoto(ctx context.Context, id string) error {
	return s.photoRepo.DeletePhoto(ctx, id)
}

// ProcessPhoto обрабатывает фотографию с заданными опциями
func (s *PhotoService) ProcessPhoto(ctx context.Context, id string, options *model.ProcessingOptions) (io.Reader, error) {
	// Получаем оригинальную фотографию
	originalPhoto, _, err := s.photoRepo.GetPhoto(ctx, id, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get original photo: %v", err)
	}
	defer originalPhoto.Close()

	// Обрабатываем изображение с заданными опциями
	processedImage, err := s.imageProcessor.ProcessImageWithOptions(ctx, originalPhoto, options)
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %v", err)
	}

	return processedImage, nil
}

// ListPhotos получает список фотографий с учетом фильтра
func (s *PhotoService) ListPhotos(ctx context.Context, filter *model.PhotoQueryFilter) ([]*model.PhotoMetadata, error) {
	return s.photoRepo.ListPhotos(ctx, filter)
}
