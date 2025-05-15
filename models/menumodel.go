package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Menu struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Menu_ID string `json:"menu_id" bson:"menu_id"`
	Menu_Name string `json:"menu_name" validate:"required,min=2,max=100" bson:"menu_name"`
	Category string `json:"category" validate:"required,min=2,max=100" bson:"category"`
	Start_date time.Time `json:"starting_date" validate:"required" bson:"starting_date"`
	End_date time.Time `json:"ending_date" validate:"required" bson:"ending_date"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}