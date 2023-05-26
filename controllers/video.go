package controllers

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadVideo(c *gin.Context) {
	// Retrieve the uploaded video file
	file, err := c.FormFile("video")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "no file uploaded",
		})
		return
	}

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"), // Replace with your AWS region
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Open the uploaded file
	srcFile, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// Specify the S3 bucket and object key
	bucket := "skikosvideo" // Replace with your S3 bucket name
	objectKey := "videos/" + file.Filename

	// Create an S3 PutObjectInput
	input := &s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(objectKey),
		Body:                 srcFile,
		ServerSideEncryption: aws.String("AES256"),  // Enable SSE with AES-256 encryption
		ACL:                  aws.String("private"), // Set the desired ACL for the object
	}

	// Upload the file to S3
	_, err = svc.PutObject(input)
	if err != nil {
		log.Fatal(err)
	}

	// Add any additional operations you need to perform with the uploaded video, such as saving metadata to a database

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "video uploaded successfully",
	})
}
