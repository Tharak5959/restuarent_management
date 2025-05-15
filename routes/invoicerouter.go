package routes
import(
	"github.com/gin-gonic/gin"
	controller "golang-restuarent_management/controllers"
)
func invoiceRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.GET("/invoice", controller.GetInvoices())
	incomingRoutes.GET("/invoice/:invoice_id", controller.GetInvoice())
	incomingRoutes.POST("/invoice", controller.CreateInvoice())
	incomingRoutes.PATCH("/invoice/:invoice_id", controller.UpdateInvoice())

}