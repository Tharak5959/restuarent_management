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
type OrderItemPack struct{
	Table_ID primitive.ObjectID `json:"table_id" validate:"required"`
	Order_Items []models.OrderItem `json:"order_items" validate:"required"`
}
var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "order_items")
func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	result,err:=	orderItemCollection.Find(context.TODO(), bson.M{})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching order items"})
		return
	}
	var allOrderItems []bson.M
	if err = result.All(ctx, &allOrderItems); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, allOrderItems)	
	}
}
func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		OrderItem_ID := c.Param("order_item_id")
		var orderItem models.OrderItem
		defer cancel()
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": OrderItem_ID}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching Order Item"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
		defer cancel()

	}
}
func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		Order_ID := c.Param("order_id")
		allOrderItems ,err := OrderItemsByOrder(Order_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error": "Error fetching order items"})

			return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}
 func OrderItemsByOrder(id string)(OrderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "food"},
		{Key: "localField", Value: "food_id"},
		{Key: "foreignField", Value: "food_id"},
		{Key: "as", Value: "food"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "order"},
		{Key: "localField", Value: "order_id"},
		{Key: "foreignField", Value: "order_id"},
		{Key: "as", Value: "order"},
	}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$order"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}
	lookupTablestage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "table"},
		{Key: "localField", Value: "order.table_id"},
		{Key: "foreignField", Value: "table_id"},
		{Key: "as", Value: "table"},
	}}}
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$table"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	projectStage := bson.D{
	{Key:"$project", Value:bson.D{
		{Key: "_id", Value: 0},
		{Key: "ammount", Value: "$food.price"},
		{Key: "total_count", Value: 1},
		{Key: "food_name", Value: "$food.food_name"},
		{Key: "food_image", Value: "$food.food_image"},
		{Key: "table_number", Value: "$table.table_number"},
		{Key: "table_id", Value: "$table.table_id"},
		{Key: "order_id", Value: "$order.order_id"},
		{Key: "quantity", Value: "$quantity"},
		{Key: "order_date", Value: "$order.order_date"},
	}},
}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "order_id", Value: "$order_id"},
			{Key: "table_id", Value: "$table_id"},
			{Key: "table_number", Value: "$table_number"},
		}},
		{Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$ammount"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
	}}}
	projectStage2 := bson.D{
		{Key:"$project", Value:  bson.D{
			{Key: "_id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: "$order_items"},
		}},
	}
	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTablestage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err) 
	}
	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}
	defer cancel()
	return OrderItems, nil
}



func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack OrderItemPack
		var order models.Order
		defer cancel()
		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if validationErr := validate.Struct(orderItemPack); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		
		order.Table_ID = orderItemPack.Table_ID.Hex()
		Order_Id := OrderItemOrderCreator(order)
		var orderItemToBeInserted []interface{}
		for _, orderItem := range orderItemPack.Order_Items {
			orderItem.Order_Id = Order_Id		
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.CreatedAt, _ = time.Parse((time.RFC3339), time.Now().Format(time.RFC3339))
			orderItem.UpdatedAt, _ = time.Parse((time.RFC3339), time.Now().Format(time.RFC3339))
			orderItem.Id = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.Id.Hex()
			var num = toFixed(orderItem.Unitprice, 2)
			orderItem.Unitprice = num
			orderItemToBeInserted = append(orderItemToBeInserted, orderItem)
		}

		_, err := orderItemCollection.InsertMany(ctx, orderItemToBeInserted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting order items"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"message": "Order items created successfully"})
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		var orderItem models.OrderItem
		defer cancel()
		orderItemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderItemId}

		var updateObj primitive.D
		if orderItem.Unitprice !=0{
			updateObj = append(updateObj, bson.E{Key: "unitprice", Value: orderItem.Unitprice,})
		} 
		if orderItem.Quantity !=0{
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
		}
		if orderItem.Food_Id != "" {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.Food_Id})
		}

		orderItem.UpdatedAt, _ = time.Parse((time.RFC3339), time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})
		upsert  := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
			
		}
	result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating order item"})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	
	}
}
