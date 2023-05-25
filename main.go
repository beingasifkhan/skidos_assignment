package main

import (
	"Assignment/controllers"
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
	authGroup.Use(controllers.AuthMiddleware())
	{
		authGroup.GET("/videos", controllers.GetVideos)
		authGroup.GET("/videos/:id", controllers.GetVideosByID)
		authGroup.POST("/videos", controllers.UploadVideo)
	}

	r.Run("localhost:8080")
	fmt.Println("server is running")

}
