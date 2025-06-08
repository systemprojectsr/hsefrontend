#!/bin/bash

# Скрипт для тестирования обработки фотографий в фотосервисе

# Проверка наличия обязательного параметра (ID фотографии)
if [ -z "$1" ]; then
    echo "Использование: $0 <id_фотографии> [output_file]"
    echo "Пример: $0 3,01ebc098bb processed.jpg"
    exit 1
fi

# ID фотографии из аргумента
PHOTO_ID=$1

# Имя выходного файла (по умолчанию processed_image.jpg)
OUTPUT_FILE=${2:-"processed_image.jpg"}

# Параметры обработки
WIDTH=400
HEIGHT=300
QUALITY=85

# URL сервиса
SERVICE_URL="http://localhost:8081"

# Проверка статуса сервиса перед выполнением запроса
echo "Проверка статуса сервиса..."
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" ${SERVICE_URL}/status)

if [ "$STATUS_CODE" != "200" ]; then
    echo "Ошибка: Сервис недоступен (статус код: $STATUS_CODE)"
    echo "Проверьте, запущен ли сервис и доступен ли он по адресу ${SERVICE_URL}"
    exit 1
fi

echo "Сервис доступен. Выполняется запрос на обработку фотографии с ID: $PHOTO_ID"
echo "Параметры обработки: ширина=$WIDTH, высота=$HEIGHT, качество=$QUALITY"
echo "Файл будет сохранен как: $OUTPUT_FILE"

# Выполнение запроса к API для обработки фотографии
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{\"resize\": true, \"width\": $WIDTH, \"height\": $HEIGHT, \"quality\": $QUALITY}" \
  "${SERVICE_URL}/photos/process?id=${PHOTO_ID}" > "$OUTPUT_FILE"

# Проверка результата
if [ -s "$OUTPUT_FILE" ]; then
    echo "Успешно! Обработанное изображение сохранено в файл: $OUTPUT_FILE"
    
    # Определение размера файла
    FILE_SIZE=$(du -h "$OUTPUT_FILE" | cut -f1)
    echo "Размер файла: $FILE_SIZE"
    
    # Проверка типа файла
    FILE_TYPE=$(file -b "$OUTPUT_FILE")
    echo "Тип файла: $FILE_TYPE"
    
    # Если это не изображение, возможно, в файле содержится сообщение об ошибке
    if [[ ! "$FILE_TYPE" =~ "image" ]]; then
        echo "Предупреждение: Полученный файл не является изображением. Содержимое может быть сообщением об ошибке:"
        cat "$OUTPUT_FILE"
    fi
else
    echo "Ошибка: Не удалось получить обработанное изображение"
fi 