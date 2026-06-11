package handler

import (
	"delivery-app/internal/service"
	ws "delivery-app/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub            *ws.Hub
	courierService service.CourierService
}

func NewWSHandler(hub *ws.Hub, courierService service.CourierService) *WSHandler {
	return &WSHandler{hub: hub, courierService: courierService}
}

func (h *WSHandler) Connect(c *gin.Context) {
	userID := c.GetString("userID")

	courier, err := h.courierService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "только для курьеров"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(h.hub, conn, courier.ID)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}
