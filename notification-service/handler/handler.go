package handler

import (
	"encoding/json"
	"net/http"
	"notification-service/domain"
	"notification-service/hub"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *hub.Hub
}

func NewHandler(hub *hub.Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) ConnectUser(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := hub.NewClient(userID, h.hub, conn)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) Notify(c *gin.Context) {
	var n domain.Notification
	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload, _ := json.Marshal(n)
	h.hub.Notify(n.UserID, payload)

	c.JSON(http.StatusOK, gin.H{"message": "уведомление отправлено"})
}
