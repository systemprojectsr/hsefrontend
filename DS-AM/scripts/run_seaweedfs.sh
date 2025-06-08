#!/bin/bash

# Скрипт для запуска SeaweedFS в локальном окружении

# Создаем директории для хранения данных SeaweedFS
mkdir -p .data/master
mkdir -p .data/volume1
mkdir -p .data/filer

# Запускаем SeaweedFS Master
weed master -ip=localhost -port=9333 -dir=.data/master -mdir=.data/master &
MASTER_PID=$!
echo "Started SeaweedFS Master with PID: $MASTER_PID"

# Даем время мастеру запуститься
sleep 2

# Запускаем Volume Server
weed volume -ip=localhost -port=8080 -dir=.data/volume1 -max=5 -mserver=localhost:9333 -dataCenter=dc1 &
VOLUME_PID=$!
echo "Started SeaweedFS Volume Server with PID: $VOLUME_PID"

# Запускаем Filer
weed filer -ip=localhost -port=8888 -master=localhost:9333 -dir=.data/filer &
FILER_PID=$!
echo "Started SeaweedFS Filer with PID: $FILER_PID"

echo "SeaweedFS is running. Press Ctrl+C to stop."

# Функция для корректного завершения
cleanup() {
    echo "Stopping SeaweedFS services..."
    kill $FILER_PID
    kill $VOLUME_PID
    kill $MASTER_PID
    wait
    echo "SeaweedFS services stopped."
    exit 0
}

# Перехватываем сигнал завершения
trap cleanup SIGINT SIGTERM

# Ждем завершения
wait 