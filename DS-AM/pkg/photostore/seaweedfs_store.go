package photostore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"ServiceApi/internal/photoservice/model"
)

// SeaweedFSConfig содержит конфигурацию для подключения к SeaweedFS
type SeaweedFSConfig struct {
	MasterURL       string `json:"master_url"`
	VolumeServerURL string `json:"volume_server_url"`
	FilerURL        string `json:"filer_url"`
	Replication     string `json:"replication"`
}

// SeaweedFSStore реализует хранилище фотографий на базе SeaweedFS
type SeaweedFSStore struct {
	config          SeaweedFSConfig
	client          *http.Client
	externalHost    string // Хост для внешнего доступа к API
	publicVolumeURL string // Публичный URL для доступа к volume server
}

// NewSeaweedFSStore создает новый экземпляр хранилища SeaweedFS
func NewSeaweedFSStore(config SeaweedFSConfig) *SeaweedFSStore {
	// Определяем хост для внешнего доступа из переменной окружения или используем значение по умолчанию
	externalHost := os.Getenv("EXTERNAL_HOST")
	if externalHost == "" {
		// В Docker Compose будем использовать хост машины
		externalHost = "localhost"
	}

	// Определяем публичный URL для доступа к volume server
	publicVolumeURL := os.Getenv("PUBLIC_VOLUME_URL")
	if publicVolumeURL == "" {
		publicVolumeURL = fmt.Sprintf("http://%s:8080", externalHost)
	}

	return &SeaweedFSStore{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		externalHost:    externalHost,
		publicVolumeURL: publicVolumeURL,
	}
}

// SavePhoto сохраняет фотографию в SeaweedFS
func (s *SeaweedFSStore) SavePhoto(ctx context.Context, filename string, data io.Reader, metadata *model.PhotoMetadata) (*model.PhotoMetadata, error) {
	// Получаем URL для загрузки от мастер-сервера
	assignURL := fmt.Sprintf("http://%s/dir/assign?replication=%s", s.config.MasterURL, s.config.Replication)
	resp, err := s.client.Get(assignURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get file ID: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to assign file ID, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var assignResult struct {
		FileID    string `json:"fid"`
		URL       string `json:"url"`
		PublicURL string `json:"publicUrl"`
		Count     int    `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&assignResult); err != nil {
		return nil, fmt.Errorf("failed to decode assign response: %v", err)
	}

	// Загружаем файл на volume server
	// Важно: assignResult.URL содержит имя хоста внутри Docker сети
	uploadURL := fmt.Sprintf("http://%s/%s", assignResult.URL, assignResult.FileID)

	pr, pw := io.Pipe()
	mpw := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer mpw.Close()

		part, err := mpw.CreateFormFile("file", filename)
		if err != nil {
			return
		}

		io.Copy(part, data)
	}()

	req, err := http.NewRequestWithContext(ctx, "POST", uploadURL, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload request: %v", err)
	}

	req.Header.Set("Content-Type", mpw.FormDataContentType())

	uploadResp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(uploadResp.Body)
		return nil, fmt.Errorf("failed to upload file, status: %d, body: %s", uploadResp.StatusCode, string(body))
	}

	// Обновляем метаданные
	metadata.ID = assignResult.FileID
	metadata.StoragePath = assignResult.FileID
	metadata.UploadedAt = time.Now()

	// Формируем публичный URL для доступа к фотографии, заменяя внутренний хост на внешний
	metadata.URLs.Original = fmt.Sprintf("%s/%s", s.publicVolumeURL, assignResult.FileID)

	return metadata, nil
}

// GetPhoto получает фотографию из SeaweedFS
func (s *SeaweedFSStore) GetPhoto(ctx context.Context, id string, size string) (io.ReadCloser, *model.PhotoMetadata, error) {
	fileURL := s.buildFileURL(id, size)

	req, err := http.NewRequestWithContext(ctx, "GET", fileURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("failed to get file, status: %d", resp.StatusCode)
	}

	// Создаем метаданные на основе заголовков
	metadata := &model.PhotoMetadata{
		ID:          id,
		ContentType: resp.Header.Get("Content-Type"),
		StoragePath: id,
		URLs: model.PhotoURLs{
			Original: s.buildPublicFileURL(id, ""),
		},
	}

	if sizeStr := resp.Header.Get("Content-Length"); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
			metadata.FileSize = size
		}
	}

	return resp.Body, metadata, nil
}

// DeletePhoto удаляет фотографию из SeaweedFS
func (s *SeaweedFSStore) DeletePhoto(ctx context.Context, id string) error {
	deleteURL := fmt.Sprintf("http://%s/%s", s.config.VolumeServerURL, id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete file, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SaveProcessedPhoto сохраняет обработанную версию фотографии
func (s *SeaweedFSStore) SaveProcessedPhoto(ctx context.Context, originalID string, size string, data io.Reader) (string, error) {
	// Создаем имя файла с суффиксом размера
	filename := fmt.Sprintf("%s_%s", originalID, size)

	// Сохраняем обработанную версию как новый файл
	metadata := &model.PhotoMetadata{
		FileName: filename,
	}

	savedMetadata, err := s.SavePhoto(ctx, filename, data, metadata)
	if err != nil {
		return "", fmt.Errorf("failed to save processed photo: %v", err)
	}

	return savedMetadata.ID, nil
}

// GetPhotoURL возвращает публичный URL для доступа к фотографии
func (s *SeaweedFSStore) GetPhotoURL(ctx context.Context, id string, size string) (string, error) {
	return s.buildPublicFileURL(id, size), nil
}

// ListPhotos реализация пока не требуется для простого хранилища
func (s *SeaweedFSStore) ListPhotos(ctx context.Context, filter *model.PhotoQueryFilter) ([]*model.PhotoMetadata, error) {
	// В базовой реализации мы не поддерживаем листинг файлов через simple volume server
	// Для полной реализации потребуется использовать SeaweedFS Filer
	return nil, fmt.Errorf("listing photos is not implemented in simple SeaweedFS implementation")
}

// buildFileURL создает внутренний URL для доступа к файлу (для использования внутри сервиса)
func (s *SeaweedFSStore) buildFileURL(id, size string) string {
	baseURL := fmt.Sprintf("http://%s", s.config.VolumeServerURL)

	if size == "" {
		return fmt.Sprintf("%s/%s", baseURL, id)
	}

	// Если запрашивается определенный размер, предполагаем, что файл имеет суффикс с размером
	parts := strings.Split(id, ".")
	ext := ""
	baseID := id

	if len(parts) > 1 {
		ext = "." + parts[len(parts)-1]
		baseID = strings.Join(parts[:len(parts)-1], ".")
	}

	return fmt.Sprintf("%s/%s_%s%s", baseURL, baseID, size, ext)
}

// buildPublicFileURL создает публичный URL для доступа к файлу (для клиентов)
func (s *SeaweedFSStore) buildPublicFileURL(id, size string) string {
	if size == "" {
		return fmt.Sprintf("%s/%s", s.publicVolumeURL, id)
	}

	// Если запрашивается определенный размер, предполагаем, что файл имеет суффикс с размером
	parts := strings.Split(id, ".")
	ext := ""
	baseID := id

	if len(parts) > 1 {
		ext = "." + parts[len(parts)-1]
		baseID = strings.Join(parts[:len(parts)-1], ".")
	}

	return fmt.Sprintf("%s/%s_%s%s", s.publicVolumeURL, baseID, size, ext)
}
