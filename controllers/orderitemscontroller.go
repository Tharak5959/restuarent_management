package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handler logic here
	}
}
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handler logic here
	}
}
func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handler logic here
	}
}
//  func ItemsByOrder(id string)(OrderItems []primitive.M, err error) {
//  }

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handler logic here
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handler logic here
	}
}
