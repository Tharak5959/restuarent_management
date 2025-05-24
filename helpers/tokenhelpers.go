package helper

import (
	"context"
	"fmt"
	"golang-restuarent_management/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type Signeddetails struct {
	Email string
	User_id string
	FirstName string
	LastName string
	jwt.StandardClaims
}
var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var SECRET_KEY string =os.Getenv("SECRET_KEY")

func GenarateAllTokens(email string,uid string,firstname string,lastname string) (signedToken string,refreshToken string,err error){
	claims:= &Signeddetails{
		Email: email,
		User_id: uid,			
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(),
				},
		FirstName: firstname,
		LastName: lastname,
			}
		refreshClaims:= &Signeddetails{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Local().Add(time.Hour * 24 * 30).Unix(),
			},
		}
		token,err:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
		refreshToken,err=jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims).SignedString([]byte(SECRET_KEY))
		if err!=nil{
			log.Panic(err)
			return
		}
		return token,refreshToken,err
	}
func UpdateAllTokens(signedToken string, signedRefreshToken string,userId string)(){
 var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)
 var UpdateObj primitive.D
 UpdateObj=append(UpdateObj, bson.E{"token",signedToken})
 UpdateObj=append(UpdateObj, bson.E{"refresh_token",signedRefreshToken})
 updated_at,err:=time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
 if err!=nil{
		log.Panic(err)
	}
	 UpdateObj=append(UpdateObj, bson.E{"updated_at",updated_at})
	upsert :=true
	filter :=bson.M{"user_id":userId}
	opt :=options.UpdateOptions{
	Upsert: &upsert,
	}
	_,err=userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
		{"$set",UpdateObj},
	},
		&opt,
	)
	defer cancel()
	if err!=nil{
		log.Panic(err)
	}
	return
	 
}
func ValidateToken(signedToken string) (claims *Signeddetails,msg string,err error){
	//parse the token

	token,err:=	jwt.ParseWithClaims(signedToken, &Signeddetails{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
		
		//token is invalid
		claims,ok:=token.Claims.(*Signeddetails)
		if !ok{
			msg=fmt.Sprintf("token is invalid")
			msg=err.Error()
			return
		}
		//token is expired
		if claims.ExpiresAt<time.Now().Local().Unix(){
			msg=fmt.Sprintf("token is expired")
			msg=err.Error()
			return
		}
		return claims,msg,err
}


