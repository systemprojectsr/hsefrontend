package database

import (
    "encoding/json"
    "log"
    "ServiceApi/internal/model"
    "go.etcd.io/bbolt"
)

func Connect() (*bbolt.DB, error) {
	// Подключение к BoltDB (файл базы данных будет создан автоматически)
	db, err := bbolt.Open("catalog.db", 0600, nil)
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
		return nil, err
	}

	// Инициализация данных
	err = initData(db)
	if err != nil {
		log.Fatalf("Could not initialize data: %v", err)
		return nil, err
	}

	return db, nil
}

func initData(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		// Создаем Bucket для услуг
		b, err := tx.CreateBucketIfNotExists([]byte("services"))
		if err != nil {
			return err
		}

		// Данные для инициализации
		services := []model.Service{
			{ID: 1, Name: "Профессиональная мойка окон", Description: "Чистые окна без разводов", Price: 2000.00, Location: "Центральный район", Rating: 4.5},
			{ID: 2, Name: "Химчистка мебели", Description: "Глубокая очистка вашей мебели", Price: 3000.00, Location: "Северный район", Rating: 4.7},
			{ID: 3, Name: "Компьютерный мастер", Description: "Ремонт и настройка компьютеров", Price: 1500.00, Location: "Западный район", Rating: 4.3},
			{ID: 4, Name: "Ремонт стиральных машин", Description: "Быстрый и качественный ремонт", Price: 2500.00, Location: "Восточный район", Rating: 4.6},
			{ID: 5, Name: "Электрика", Description: "Услуги электрика на дому", Price: 1800.00, Location: "Южный район", Rating: 4.4},
			{ID: 6, Name: "Слесарные работы", Description: "Ремонт и установка сантехники", Price: 2200.00, Location: "Центральный район", Rating: 4.2},
			{ID: 7, Name: "Сантехника", Description: "Устранение засоров и протечек", Price: 2000.00, Location: "Северный район", Rating: 4.5},
			{ID: 8, Name: "Вывоз мусора", Description: "Быстрый вывоз строительного мусора", Price: 1000.00, Location: "Западный район", Rating: 4.1},
		}

		// Добавляем данные в Bucket
		for _, service := range services {
			data, err := json.Marshal(service)
			if err != nil {
				return err
			}
			err = b.Put([]byte(string(service.ID)), data)
			if err != nil {
				return err
			}
			log.Printf("Добавлена услуга: %+v\n", service)
		}

		return nil
	})
}