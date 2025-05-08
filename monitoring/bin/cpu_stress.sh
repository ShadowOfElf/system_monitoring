#!/bin/bash

if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Использование: $0 <количество_потоков> <время_в_секундах>"
  echo "Пример: $0 4 10"
  exit 1
fi

THREADS=$1
TIME=$2

echo "Запуск нагрузки на $THREADS потока(ов) на $TIME секунд..."

stress_cpu() {
  while true; do :; done
}

for ((i = 1; i <= THREADS; i++)); do
  stress_cpu &
  echo "Запущен поток $i с PID $!"
done

sleep $TIME
echo "Время истекло. Останавливаем нагрузку..."
killall cpu_stress.sh
