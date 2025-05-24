package controller

import (
	"context"
	"golang-restuarent_management/database"
	"golang-restuarent_management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		recordPerPage,err:=strconv.Atoi(c.Query("recordpage"))
		if err!=nil||recordPerPage<1{
			recordPerPage=10
		}
		page,err1:=strconv.Atoi(c.Query("page"))
		if err1!=nil||page<1{
			page=1
		}
		startIndex := (page-1)*recordPerPage
		if queryStartIndex := c.Query("startIndex"); queryStartIndex != "" {
			startIndex, _ = strconv.Atoi(queryStartIndex)
		}

		matchstage := bson.D{{Key: "$match", Value: bson.D{}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "totalCount", Value: 1},
				{Key: "user_items", Value: bson.D{
					{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
				}},
			}},
		}
		result,err:=userCollection.Aggregate(ctx,mongo.Pipeline{
			matchstage,projectStage})
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while listing users"})
			return
		}
		var allUsers []bson.M
		if err=result.All(ctx,allUsers);err!=nil{
			log.Fatal(err)
	
		}
		c.JSON(http.StatusOK,gin.H{"data":allUsers[0]})
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		userId := c.Param("user_id")
		var user models.User
		err := userCollection.FindOne(ctx,bson.M{"user_id":userId}).Decode(&user)
		
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while getting user"})
			return
		}
		c.JSON(http.StatusOK,gin.H{"data":user})
	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//convert the json data from postman to something that we can use
		//validate the data  based on the struct
		//check if the email already exists
		//hash the password
		//phone number  is alredy used by another user or not
		//get some extra details from the user, -createdat,updtaedat,id
		//genarate token and refresh token
		//if all ok,then you insert new user to the database
		//result status ok andd send the resullt back to the user
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//convert the json data from postman to something that we can use
		//find the user by email and seee if that user exists
		//then u will verify password
		//if all goes well,then you'll generate 
		//update and refresh token
		//retrun status ok

	}
}
func HashPassword(password string) (string, error) {
	// Hashing logic here
	return password, nil // Replace with actual hashing logic
}
func VerifyPassword(userPassword string, providedpassword string) (bool, string) {
	// Password verification logic here
	return true, "" // Replace with actual verification logic
}