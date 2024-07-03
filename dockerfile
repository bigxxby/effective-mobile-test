# Используем официальный образ Golang 1.22.4
FROM golang:1.22.4

# Установка рабочей директории внутри контейнера
WORKDIR /app

# Копирование всех файлов проекта в контейнер
COPY . .

# Установка PostgreSQL клиента
RUN apt-get update && \
    apt-get install -y postgresql-client

# Копирование .env.example в .env
COPY .env.example .env

# Применение переменных окружения для подключения к PostgreSQL
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=your_user
ENV DB_PASSWORD=your_password
ENV DB_NAME=your_database

# Команда для запуска приложения через go run ./cmd
CMD ["go", "run", "./cmd"]
