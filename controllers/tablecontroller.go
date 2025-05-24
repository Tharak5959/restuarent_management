package controller

import (
	"context"
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
var tableColl *mongo.Collection = database.OpenCollection(database.Client, "table")
func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel =context.WithTimeout(context.Background(), 10*time.Second)

		result, err := tableColl.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Tables"})
			return
		}
		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allTables)
		defer cancel()
	}
}
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		Table_ID := c.Param("table_id")
		var Table models.Table
		defer cancel()
		err := orderCollection.FindOne(ctx, bson.M{"table_id": Table_ID}).Decode(&Table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching Table"})
			return
	}
		c.JSON(http.StatusOK, Table)
	}
}
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if validationErr := validate.Struct(table); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		defer cancel()
		table.ID = primitive.NewObjectID()
		table.Table_Id = table.ID.Hex()
		table.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		result, err := tableColl.InsertOne(ctx, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while creating Table"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var table models.Table
		Table_ID := c.Param("table_id")
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		var updateObj primitive.D
		if table.NumberofGuests != nil {
			updateObj = append(updateObj, bson.E{Key: "numberoftables", Value: table.NumberofGuests})
		}
		if table.TableNumber != nil {
			updateObj = append(updateObj, bson.E{Key: "tablenumber", Value: table.TableNumber})
		}

		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		Upsert := true
		filter := bson.M{"table_id": Table_ID}
		opts := options.UpdateOptions{
			Upsert: &Upsert,
		}
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.Updated_At})
		result, err := tableColl.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},

			&opts,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating Table"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}
