package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Order struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	OrderDate time.Time `bson:"order_date,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
	Order_ID string `bson:"order_id,omitempty"`
	Table_ID string `bson:"table_id,omitempty"`
}
