package main

import (
	"context"

	middleware "golang-restuarent_management/middleware"
	// "golang-restuarent_management/mongo"
	routes "golang-restuarent_management/routes"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
 func main(){
	port  := os.Getenv("PORT")
	if port ==""{	
		port="8080"
 	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Autentication())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	// router.itemRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":"+port)
 }

var Client *mongo.Client

func init() {
    // Initialize MongoDB client
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    Client = client
}