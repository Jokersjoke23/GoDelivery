# GoDelivery

REST API для сервиса доставки еды, написанный на Go.

## Технологии

- **Go** + **Gin** — веб-фреймворк
- **PostgreSQL** + **GORM** — база данных
- **JWT** — авторизация
- **WebSocket** — уведомления в реальном времени
- **Docker** — контейнеризация
- **Swagger** — документация API

## Архитектура

Проект построен на Clean Architecture:
Handler → Service → Repository → PostgreSQL

Также реализован отдельный микросервис уведомлений (`notification-service`).

## Функционал

- JWT авторизация с ролями (customer, courier, admin)
- Управление ресторанами и меню
- Система заказов с полным флоу статусов
- Курьеры — самостоятельный выбор заказов
- Автоматическая смена статусов доставки
- WebSocket уведомления для курьеров (новые заказы)
- WebSocket уведомления для пользователей (статус доставки)
- Восстановление пароля через email
- Экспорт/импорт данных в Excel
- Поиск ресторанов и блюд
- Пагинация
- Unit тесты

## Запуск через Docker

```bash
git clone https://github.com/Jokersjoke23/GoDelivery.git
cd GoDelivery
docker-compose up --build
```

Сервер запустится на `http://localhost:8080`

## Запуск локально

### Требования

- Go 1.21+
- PostgreSQL 16+

### Установка

```bash
git clone https://github.com/Jokersjoke23/GoDelivery.git
cd GoDelivery
cp .env.example .env
# заполни .env своими данными
go mod download
```

### Миграции

```bash
psql -U postgres -d delivery -f migrations/001_init.sql
psql -U postgres -d delivery -f migrations/002_exports.sql
psql -U postgres -d delivery -f migrations/003_delivery_price.sql
psql -U postgres -d delivery -f migrations/004_password_resets.sql
```

### Запуск

```bash
go run cmd/server/main.go
```

## Swagger документация
http://localhost:8080/swagger/index.html

## Notification Microservice

Отдельный сервис для WebSocket уведомлений пользователей.

```bash
cd notification-service
go run cmd/main.go
```

Запускается на порту `:8085`

## Тесты

```bash
go test ./internal/service/... -v
```

## Роли

| Роль | Возможности |
|------|------------|
| customer | Создание заказов, просмотр меню |
| courier | Выбор заказов, обновление статуса доставки |
| admin | Полный доступ, экспорт/импорт данных |

## Структура проекта

GoDelivery/

├── cmd/server/          # Точка входа
├── config/              # Конфигурация
├── internal/
│   ├── domain/          # Модели и DTO
│   ├── handler/         # HTTP хендлеры
│   ├── middleware/       # Auth, Role, Lang middleware
│   ├── repository/      # Работа с БД
│   ├── service/         # Бизнес логика
│   └── websocket/       # WebSocket для курьеров
├── migrations/          # SQL миграции
├── notification-service/ # Микросервис уведомлений
├── pkg/
│   ├── hasher/          # Хеширование паролей
│   ├── jwt/             # JWT токены
│   ├── mailer/          # Отправка email
│   └── response/        # Формат ответов API
├── Dockerfile
└── docker-compose.yml