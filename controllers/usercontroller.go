package controller

import (
	"context"
	"fmt"


	"golang-restuarent_management/database"
	helper "golang-restuarent_management/helpers"
	"golang-restuarent_management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var user models.User
		//convert the json data from postman to something that we can use
		if err := c.BindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":"error while binding json"})
			return
		}

		//validate the data  based on the struct
		validateErr := validate.Struct(user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}

		//check if the email already exists
		count,err:=userCollection.CountDocuments(ctx,bson.M{"email":user.Email})
		defer cancel()
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while checking for the email"})
			return
		}

		//hash the password
		hashedPassword, hashErr := HashPassword(user.Password)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while hashing password"})
			return
		}
		user.Password = hashedPassword

		//phone number  is alredy used by another user or not
		count,err=userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while checking for the phone number"})
			return
		}
		if count>0{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"phone number already exists"})
			return
		}

		//get some extra details from the user, -createdat,updtaedat,id
		user.CreatedAt,_=time.Parse((time.RFC3339),time.Now().Format(time.RFC3339))
		user.UpdatedAt,_=time.Parse((time.RFC3339),time.Now().Format(time.RFC3339))
		user.Id=primitive.NewObjectID()
		user.User_Id=user.Id.Hex()

		//genarate token and refresh token
	token,refreshToken,_:=	helper.GenarateAllTokens(*&user.Email,*&user.User_Id,*&user.FirstName,*&user.LastName)
		user.Token=token
		user.Refresh_token=refreshToken

		//if all ok,then you insert new user to the database
		resultInsertionNumber,err:=userCollection.InsertOne(ctx,user)
		if err != nil {
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		//result status ok andd send the resullt back to the user
		c.JSON(http.StatusOK, gin.H{"data": resultInsertionNumber})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
		var user models.User
		var foundUser models.User
		//convert the json data from postman to something that we can use
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":"error while binding json"})
			return
		}
		//find the user by email and seee if that user exists
		err:=userCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error while finding user"})
			return
		}

		//then u will verify password
		passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		// if all goes well,then you'll generate 
		token,refreshtoken,_:=helper.GenarateAllTokens(*&foundUser.Email,*&foundUser.User_Id,*&foundUser.FirstName,*&foundUser.LastName)
		//update and refresh token
		helper.UpdateAllTokens(token,refreshtoken,foundUser.User_Id)
		//retrun status ok
		c.JSON(http.StatusOK, gin.H{"data": foundUser})

	}
}
func HashPassword(password string) (string, error) {
	bytes,err:= bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes), err

}
func VerifyPassword(userPassword string, providedpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedpassword), []byte(userPassword))
	check := true
	msg := "login success"
	if err != nil {
		msg = fmt.Sprintf("login failed")
		check = false
	}
	return check, msg

}