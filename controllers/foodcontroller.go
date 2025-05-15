package controller

import (
	"context"
	"fmt"
	"golang-restuarent_management/database"
	"golang-restuarent_management/models"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerpage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		startIndex,err= strconv.Atoi(c.Query("startIndex"))
		matchStage := primitive.D{{Key: "$match", Value: primitive.D{}}}
		GroupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
				{Key: "totalcount", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
			}},
		}
		projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "totalcount", Value: 1},
				{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			},
		},
	}
	result ,err :=foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, GroupStage, projectStage})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching food"})
		return
	}
	var allFoods []bson.M
	if err =result.All(ctx, &allFoods); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, allFoods[0])
}
}
func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Handler logic here
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		foodId := c.Param("foodId")
		var food models.Food
		defer cancel()
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching food"})
	}
		c.JSON(http.StatusOK, food)
	}
}
func CreateFood() gin.HandlerFunc{
	return func(c *gin.Context) {
		// Handler logic here
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var food models.Food
		var menu models.Menu
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while binding food"})
			return
		}

		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while validating food"})
			return
		}

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_ID}).Decode(&menu)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		food.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Food_ID = primitive.NewObjectID().Hex()
		var num = toFixed(food.Price, 2)
		food.Price = num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("food was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return int(num+math.Copysign(0.5, num))
}
func toFixed(num float64, precision int) float64 {
	output:=math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}
func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		var food models.Food
		food.Food_ID = c.Param("food_id")
		if err :=c.BindJSON(&menu);err!=nil {
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		var UpdateObj primitive.D
		if food.Name != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "name", Value: food.Name})
		}
		if food.Food_Image != "" {
			UpdateObj = append(UpdateObj, bson.E{Key: "food_image", Value: food.Food_Image})
		}
		if food.Price != 0 {
			UpdateObj = append(UpdateObj, bson.E{Key: "price", Value: food.Price})
		}
		if food.Menu_ID != "" {
			err := menuCollection.FindOne(ctx, bson.E{Key: "menu_id", Value: food.Menu_ID}).Decode(&menu)
			defer cancel()
			if err != nil {

				msg:= "menu not found"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			UpdateObj = append(UpdateObj, bson.E{Key: "menu",Value: food.Price}) 
		}
		food.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		UpdateObj = append(UpdateObj, bson.E{Key: "updated_at", Value: food.Updated_At})
		upsert := true
		filter := bson.M{"food_id": food.Food_ID}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := foodCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: UpdateObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("food was not updated")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, result)
			

}}
// func DeleteFood() gin.HandlerFunc {

// }