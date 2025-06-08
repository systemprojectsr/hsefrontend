package model

import (
	"time"
)

// PhotoMetadata содержит метаданные о сохраненной фотографии
type PhotoMetadata struct {
	ID           string    `json:"id"`
	FileName     string    `json:"file_name"`
	FileSize     int64     `json:"file_size"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	UserID       string    `json:"user_id"`
	CompanyID    string    `json:"company_id,omitempty"`
	TaskID       string    `json:"task_id,omitempty"`
	IsTaskResult bool      `json:"is_task_result"`
	StoragePath  string    `json:"storage_path"`
	URLs         PhotoURLs `json:"urls"`
}

// PhotoURLs содержит URL-адреса для доступа к различным размерам фотографии
type PhotoURLs struct {
	Original  string `json:"original"`
	Thumbnail string `json:"thumbnail"`
	Medium    string `json:"medium"`
	Large     string `json:"large"`
}

// UploadPhotoRequest представляет запрос на загрузку фотографии
type UploadPhotoRequest struct {
	UserID       string `json:"user_id"`
	CompanyID    string `json:"company_id,omitempty"`
	TaskID       string `json:"task_id,omitempty"`
	IsTaskResult bool   `json:"is_task_result"`
}

// ProcessingOptions определяет параметры обработки изображения
type ProcessingOptions struct {
	Resize     bool `json:"resize"`
	Width      int  `json:"width,omitempty"`
	Height     int  `json:"height,omitempty"`
	Quality    int  `json:"quality,omitempty"`
	Watermark  bool `json:"watermark,omitempty"`
	Crop       bool `json:"crop,omitempty"`
	CropX      int  `json:"crop_x,omitempty"`
	CropY      int  `json:"crop_y,omitempty"`
	CropWidth  int  `json:"crop_width,omitempty"`
	CropHeight int  `json:"crop_height,omitempty"`
}

// PhotoQueryFilter содержит параметры для фильтрации при поиске фотографий
type PhotoQueryFilter struct {
	UserID       string `json:"user_id,omitempty"`
	CompanyID    string `json:"company_id,omitempty"`
	TaskID       string `json:"task_id,omitempty"`
	IsTaskResult *bool  `json:"is_task_result,omitempty"`
	FromDate     string `json:"from_date,omitempty"`
	ToDate       string `json:"to_date,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Offset       int    `json:"offset,omitempty"`
}
