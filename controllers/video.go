package controllers

import (
	"Assignment/models"
	"Assignment/storage"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var videos []models.Video

func GetVideos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, videos)
}

func GetVideosByID(c *gin.Context) {
	ID := c.Param("id")
	for _, v := range videos {
		if v.ID == ID {
			c.IndentedJSON(http.StatusOK, v)
			return
		}
	}
}

func UploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "no file uploaded",
		})
		return
	}
	videoStorage := storage.NewStorage("your-storage-path")
	err = videoStorage.SaveVideo(file)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save video",
		})
		return
	}
	video := models.Video{
		ID:       generateID(),
		Title:    c.PostForm("title"),
		Duration: calculateDuration(file),
		Size:     int(file.Size / 1024),
	}
	videos = append(videos, video)

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "video uploaded sucessfully",
	})
}

func generateID() string {
	timestamp := time.Now().UnixNano()
	randomNum := rand.Intn(1000)
	id := fmt.Sprintf("%d%d", timestamp, randomNum)
	return id
}

func calculateDuration(file *multipart.FileHeader) int {
	rand.Seed(time.Now().UnixNano())
	duration := rand.Intn(10) + 1
	return duration
}
