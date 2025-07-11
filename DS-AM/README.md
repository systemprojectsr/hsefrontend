# Микросервис для хранения и обработки фотографий

Данный микросервис предназначен для хранения и обработки фотографий с использованием SeaweedFS в качестве S3-совместимого хранилища.

## Особенности

- Загрузка фотографий и автоматическое создание нескольких размеров (thumbnail, medium, large)
- Получение фотографий разных размеров по ID
- Обработка изображений: изменение размера, обрезка
- Хранение метаданных фотографий
- Фильтрация фотографий по пользователю, компании, задаче, и т.д.
- Докеризация всех компонентов для простоты запуска и развертывания

## Требования

- Docker и Docker Compose
- 1 ГБ свободной оперативной памяти (минимум)
- 2 ГБ свободного дискового пространства (минимум)

## Быстрый старт

### Запуск сервиса

```bash
# Клонировать репозиторий
git clone <repository-url>
cd ServiceApi

# Запустить все сервисы в фоновом режиме
make up

# Или вручную через docker-compose
docker-compose up -d
```

После запуска сервис будет доступен по адресу: http://localhost:8081

### Проверка статуса

```bash
# Проверить статус всех контейнеров
make ps

# Просмотр логов
make logs
```

### Остановка сервиса

```bash
# Остановить все контейнеры
make down

# Удалить все данные и контейнеры
make clean
```

## Структура проекта в Docker

- **photoservice**: Основной микросервис для обработки и хранения фотографий
- **seaweedfs-master**: Мастер-сервер SeaweedFS, управляет распределением и хранением данных
- **seaweedfs-volume**: Volume-сервер SeaweedFS, хранит физические данные (фотографии)
- **seaweedfs-filer**: Filer-сервер SeaweedFS, предоставляет файловый интерфейс и метаданные

## API

### Загрузка фотографии
```
POST /photos/upload
Content-Type: multipart/form-data

Параметры:
- file: Файл изображения
- user_id: ID пользователя (обязательно)
- company_id: ID компании (опционально)
- task_id: ID задачи (опционально)
- is_task_result: Является ли результатом задачи (true/false)
```

### Получение фотографии
```
GET /photos/{id}?size=[original|thumbnail|medium|large]

Параметры:
- size: Размер изображения (по умолчанию: original)
- metadata: Если true, возвращает только метаданные без самого изображения
```

### Список фотографий
```
GET /photos?user_id={user_id}&company_id={company_id}&task_id={task_id}&is_task_result={true|false}&from_date={date}&to_date={date}&limit={limit}&offset={offset}

Все параметры опциональны и используются для фильтрации.
```

### Обработка фотографии
```
POST /photos/process?id={photo_id}
Content-Type: application/json

Пример тела запроса:
{
  "resize": true,
  "width": 800,
  "height": 600,
  "quality": 85,
  "crop": false
}
```

## Переменные окружения

Фотосервис использует следующие переменные окружения:

- `EXTERNAL_HOST`: Внешний хост для доступа к API (по умолчанию: localhost)
- `PUBLIC_VOLUME_URL`: Публичный URL для доступа к volume server (по умолчанию: http://{EXTERNAL_HOST}:8080)

## Управление через Makefile

В проекте есть Makefile со следующими командами:

- `make up`: Запуск всех сервисов
- `make down`: Остановка всех сервисов
- `make clean`: Удаление всех контейнеров и данных
- `make logs`: Просмотр логов всех сервисов
- `make logs-photo`: Просмотр логов только фотосервиса
- `make logs-seaweed`: Просмотр логов только SeaweedFS
- `make ps`: Проверка статуса контейнеров
- `make restart-photo`: Перезапуск фотосервиса
- `make rebuild`: Полное пересоздание и запуск системы
- `make disk-usage`: Проверка использования дискового пространства 