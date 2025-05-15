package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"golang-restuarent_management/database"
	"golang-restuarent_management/routes"
	"golang-restuarent_management/middleware"
	// "golang-restuarent_management/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)
var foodcollection *mongo.collection - database.Opencollection(database.client,"food")
 func main(){
	port  := os.Getenv("PORT")
	if port ==""{	
		port="8080"
 	}
	router := gin.new()
	router.Use(gin.Logger())
	router.UserRoutes(router)
	router.Use(middleware.Autentication())
	router.foodRoutes(router)
	router.menuRoutes(router)
	router.tableRoutes(router)
	router.itemRoutes(router)
	router.orderRoutes(router)
	router.orderitemRoutes(router)
	router.invoiceRoutes(router)

	router.Run(":"+port)
 }