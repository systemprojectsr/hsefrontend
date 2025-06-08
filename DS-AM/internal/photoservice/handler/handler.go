package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"ServiceApi/internal/photoservice/model"
	"ServiceApi/internal/photoservice/service"
)

// Handler представляет собой обработчик HTTP запросов для фотосервиса
type Handler struct {
	service *service.Service
}

// NewHandler создает новый экземпляр обработчика
func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

// SetupRoutes настраивает маршруты для API фотосервиса
func (h *Handler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Маршрут для корневого URL
	mux.HandleFunc("/", h.RootHandler)

	// Маршрут для проверки статуса
	mux.HandleFunc("/status", h.StatusHandler)

	// Маршруты для работы с фотографиями
	mux.HandleFunc("/photos/upload", h.UploadPhoto)
	mux.HandleFunc("/photos", h.ListPhotos)
	mux.HandleFunc("/photos/", h.GetPhoto)
	mux.HandleFunc("/photos/process", h.ProcessPhoto)

	// Middleware для логирования и обработки CORS
	return h.applyMiddleware(mux)
}

// RootHandler обрабатывает запросы к корневому URL
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	// Если запрос именно на корневой URL, отдаем статус и информацию о сервисе
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "ok",
		"service": "Photo Service API",
		"version": "1.0.0",
		"endpoints": []string{
			"/photos/upload - Upload a new photo",
			"/photos/{id} - Get a photo by ID",
			"/photos - List photos with filtering",
			"/photos/process - Process an existing photo",
		},
	}
	json.NewEncoder(w).Encode(response)
}

// UploadPhoto обрабатывает загрузку новой фотографии
func (h *Handler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ограничиваем размер загружаемого файла (10 МБ)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем дополнительные параметры из формы
	uploadRequest := &model.UploadPhotoRequest{
		UserID:       r.FormValue("user_id"),
		CompanyID:    r.FormValue("company_id"),
		TaskID:       r.FormValue("task_id"),
		IsTaskResult: r.FormValue("is_task_result") == "true",
	}

	// Проверяем обязательные параметры
	if uploadRequest.UserID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Загружаем и обрабатываем фотографию
	metadata, err := h.service.PhotoService.UploadPhoto(r.Context(), file, header, uploadRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload photo: %v", err), http.StatusInternalServerError)
		return
	}

	// Возвращаем метаданные загруженной фотографии
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

// GetPhoto обрабатывает запрос на получение фотографии
func (h *Handler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID фотографии из URL
	id := r.URL.Path[len("/photos/"):]
	if id == "" {
		http.Error(w, "Photo ID is required", http.StatusBadRequest)
		return
	}

	// Получаем размер фотографии из параметров запроса
	size := r.URL.Query().Get("size")
	// Если размер не указан, возвращаем оригинал
	if size == "" {
		size = "original"
	}

	// Проверяем, нужны ли только метаданные
	metadataOnly := r.URL.Query().Get("metadata") == "true"
	if metadataOnly {
		// Получаем фотографию только для метаданных
		_, metadata, err := h.service.PhotoService.GetPhoto(r.Context(), id, "")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get photo metadata: %v", err), http.StatusInternalServerError)
			return
		}

		// Возвращаем метаданные
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metadata)
		return
	}

	// Получаем фотографию из хранилища
	photoReader, metadata, err := h.service.PhotoService.GetPhoto(r.Context(), id, size)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get photo: %v", err), http.StatusInternalServerError)
		return
	}
	defer photoReader.Close()

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", metadata.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", metadata.FileName))
	w.Header().Set("Content-Length", strconv.FormatInt(metadata.FileSize, 10))

	// Отправляем содержимое файла клиенту
	io.Copy(w, photoReader)
}

// ListPhotos обрабатывает запрос на получение списка фотографий
func (h *Handler) ListPhotos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Создаем фильтр на основе параметров запроса
	filter := &model.PhotoQueryFilter{
		UserID:    r.URL.Query().Get("user_id"),
		CompanyID: r.URL.Query().Get("company_id"),
		TaskID:    r.URL.Query().Get("task_id"),
		FromDate:  r.URL.Query().Get("from_date"),
		ToDate:    r.URL.Query().Get("to_date"),
	}

	// Парсим лимит и смещение
	if limit := r.URL.Query().Get("limit"); limit != "" {
		limitVal, err := strconv.Atoi(limit)
		if err == nil && limitVal > 0 {
			filter.Limit = limitVal
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		offsetVal, err := strconv.Atoi(offset)
		if err == nil && offsetVal >= 0 {
			filter.Offset = offsetVal
		}
	}

	// Парсим флаг IsTaskResult
	if isTaskResult := r.URL.Query().Get("is_task_result"); isTaskResult != "" {
		boolVal := isTaskResult == "true"
		filter.IsTaskResult = &boolVal
	}

	// Получаем список фотографий
	photos, err := h.service.PhotoService.ListPhotos(r.Context(), filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list photos: %v", err), http.StatusInternalServerError)
		return
	}

	// Возвращаем список фотографий
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(photos)
}

// ProcessPhoto обрабатывает запрос на обработку фотографии
func (h *Handler) ProcessPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID фотографии из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Photo ID is required", http.StatusBadRequest)
		return
	}

	// Декодируем опции обработки из тела запроса
	var options model.ProcessingOptions
	if err := json.NewDecoder(r.Body).Decode(&options); err != nil {
		http.Error(w, "Failed to parse processing options", http.StatusBadRequest)
		return
	}

	// Обрабатываем фотографию
	processedImage, err := h.service.PhotoService.ProcessPhoto(r.Context(), id, &options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process photo: %v", err), http.StatusInternalServerError)
		return
	}

	// Определяем тип содержимого на основе опций
	contentType := "image/jpeg" // По умолчанию
	// TODO: Определить тип содержимого на основе формата выходного файла

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "inline")

	// Отправляем обработанное изображение клиенту
	io.Copy(w, processedImage)
}

// StatusHandler обрабатывает запросы на проверку статуса сервисов
func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем статус хранилища
	storageStatus := "online"
	storageError := ""

	// Пробуем выполнить простой запрос к хранилищу
	_, err := h.service.PhotoService.ListPhotos(r.Context(), &model.PhotoQueryFilter{Limit: 1})
	if err != nil {
		storageStatus = "offline"
		storageError = err.Error()
	}

	// Формируем ответ
	status := map[string]interface{}{
		"service":   "Photo Service API",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"storage": map[string]interface{}{
			"status": storageStatus,
			"error":  storageError,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// applyMiddleware применяет промежуточное ПО к обработчикам
func (h *Handler) applyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логирование запроса
		fmt.Printf("[%s] %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)

		// Заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 часа

		// Обработка предварительных запросов CORS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверка, поддерживается ли метод для данного пути
		if r.URL.Path == "/" {
			// Для корневого пути поддерживаем только GET
			if r.Method != http.MethodGet && r.Method != http.MethodOptions {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}

		// Вызов следующего обработчика
		next.ServeHTTP(w, r)
	})
}
