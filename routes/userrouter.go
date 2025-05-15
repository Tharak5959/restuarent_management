package routes
import(
	"github.com/gin-gonic/gin"
	controller "golang-restuarent_management/controllers"
)
func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/users", controller.Getusers())
	incomingRoutes.GET("/users/:user_id", controller.Getuser())
	incomingRoutes.POST("/users/signup", controller.Signup())
	incomingRoutes.POST("/users/login", controller.Login())

}