package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "golang-restuarent_management/helpers"
)
func Autentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No token provided"})
			c.Abort()
			return
		}
		claims, err, _ := helper.ValidateToken(clientToken)
		if err!=""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("uid", claims.User_id)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Next()
	}
}