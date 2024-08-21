# Используем официальный образ PostgreSQL
FROM postgres:latest

# Устанавливаем переменные окружения для настройки PostgreSQL
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=admin
ENV POSTGRES_DB=test

# Открываем порт 5432 для внешнего доступа
EXPOSE 5432

# Запускаем PostgreSQL
CMD ["postgres"]
