package models
import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Invoice struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	InvoiceId string `bson:"invoice_number,omitempty"`
	OrderID string `bson:"order_id,omitempty"`
	Payment_method string `json:"payment_method" validate:"eq=CARD|eq=CASH|eq=" bson:"payment_method"`
	Payment_due_date time.Time `json:"payment_due_date" validate:"required,eq=pending|eq=paid|eq=overdue"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}