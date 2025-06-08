#!/bin/bash

# Скрипт для тестирования загрузки фотографий в фотосервис

# Проверяем, что предоставлен путь к файлу
if [ -z "$1" ]; then
  echo "Использование: $0 <путь_к_изображению>"
  echo "Пример: $0 ./test.jpg"
  exit 1
fi

IMAGE_PATH="$1"
API_URL="http://localhost:8081"

# Проверяем, что файл существует
if [ ! -f "$IMAGE_PATH" ]; then
  echo "Ошибка: Файл $IMAGE_PATH не найден"
  exit 1
fi

# Проверяем, что сервис доступен
echo "Проверка доступности сервиса..."
if ! curl -s --head --fail "$API_URL" > /dev/null; then
  echo "Ошибка: Сервис недоступен по адресу $API_URL"
  echo "Убедитесь, что сервис запущен: make up"
  exit 1
fi

echo "Загрузка файла $IMAGE_PATH на сервер..."

# Загрузка фотографии с тестовыми параметрами
RESPONSE=$(curl -s -X POST \
  -F "file=@$IMAGE_PATH" \
  -F "user_id=test_user" \
  -F "company_id=test_company" \
  -F "is_task_result=false" \
  "$API_URL/photos/upload")

# Проверяем, содержит ли ответ ID фотографии
if echo "$RESPONSE" | grep -q "id"; then
  echo "Загрузка успешна! Ответ сервера:"
  echo "$RESPONSE" | python -m json.tool 2>/dev/null || echo "$RESPONSE"
  
  # Извлекаем ID для последующих запросов
  PHOTO_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*' | sed 's/"id":"//g')
  
  echo "ID фотографии: $PHOTO_ID"
  echo
  echo "Для просмотра оригинала: $API_URL/photos/$PHOTO_ID?size=original"
  echo "Для просмотра превью: $API_URL/photos/$PHOTO_ID?size=thumbnail"
  echo "Для просмотра среднего размера: $API_URL/photos/$PHOTO_ID?size=medium"
  echo "Для просмотра большого размера: $API_URL/photos/$PHOTO_ID?size=large"
else
  echo "Ошибка при загрузке. Ответ сервера:"
  echo "$RESPONSE"
  exit 1
fi 