package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Username string `bson:"username,omitempty"`
	FirstName string `bson:"first_name,omitempty"`
	LastName string `bson:"last_name,omitempty"`
	Avatar string `bson:"avatar,omitempty"`
	Phone string `bson:"phone,omitempty"`
	Email string `bson:"email,omitempty"`
	Password string `bson:"password,omitempty"`
	Token string `bson:"token,omitempty"`
	Refresh_token string `bson:Refresh_token`
	Role string `bson:"role,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
	User_Id string `bson:"user_id,omitempty"`
}