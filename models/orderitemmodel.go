package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type OrderItem struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Quantity int `json:"quantity" validate:"required,min=1,max=100" bson:"quantity"`
	Unitprice float64 `json:"unitprice" validate:"required,min=1,max=100" bson:"unitprice"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Food_Id string `json:"food_id" bson:"food_id"`
	Order_item_id string `json:"order_item_id" bson:"order_item_id"`
	Order_Id string `json:"order_id" bson:"order_id"`
}