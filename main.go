package main

import (
	"assignment_skidos/controllers"
	"assignment_skidos/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//Register/Login
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.Login)

	//Videos
	authGroup := r.Group("/api")
	authGroup.Use(utils.AuthMiddleware())
	{
		authGroup.POST("/videos", controllers.UploadVideo)
	}
	//bitrate
	r.GET("/stream", controllers.StreamVideo)

	r.Run("localhost:8080")
	fmt.Println("server is running")

}
