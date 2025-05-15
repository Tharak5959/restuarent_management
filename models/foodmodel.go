package models
import (
	 "time"
	 "go.mongodb.org/mongo-driver/bson/primitive"
)
type Food struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" validate:"required,min=2,max=100" bson:"name"`
	Price float64 `json:"price" validate:"required" bson:"price"`
	Food_Image string `json:"image" bson:"image"`
	Created_At time.Time `json:"created_at" bson:"created_at"`
	Updated_At time.Time `json:"updated_at" bson:"updated_at"`
	Food_ID string `json:"food_id" bson:"food_id"`
	Menu_ID string `json:"menu_id" validate:"required" bson:"menu_id"`
}
	
		