package repository

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"ServiceApi/internal/model"
	"go.etcd.io/bbolt"
)

type Repository struct {
	db *bbolt.DB
}

func NewRepository(db *bbolt.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindServices(filters map[string]string) ([]model.Service, error) {
	var services []model.Service

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("services"))
		if b == nil {
			log.Println("Bucket 'services' не найден")
			return nil
		}

		// Итерация по всем записям в Bucket
		return b.ForEach(func(k, v []byte) error {
			var service model.Service
			if err := json.Unmarshal(v, &service); err != nil {
				return err
			}

			// Применение фильтров
			if applyFilters(service, filters) {
				services = append(services, service)
			}

			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Найдено услуг: %d\n", len(services))
	return services, nil
}

// applyFilters проверяет, соответствует ли услуга фильтрам
func applyFilters(service model.Service, filters map[string]string) bool {
	// Фильтр по цене (от и до)
	if priceFrom, ok := filters["price_from"]; ok && priceFrom != "" {
		priceFromFilter, err := strconv.ParseFloat(priceFrom, 64)
		if err != nil {
			log.Printf("Ошибка парсинга price_from: %v\n", err)
			return false
		}
		if service.Price < priceFromFilter {
			log.Printf("Услуга %s не прошла фильтр price_from\n", service.Name)
			return false
		}
	}

	if priceTo, ok := filters["price_to"]; ok && priceTo != "" {
		priceToFilter, err := strconv.ParseFloat(priceTo, 64)
		if err != nil {
			log.Printf("Ошибка парсинга price_to: %v\n", err)
			return false
		}
		if service.Price > priceToFilter {
			log.Printf("Услуга %s не прошла фильтр price_to\n", service.Name)
			return false
		}
	}

	// Фильтр по району (нечувствительный к регистру)
	if location, ok := filters["location"]; ok && location != "" {
		if !strings.EqualFold(service.Location, location) {
			log.Printf("Услуга %s не прошла фильтр location\n", service.Name)
			return false
		}
	}

	// Фильтр по рейтингу
	if rating, ok := filters["rating"]; ok && rating != "" {
		ratingFilter, err := strconv.ParseFloat(rating, 64)
		if err != nil {
			log.Printf("Ошибка парсинга rating: %v\n", err)
			return false
		}
		if service.Rating < ratingFilter {
			log.Printf("Услуга %s не прошла фильтр rating\n", service.Name)
			return false
		}
	}

	log.Printf("Услуга %s прошла все фильтры\n", service.Name)
	return true
}