package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Table struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	NumberofGuests *int `json:"numberoftables" bson:"numberoftables"`
	TableNumber *int `json:"tablenumber" bson:"tablenumber"`
	CreatedAt time.Time `json:"createdat" bson:"createdat"`
	UpdatedAt time.Time `json:"updatedat" bson:"updatedat"`
	Table_Id string `json:"table_id" bson:"table_id"`
}