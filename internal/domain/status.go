package domain

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleCustomer UserRole = "customer"
	UserRoleCourier  UserRole = "courier"
)

type RestaurantStatus string

const (
	RestaurantStatusActive   RestaurantStatus = "active"
	RestaurantStatusInactive RestaurantStatus = "inactive"
	RestaurantStatusClosed   RestaurantStatus = "closed"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusAccepted  OrderStatus = "accepted"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusPickedUp  OrderStatus = "picked_up"
	OrderStatusOnTheWay  OrderStatus = "on_the_way"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRejected  OrderStatus = "rejected"
)

type DeliveryStatus string

const (
	DeliveryStatusWaiting   DeliveryStatus = "waiting"
	DeliveryStatusAssigned  DeliveryStatus = "assigned"
	DeliveryStatusPickedUp  DeliveryStatus = "picked_up"
	DeliveryStatusOnTheWay  DeliveryStatus = "on_the_way"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusFailed    DeliveryStatus = "failed"
)

type CourierStatus string

const (
	CourierStatusOnline  CourierStatus = "online"
	CourierStatusOffline CourierStatus = "offline"
	CourierStatusBusy    CourierStatus = "busy"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodOnline PaymentMethod = "online"
)
