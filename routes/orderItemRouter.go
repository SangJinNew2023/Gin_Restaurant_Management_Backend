package routes

import (
	controller "golang-restaurant-management/controllers" //모듈 이름 지정

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orderItems", controller.GetOrderItems())
	incomingRoutes.GET("/orderItems/:orderItem_id", controller.GetOrderItem())
	incomingRoutes.GET("/orderItems-order/:order_id", controller.GetOrderItemByOrder())
	incomingRoutes.POST("/orderItems", controller.CreateOrderItem())
	incomingRoutes.PATCH("/orderItems/:orderItem_id", controller.UpdateOrderItem())
}
