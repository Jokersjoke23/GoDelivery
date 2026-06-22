package main

import (
	"log"
	"notification-service/handler"
	"notification-service/hub"

	"github.com/gin-gonic/gin"
)

func main() {
	h := hub.NewHub()
	go h.Run()

	hdl := handler.NewHandler(h)

	engine := gin.Default()

	engine.GET("/ws/user", hdl.ConnectUser)
	engine.POST("/notify", hdl.Notify)

	log.Println("notification-service запущен на порту 8085")
	if err := engine.Run(":8085"); err != nil {
		log.Fatalf("ошибка запуска: %v", err)
	}
}
