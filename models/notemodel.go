package models
import(
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type note struct{
	Id	primitive.ObjectID `bson:"_id"`
	Text string `json:"text"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time	`json:"UpdatedAT"`
	Note_Id string `json:"Note_id"`
}