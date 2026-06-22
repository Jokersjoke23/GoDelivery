package handler

import (
	"delivery-app/docs"
	"delivery-app/internal/middleware"
	"delivery-app/pkg/jwt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	auth       *AuthHandler
	user       *UserHandler
	restaurant *RestaurantHandler
	order      *OrderHandler
	courier    *CourierHandler
	delivery   *DeliveryHandler
	export     *ExportHandler
	jwt        jwt.JWT
	ws         *WSHandler
}

func NewRouter(
	auth *AuthHandler,
	user *UserHandler,
	restaurant *RestaurantHandler,
	order *OrderHandler,
	courier *CourierHandler,
	delivery *DeliveryHandler,
	export *ExportHandler,
	jwtManager jwt.JWT,
	ws *WSHandler,
) *Router {
	return &Router{
		auth: auth, user: user, restaurant: restaurant, order: order,
		courier: courier, delivery: delivery, export: export, jwt: jwtManager, ws: ws,
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api"

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api")

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", r.auth.Register)
		authGroup.POST("/login", r.auth.Login)
		authGroup.POST("/forgot-password", r.auth.ForgotPassword)
		authGroup.POST("/reset-password", r.auth.ResetPassword)
	}

	api.GET("/restaurants", r.restaurant.GetAll)
	api.GET("/restaurants/search", r.restaurant.Search)
	api.GET("/restaurants/:id", r.restaurant.GetByID)
	api.GET("/restaurants/menu/search", r.restaurant.SearchMenu)
	api.GET("/restaurants/:id/menu", r.restaurant.GetMenu)
	api.GET("/restaurants/menu/:itemID", r.restaurant.GetMenuItem)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(r.jwt))
	{
		protected.GET("users/me", r.user.GetMe)
		protected.PUT("users/me", r.user.UpdateMe)

		protected.POST("orders", r.order.Create)
		protected.GET("orders/my", r.order.GetMyOrders)
		protected.GET("orders/:id", r.order.GetByID)
		protected.POST("orders/:id/cancel", r.order.Cancel)
		protected.GET("orders/:id/delivery", r.delivery.GetByOrderID)

		protected.GET("ws", r.ws.Connect)

		protected.POST("restaurants", r.restaurant.Create)
		protected.PUT("restaurants/:id", r.restaurant.Update)
		protected.DELETE("restaurants/:id", r.restaurant.Delete)
		protected.POST("restaurants/:id/menu", r.restaurant.AddMenuItem)
		protected.PUT("restaurants/menu/:itemID", r.restaurant.UpdateMenuItem)
		protected.DELETE("restaurants/menu/:itemID", r.restaurant.DeleteMenuItem)
		protected.GET("restaurants/:id/orders", r.order.GetByRestaurant)

		protected.POST("couriers", r.courier.Create)
		protected.GET("couriers/me", r.courier.GetMe)
		protected.PUT("couriers/me/status", r.courier.UpdateStatus)
		protected.GET("couriers/orders/available", r.courier.GetAvailableOrders)
		protected.POST("couriers/orders/:id/take", r.courier.TakeOrder)
		protected.GET("couriers/me/location", r.courier.GetMyLocation)
		protected.PUT("couriers/me/location", r.courier.UpdateLocation)
		protected.GET("couriers/me/deliveries", r.delivery.GetMyCourierDeliveries)
		protected.DELETE("couriers/me", r.courier.DeleteProfile)

		protected.GET("deliveries/:id", r.delivery.GetByID)
		protected.POST("deliveries/:id/next-status", r.delivery.NextStatus)

		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("users", r.user.GetAll)
			admin.DELETE("users/:id", r.user.Delete)
			admin.PUT("orders/:id/status", r.order.UpdateStatus)
			admin.GET("orders", r.order.GetAll)
			admin.GET("couriers/online", r.courier.GetAllOnline)
			admin.PUT("users/:id/role", r.user.AssignRole)
			admin.POST("exports", r.export.Create)
			admin.GET("exports", r.export.GetAll)
			admin.GET("exports/:id", r.export.GetByID)
			admin.GET("exports/:id/download", r.export.Download)
			admin.POST("restaurants/import", r.restaurant.ImportFromExcel)
			admin.POST("restaurants/menu/import", r.restaurant.ImportMenuFromExcel)
			admin.PUT("deliveries/:id/status", r.delivery.UpdateStatus)
		}
	}
}
