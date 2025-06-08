package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"

	"ServiceApi/internal/photoservice/model"
)

// ImageProcessor представляет собой сервис для обработки изображений
type ImageProcessor struct {
	thumbnailWidth  int
	thumbnailHeight int
	mediumWidth     int
	mediumHeight    int
	largeWidth      int
	largeHeight     int
	quality         map[string]int
}

// NewImageProcessor создает новый экземпляр сервиса для обработки изображений
func NewImageProcessor(config map[string]interface{}) *ImageProcessor {
	thumbnailConfig := config["thumbnail"].(map[string]interface{})
	mediumConfig := config["medium"].(map[string]interface{})
	largeConfig := config["large"].(map[string]interface{})

	quality := make(map[string]int)
	quality["thumbnail"] = int(thumbnailConfig["quality"].(float64))
	quality["medium"] = int(mediumConfig["quality"].(float64))
	quality["large"] = int(largeConfig["quality"].(float64))

	return &ImageProcessor{
		thumbnailWidth:  int(thumbnailConfig["width"].(float64)),
		thumbnailHeight: int(thumbnailConfig["height"].(float64)),
		mediumWidth:     int(mediumConfig["width"].(float64)),
		mediumHeight:    int(mediumConfig["height"].(float64)),
		largeWidth:      int(largeConfig["width"].(float64)),
		largeHeight:     int(largeConfig["height"].(float64)),
		quality:         quality,
	}
}

// ProcessImage обрабатывает изображение и создает различные его размеры
func (p *ImageProcessor) ProcessImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (map[string]io.Reader, error) {
	// Читаем содержимое файла
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %v", err)
	}

	// Определяем формат файла
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentTypeFromExtension(fileHeader.Filename)
	}

	// Декодируем изображение
	img, format, err := image.Decode(bytes.NewReader(fileContent))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Создаем разные размеры изображения
	results := make(map[string]io.Reader)

	// Оригинал
	results["original"] = bytes.NewReader(fileContent)

	// Thumbnail
	thumbnailImg := imaging.Resize(img, p.thumbnailWidth, p.thumbnailHeight, imaging.Lanczos)
	thumbnailBuf := new(bytes.Buffer)
	err = p.encodeImage(thumbnailImg, thumbnailBuf, format, p.quality["thumbnail"])
	if err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %v", err)
	}
	results["thumbnail"] = bytes.NewReader(thumbnailBuf.Bytes())

	// Medium
	mediumImg := imaging.Resize(img, p.mediumWidth, p.mediumHeight, imaging.Lanczos)
	mediumBuf := new(bytes.Buffer)
	err = p.encodeImage(mediumImg, mediumBuf, format, p.quality["medium"])
	if err != nil {
		return nil, fmt.Errorf("failed to encode medium image: %v", err)
	}
	results["medium"] = bytes.NewReader(mediumBuf.Bytes())

	// Large
	largeImg := imaging.Resize(img, p.largeWidth, p.largeHeight, imaging.Lanczos)
	largeBuf := new(bytes.Buffer)
	err = p.encodeImage(largeImg, largeBuf, format, p.quality["large"])
	if err != nil {
		return nil, fmt.Errorf("failed to encode large image: %v", err)
	}
	results["large"] = bytes.NewReader(largeBuf.Bytes())

	return results, nil
}

// ProcessImageWithOptions обрабатывает изображение с заданными опциями
func (p *ImageProcessor) ProcessImageWithOptions(ctx context.Context, originalImage io.Reader, options *model.ProcessingOptions) (io.Reader, error) {
	// Читаем содержимое файла
	fileContent, err := ioutil.ReadAll(originalImage)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %v", err)
	}

	// Декодируем изображение
	img, format, err := image.Decode(bytes.NewReader(fileContent))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Применяем операции обработки
	processed := img

	// Обрезка изображения
	if options.Crop {
		processed = imaging.Crop(processed, image.Rect(options.CropX, options.CropY, options.CropX+options.CropWidth, options.CropY+options.CropHeight))
	}

	// Изменение размера
	if options.Resize {
		processed = imaging.Resize(processed, options.Width, options.Height, imaging.Lanczos)
	}

	// Кодируем результат
	buf := new(bytes.Buffer)
	quality := 85 // По умолчанию
	if options.Quality > 0 {
		quality = options.Quality
	}
	err = p.encodeImage(processed, buf, format, quality)
	if err != nil {
		return nil, fmt.Errorf("failed to encode processed image: %v", err)
	}

	return bytes.NewReader(buf.Bytes()), nil
}

// encodeImage кодирует изображение в нужный формат
func (p *ImageProcessor) encodeImage(img image.Image, w io.Writer, format string, quality int) error {
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		return encoder.Encode(w, img)
	default:
		// По умолчанию используем JPEG
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	}
}

// getContentTypeFromExtension определяет тип содержимого по расширению файла
func getContentTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
