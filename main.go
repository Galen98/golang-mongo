package main

import (
	"crudgolang/configs"
	"crudgolang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run db
	configs.ConnectDB()
	//routes
	routes.UserRoute(router)
	routes.ProductRoute(router)
	// router.GET("/", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"data": configs.ConnectDB(),
	// 	})
	// })

	router.Run()
}
