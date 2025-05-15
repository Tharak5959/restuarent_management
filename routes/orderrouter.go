package routes
import(
	"github.com/gin-gonic/gin"
	controller"golang-restuarent_management/controllers"
)
func OrderRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/order", controller.GetOrders())
	incomingRoutes.GET("/order/:order_id", controller.GetOrders())
	incomingRoutes.POST("/order", controller.CreateOrder())
	incomingRoutes.PATCH("/order/:order_id", controller.UpdateOrder())

}