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

	"go.mongodb.org/mongo-driver/mongo/options"
)

var MenuCollection = database.OpenCollection(database.Client, "menu")

// inTimeSpan checks if the current time is within the start and end time range.
func inTimeSpan(start, end, check time.Time) bool {
	// return check.After(start) && check.Before(end)
	return start.After(time.Now()) && end.After(start)
}
func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := MenuCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			return
		}

		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allMenus)
	}
}
func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
			menu_id := c.Param("menu_id")
			var menu bson.M
			defer cancel()
		err := MenuCollection.FindOne(ctx, bson.M{"menu_id": menu_id}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found"})
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var menu models.Menu
		defer cancel()
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while binding food"})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while validating menu"})
			return
		}
		menu.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Id = primitive.NewObjectID()
		menu.Menu_ID = menu.Id.Hex()
		result, insertErr := MenuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			msg := fmt.Sprintf("menu was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

			if !inTimeSpan(menu.Start_date, menu.End_date, time.Now()) {
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			if !inTimeSpan(menu.Start_date, menu.End_date, time.Now()) {
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			var updatedObj bson.D
			updatedObj = append(updatedObj, bson.E{Key: "starting_date", Value: menu.Start_date})
			updatedObj = append(updatedObj, bson.E{Key: "ending_date", Value: menu.End_date})
			if menu.Menu_Name != "" {
				updatedObj = append(updatedObj, bson.E{Key: "menu_name", Value: menu.Menu_Name})
			}
			if menu.Category != "" {
				updatedObj = append(updatedObj, bson.E{Key: "category", Value: menu.Category})
			}
			menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: menu.Updated_At})

			upsert := true
			opt := options.UpdateOptions{
				Upsert: &upsert,
			}
			filter := bson.M{"menu_id": menu.Menu_ID}
			updateObj := bson.D{
				{Key: "$set", Value: updatedObj},
			}

			result, err := MenuCollection.UpdateOne(
				ctx,
				filter,
				updateObj,
				&opt,
			)
			if err != nil {
				msg := "menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			c.JSON(http.StatusOK, result)
		}
	}

