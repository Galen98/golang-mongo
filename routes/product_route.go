package routes

import (
	"crudgolang/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoute(router *gin.Engine) {
	//insert product
	router.POST("/insertProduct", controllers.CreateProduct())
	//get all
	router.GET("/allproducts", controllers.GetAllProducts())
	//get by parameter
	router.GET("/product/:productId", controllers.GetProducts())
	//edit products
	router.PUT("/product/:productId", controllers.EditProduct())
	//delete product
	router.DELETE("/product/:productId", controllers.DeleteProduct())
}
