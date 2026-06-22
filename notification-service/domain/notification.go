package domain

type NotificationType string

const (
	NotificationOrderAssigned  NotificationType = "order_assigned"
	NotificationOrderPickedUp  NotificationType = "order_picked_up"
	NotificationOrderOnTheWay  NotificationType = "order_on_the_way"
	NotificationOrderDelivered NotificationType = "order_delivered"
)

type Notification struct {
	Type    NotificationType `json:"type"`
	UserID  string           `json:"user_id"`
	OrderID string           `json:"order_id"`
	Message string           `json:"message"`
}
