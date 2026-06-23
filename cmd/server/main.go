// @title           GoDelivery API
// @version         1.0
// @description     REST API для сервиса доставки еды
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"delivery-app/config"
	"delivery-app/internal/domain"
	"delivery-app/internal/handler"
	"delivery-app/internal/middleware"
	"delivery-app/internal/repository"
	"delivery-app/internal/service"
	"delivery-app/internal/websocket"
	"delivery-app/pkg/hasher"
	"delivery-app/pkg/jwt"
	"delivery-app/pkg/mailer"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден")
	}

	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	log.Println("подключение к БД успешно")

	os.MkdirAll("exports", 0755)
	os.MkdirAll("imports", 0755)

	jwtManager := jwt.NewJWT(cfg.JWT.SecretKey, cfg.JWT.ExpiresIn)
	passwordHasher := hasher.NewHasher()

	userRepo := repository.NewUserRepository(db)
	restaurantRepo := repository.NewRestaurantRepository(db)
	menuItemRepo := repository.NewMenuItemRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	courierRepo := repository.NewCourierRepository(db)
	deliveryRepo := repository.NewDeliveryRepository(db)
	exportRepo := repository.NewExportRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	if err := db.AutoMigrate(&domain.PasswordReset{}); err != nil {
		log.Fatalf("ошибка миграции password_resets: %v", err)
	}

	hub := websocket.NewHub()
	go hub.Run()

	mailClient := mailer.NewMailer(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Password,
		cfg.SMTP.From,
	)

	authSvc := service.NewAuthService(userRepo, passwordResetRepo, passwordHasher, jwtManager, mailClient)
	userSvc := service.NewUserService(userRepo)
	restaurantSvc := service.NewRestaurantService(restaurantRepo, menuItemRepo)
	orderSvc := service.NewOrderService(orderRepo, restaurantRepo, menuItemRepo, courierRepo, deliveryRepo, hub)
	courierSvc := service.NewCourierService(courierRepo, userRepo)
	deliverySvc := service.NewDeliveryService(deliveryRepo, orderRepo, courierRepo)
	exportSvc := service.NewExportService(exportRepo, orderRepo, courierRepo, restaurantRepo, menuItemRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	restaurantHandler := handler.NewRestaurantHandler(restaurantSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	courierHandler := handler.NewCourierHandler(courierSvc, orderSvc, deliverySvc)
	deliveryHandler := handler.NewDeliveryHandler(deliverySvc, courierSvc)
	exportHandler := handler.NewExportHandler(exportSvc)
	wsHandler := handler.NewWSHandler(hub, courierSvc)

	router := handler.NewRouter(
		authHandler,
		userHandler,
		restaurantHandler,
		orderHandler,
		courierHandler,
		deliveryHandler,
		exportHandler,
		jwtManager,
		wsHandler,
	)

	engine := gin.Default()
	engine.Use(middleware.LangMiddleware())
	router.Setup(engine)

	log.Printf("сервер запущен на порту %s", cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
