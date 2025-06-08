.PHONY: build up down logs ps clean

# Запустить сервисы в фоновом режиме
up:
	docker-compose up -d

# Остановить и удалить контейнеры
down:
	docker-compose down

# Остановить и удалить контейнеры вместе с томами (данные будут потеряны)
clean:
	docker-compose down -v

# Сборка только фотосервиса (без запуска)
build:
	docker-compose build photoservice

# Вывод логов всех сервисов
logs:
	docker-compose logs -f

# Вывод логов только фотосервиса
logs-photo:
	docker-compose logs -f photoservice

# Вывод логов только SeaweedFS
logs-seaweed:
	docker-compose logs -f seaweedfs-master seaweedfs-volume seaweedfs-filer

# Проверка статуса запущенных контейнеров
ps:
	docker-compose ps

# Перезапуск фотосервиса
restart-photo:
	docker-compose restart photoservice

# Полный пересоздание и запуск системы (с перекомпиляцией)
rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Проверка дискового пространства, используемого контейнерами и образами
disk-usage:
	docker system df 