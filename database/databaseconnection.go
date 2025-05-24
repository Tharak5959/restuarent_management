package database
import(
	"fmt"
	"log"
	"time" 
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"


)
func DBinstance() *mongo.Client{
	MongoDb :="mongodb://localhost:27017 "
	fmt.Print(MongoDb)
	client,err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err!=nil{
		log.Fatal(err)
	}
	ctx,cancel:= context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	err =client.Connect(ctx)
	if err!=nil {
		log.Fatal(err)
	}	
	fmt.Println("connected to mongodb")
	return client
}
var client *mongo.Client =DBinstance()
var MenuCollection *mongo.Collection = OpenCollection(client,"menu")
var FoodCollection *mongo.Collection = OpenCollection(client,"food")
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("restaurent").Collection(collectionName)
	return collection
}

var Client *mongo.Client

func init() {
    // Initialize MongoDB client
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI("your-mongodb-uri"))
    if err != nil {
        log.Fatal(err)
    }

    Client = client
}