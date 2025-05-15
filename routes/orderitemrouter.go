package routes
import(
	"github.com/gin-gonic/gin"
	controller "golang-restuarent_management/controllers"
)
func orderitemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderitem", controller.GetOrderItems)
	incomingRoutes.GET("/orderitem/:orderitem_id", controller.GetOrderItemByID)
	incomingRoutes.POST("/orderitem", controller.CreateOrderitem())
	incomingRoutes.PATCH("/orderitem/:orderitem_id", controller.UpdateOrderitem())

}