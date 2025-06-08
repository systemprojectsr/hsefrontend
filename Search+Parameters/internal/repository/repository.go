package repository

import (
	"ServiceApi/internal/model"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"strings"

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

		return b.ForEach(func(k, v []byte) error {
			var service model.Service
			if err := json.Unmarshal(v, &service); err != nil {
				return err
			}

			if applyFilters(service, filters) {
				services = append(services, service)
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	// Сортировка
	if sortParam, ok := filters["sort"]; ok {
		switch sortParam {
		case "rating_low":
			sort.Slice(services, func(i, j int) bool {
				return services[i].Rating < services[j].Rating
			})
		case "rating_high":
			sort.Slice(services, func(i, j int) bool {
				return services[i].Rating > services[j].Rating
			})
		case "price_low":
			sort.Slice(services, func(i, j int) bool {
				return services[i].Price < services[j].Price
			})
		case "price_high":
			sort.Slice(services, func(i, j int) bool {
				return services[i].Price > services[j].Price
			})
		default:
			log.Printf("Неизвестный параметр сортировки: %s", sortParam)
		}
	}

	return services, nil
}

func applyFilters(service model.Service, filters map[string]string) bool {
	// Фильтр по цене (от и до)
	if priceFrom, ok := filters["price_from"]; ok && priceFrom != "" {
		priceFromFilter, err := strconv.ParseFloat(priceFrom, 64)
		if err != nil {
			log.Printf("Ошибка парсинга price_from: %v\n", err)
			return false
		}
		if service.Price < priceFromFilter {
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
			return false
		}
	}

	// Остальные фильтры без изменений
	if location, ok := filters["location"]; ok && location != "" {
		if !strings.EqualFold(service.Location, location) {
			return false
		}
	}

	if rating, ok := filters["rating"]; ok && rating != "" {
		ratingFilter, err := strconv.ParseFloat(rating, 64)
		if err != nil {
			log.Printf("Ошибка парсинга rating: %v\n", err)
			return false
		}
		if service.Rating < ratingFilter {
			return false
		}
	}

	return true
}
