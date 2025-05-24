package routes
import(
	"github.com/gin-gonic/gin"
	controller "golang-restuarent_management/controllers"
)
func OrderItemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/orderitem", controller.GetOrderItems())
	incomingRoutes.GET("/orderitem/:orderitem_id", controller.GetOrderItem())
	incomingRoutes.GET("/orderitem-order/:orderitem_id", controller.GetOrderItemByOrder())
	incomingRoutes.POST("/orderitem", controller.CreateOrderItem())
	incomingRoutes.PATCH("/orderitem/:orderitem_id", controller.UpdateOrderItem())

}