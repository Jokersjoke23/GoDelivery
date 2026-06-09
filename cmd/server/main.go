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
	"delivery-app/internal/handler"
	"delivery-app/internal/middleware"
	"delivery-app/internal/repository"
	"delivery-app/internal/service"
	"delivery-app/pkg/hasher"
	"delivery-app/pkg/jwt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("ошибка подключения к БД: %v", err)
	}
	log.Println("подключение к БД успешно")

	jwtManager := jwt.NewJWT(cfg.JWT.SecretKey, cfg.JWT.ExpiresIn)
	passwordHasher := hasher.NewHasher()

	userRepo := repository.NewUserRepository(db)
	restaurantRepo := repository.NewRestaurantRepository(db)
	menuItemRepo := repository.NewMenuItemRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	courierRepo := repository.NewCourierRepository(db)
	deliveryRepo := repository.NewDeliveryRepository(db)

	authSvc := service.NewAuthService(userRepo, passwordHasher, jwtManager)
	userSvc := service.NewUserService(userRepo)
	restaurantSvc := service.NewRestaurantService(restaurantRepo, menuItemRepo)
	orderSvc := service.NewOrderService(orderRepo, restaurantRepo, menuItemRepo, courierRepo, deliveryRepo)
	courierSvc := service.NewCourierService(courierRepo, userRepo)
	deliverySvc := service.NewDeliveryService(deliveryRepo, orderRepo, courierRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	restaurantHandler := handler.NewRestaurantHandler(restaurantSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	courierHandler := handler.NewCourierHandler(courierSvc, orderSvc, deliverySvc)
	deliveryHandler := handler.NewDeliveryHandler(deliverySvc, courierSvc)

	router := handler.NewRouter(
		authHandler,
		userHandler,
		restaurantHandler,
		orderHandler,
		courierHandler,
		deliveryHandler,
		jwtManager,
	)

	engine := gin.Default()
	engine.Use(middleware.LangMiddleware())
	router.Setup(engine)

	log.Printf("сервер запущен на порту %s", cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
