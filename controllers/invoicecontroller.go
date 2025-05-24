package controller

import (
	"context"
	"fmt"
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
type InvoiceViewformat struct {
	InvoiceID string `json:"invoice_id"`
	Payment_method	string `json:"payment_method"`
	Order_id	string `json:"order_id"`
	PaymentStatus	string `json:"payment_status"`
	Payment_due	interface{} `json:"payment_due"`
	PaymentDueDate time.Time `json:"payment_due_date"`
	TableNumber interface{} `json:"table_number"`
	Order_details interface{} `json:"order_details"`
}
var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func ItemsByOrder(orderID string) ([]bson.M, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := orderCollection.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	var items []bson.M
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}


func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result, err := invoiceCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching invoices"})
			return
		}
		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)

		// if _, err := ItemsByOrder(invoice.OrderID); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching order items"})
		// 	return
		// }
		defer cancel()
	}
}
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		Invoice_ID := c.Param("invoice_id")
		var invoice models.Invoice
		err:= invoiceCollection.FindOne(ctx, bson.M{"invoice_id": Invoice_ID}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching Invoice"})
			return
		}
		var invoiceViewformat InvoiceViewformat
		allOrderItems,err := ItemsByOrder(invoice.OrderID)
		invoiceViewformat.Order_id= invoice.OrderID
		invoiceViewformat.PaymentDueDate= invoice.Payment_due_date
		invoiceViewformat.Payment_method = "null"
		if invoice.Payment_method != "" {
			invoiceViewformat.Payment_method= *&invoice.Payment_method
		}
		invoiceViewformat.InvoiceID= invoice.InvoiceId
		invoiceViewformat.PaymentStatus=*&invoice.PaymentStatus
		invoiceViewformat.Payment_due= allOrderItems[0]["payment_due"]
		invoiceViewformat.TableNumber= allOrderItems[0]["table_number"]
		invoiceViewformat.Order_details= allOrderItems[0]["order_details"]
		c.JSON(http.StatusOK, invoiceViewformat)
	}
}
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.OrderID}).Decode(&order)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching Order"})
			return
		}
		status:="pending"
		if invoice.PaymentStatus == "" {
			invoice.PaymentStatus = status
		}
		invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
		invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Id = primitive.NewObjectID()
		invoice.InvoiceId = invoice.Id.Hex()
		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		result, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			msg := "invoice was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)


	}

}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceID := c.Param("invoice_id")
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		filter := bson.M{"invoice_id": invoiceID}
		var updatedObj bson.D
		if invoice.Payment_method != "" {
			updatedObj = append(updatedObj, bson.E{Key: "payment_method", Value: invoice.Payment_method})
		}
		if invoice.PaymentStatus != "" {
			updatedObj = append(updatedObj, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})
		}
		invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})
		upsert :=true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status:="pending"		
		if invoice.PaymentStatus != "" {
			invoice.PaymentStatus= status
		}
		_,err:=invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: updatedObj},
			},
			&opt,
		)
		if err != nil {
			msg := fmt.Sprintf("invoice with id %v not updated", invoiceID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, invoice)
	}
}
