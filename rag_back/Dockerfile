# Используем базовый образ Python 3.10 slim
FROM python:3.12

WORKDIR /app


# Копируем файл с зависимостями в контейнер
COPY requirements.txt .

# Устанавливаем зависимости Python
RUN pip install --no-cache-dir -r requirements.txt

# Копируем все файлы проекта в контейнер
COPY . .

