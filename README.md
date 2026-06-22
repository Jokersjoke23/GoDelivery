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
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── domain/
│   │   ├── courier.go
│   │   ├── courier_dto.go
│   │   ├── delivery.go
│   │   ├── delivery_dto.go
│   │   ├── export.go
│   │   ├── order.go
│   │   ├── order_dto.go
│   │   ├── password_reset.go
│   │   ├── restaurant.go
│   │   ├── restaurant_dto.go
│   │   ├── status.go
│   │   ├── user.go
│   │   └── user_dto.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── courier_handler.go
│   │   ├── delivery_handler.go
│   │   ├── export_handler.go
│   │   ├── order_handler.go
│   │   ├── restaurant_handler.go
│   │   ├── router.go
│   │   ├── user_handler.go
│   │   └── ws_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── lang.go
│   │   └── role.go
│   ├── repository/
│   │   ├── courier_repo.go
│   │   ├── delivery_repo.go
│   │   ├── export_repo.go
│   │   ├── menu_item_repo.go
│   │   ├── order_repo.go
│   │   ├── password_reset_repo.go
│   │   ├── restaurant_repo.go
│   │   └── user_repo.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── auth_service_test.go
│   │   ├── courier_service.go
│   │   ├── courier_service_test.go
│   │   ├── delivery_service.go
│   │   ├── errors_test.go
│   │   ├── export_service.go
│   │   ├── order_service.go
│   │   ├── order_service_test.go
│   │   └── restaurant_service.go
│   └── websocket/
│       ├── client.go
│       └── hub.go
├── migrations/
│   ├── 001_init.sql
│   ├── 002_exports.sql
│   ├── 003_delivery_price.sql
│   └── 004_password_resets.sql
├── notification-service/
│   ├── cmd/
│   │   └── main.go
│   ├── domain/
│   │   └── notification.go
│   ├── handler/
│   │   └── handler.go
│   └── hub/
│       ├── client.go
│       └── hub.go
├── pkg/
│   ├── hasher/
│   │   └── hasher.go
│   ├── jwt/
│   │   └── jwt.go
│   ├── mailer/
│   │   └── mailer.go
│   └── response/
│       └── response.go
├── docs/
├── exports/
├── imports/
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md