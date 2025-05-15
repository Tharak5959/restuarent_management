package routes
import(
	"github.com/gin-gonic/gin"
	controller "golang-restuarent_management/controllers"
)
func itemRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/item",controller.GetItems())
	incomingRoutes.GET("/item/:item_id",controller.GetItem())
	incomingRoutes.POST("/item",controller.CreateItem())
	incomingRoutes.PATCH("/item/:item_id",controller.UpdateItem())

}