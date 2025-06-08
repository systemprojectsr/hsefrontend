#!/usr/bin/env python3
"""
Скрипт для тестирования обработки фотографий в фотосервисе
"""

import argparse
import json
import os
import sys
import requests
from PIL import Image
import io

def main():
    # Парсинг аргументов командной строки
    parser = argparse.ArgumentParser(description='Обработка фотографий через API фотосервиса')
    parser.add_argument('photo_id', help='ID фотографии для обработки')
    parser.add_argument('--output', '-o', default='processed_image.jpg', help='Имя выходного файла (по умолчанию: processed_image.jpg)')
    parser.add_argument('--width', '-w', type=int, default=400, help='Ширина для изменения размера (по умолчанию: 400)')
    parser.add_argument('--height', '-ht', type=int, default=300, help='Высота для изменения размера (по умолчанию: 300)')
    parser.add_argument('--quality', '-q', type=int, default=85, help='Качество JPEG (1-100, по умолчанию: 85)')
    parser.add_argument('--url', '-u', default='http://localhost:8081', help='URL сервиса (по умолчанию: http://localhost:8081)')
    parser.add_argument('--no-resize', action='store_true', help='Не изменять размер изображения')
    parser.add_argument('--verbose', '-v', action='store_true', help='Вывод подробной информации')
    
    args = parser.parse_args()
    
    # URL сервиса
    service_url = args.url
    
    # Проверка статуса сервиса перед выполнением запроса
    print("Проверка статуса сервиса...")
    try:
        status_response = requests.get(f"{service_url}/status", timeout=5)
        status_code = status_response.status_code
    except requests.RequestException as e:
        print(f"Ошибка при проверке статуса сервиса: {e}")
        return 1
    
    if status_code != 200:
        print(f"Ошибка: Сервис недоступен (статус код: {status_code})")
        print(f"Проверьте, запущен ли сервис и доступен ли он по адресу {service_url}")
        return 1
    
    print(f"Сервис доступен. Выполняется запрос на обработку фотографии с ID: {args.photo_id}")
    
    # Параметры обработки
    process_params = {
        "resize": not args.no_resize,
        "width": args.width,
        "height": args.height,
        "quality": args.quality
    }
    
    if args.verbose:
        print(f"Параметры обработки: {json.dumps(process_params, indent=2)}")
        print(f"Файл будет сохранен как: {args.output}")
    
    # Выполнение запроса к API для обработки фотографии
    try:
        response = requests.post(
            f"{service_url}/photos/process?id={args.photo_id}",
            json=process_params,
            headers={"Content-Type": "application/json"},
            timeout=30
        )
        
        if args.verbose:
            print(f"Статус ответа: {response.status_code}")
            print(f"Заголовки ответа: {response.headers}")
        
        # Проверка успешности запроса
        response.raise_for_status()
        
        # Сохранение полученного изображения
        with open(args.output, 'wb') as f:
            f.write(response.content)
        
        # Анализ результата
        file_size = os.path.getsize(args.output)
        file_size_kb = file_size / 1024.0
        
        if file_size > 0:
            print(f"Успешно! Обработанное изображение сохранено в файл: {args.output}")
            print(f"Размер файла: {file_size_kb:.2f} KB")
            
            # Проверка, является ли файл изображением
            try:
                img = Image.open(io.BytesIO(response.content))
                print(f"Тип изображения: {img.format}")
                print(f"Размер изображения: {img.width}x{img.height}")
                print(f"Цветовой режим: {img.mode}")
            except Exception as e:
                print(f"Предупреждение: Полученный файл может не быть изображением: {e}")
                if len(response.content) < 1000:  # Если файл небольшой, вероятно это текст
                    print("Содержимое файла:")
                    print(response.text)
        else:
            print(f"Ошибка: Полученный файл пустой")
            
    except requests.RequestException as e:
        print(f"Ошибка при выполнении запроса: {e}")
        if hasattr(e, 'response') and e.response:
            print(f"Статус код: {e.response.status_code}")
            try:
                error_data = e.response.json()
                print(f"Ответ сервера: {json.dumps(error_data, indent=2)}")
            except ValueError:
                print(f"Ответ сервера: {e.response.text}")
        return 1
    
    return 0

if __name__ == "__main__":
    sys.exit(main()) 