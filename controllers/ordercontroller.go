package controller

import (
	"context"
	"fmt"
	"golang-restuarent_management/database"
	"golang-restuarent_management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel =context.WithTimeout(context.Background(), 10*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching orders"})
			return
		}
		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrders)
		defer cancel()

	}
}
func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		Order_ID := c.Param("order_id")
		var order models.Order
		defer cancel()
		err := orderCollection.FindOne(ctx, bson.M{"order_id": Order_ID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching Order"})
	}
		c.JSON(http.StatusOK, order)
	}
}
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var table models.Table
		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if validationErr := validate.Struct(order); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		if order.Table_ID != "" {
			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_ID}).Decode(&table)
			defer cancel()
			if err != nil {
				msg := fmt.Sprintf("table with id %v not found", order.Table_ID)
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			order.Id = primitive.NewObjectID()
			order.Order_ID = order.Id.Hex()
			result,inseterr := orderCollection.InsertOne(ctx, order)
			if inseterr != nil {
				msg := "order was not created"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
					defer cancel()
			c.JSON(http.StatusOK, result)
		}
	}
}
func UpdateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var table models.Table
		var order models.Order
		var updateObj primitive.D
		orderId := c.Param("order_id")
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if order.Table_ID != "" {
			err := menuCollection.FindOne(ctx, bson.M{"table_id": order.Table_ID}).Decode(&table)
			defer cancel()
		if err != nil {
			msg:=fmt.Sprintf("order with id %v not found", orderId)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		updateObj = append(updateObj, bson.E{Key: "menu", Value: order.Table_ID})
	}
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})
	upsert := true
	filter := bson.M{"order_id": orderId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	result, err := orderCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	if err != nil {
		msg := fmt.Sprintf("order with id %v not updated", orderId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)

}
}
func OrderItemOrderCreator(order models.Order) string{
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Id = primitive.NewObjectID()
	order.Order_ID = order.Id.Hex()
	orderCollection.InsertOne(ctx, order)
	defer cancel()
	return order.Order_ID

} 