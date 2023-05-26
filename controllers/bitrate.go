package controllers

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// StreamVideo streams the video with adaptive bitrate.
func StreamVideo(c *gin.Context) {
	videoPath := "C:/Users/ANAS KHAN/Videos/Power of stocks course/A complete Trader/A Complete Trader Program 2nd day session -1._HD.mp4"

	// Define the output directory for transcoded video segments
	segmentDir := "C:\\bitrate_video"

	// Define the available bitrates and resolutions for adaptive streaming
	bitrates := []string{"500k", "1000k", "2000k"}              // Example bitrates
	resolutions := []string{"640x360", "1280x720", "1920x1080"} // Example resolutions

	// Transcode the video into different quality levels and resolutions
	for i, bitrate := range bitrates {
		resolution := resolutions[i]
		segmentFilename := fmt.Sprintf("video_%s_%s.ts", bitrate, resolution)
		segmentPath := filepath.Join(segmentDir, segmentFilename)

		// Execute FFmpeg command to transcode the video segment
		cmd := exec.Command("ffmpeg", "-i", videoPath, "-c:v", "libx264", "-b:v", bitrate, "-s", resolution, "-f", "mpegts", "-y", segmentPath)
		err := cmd.Run()
		if err != nil {
			// Handle error
			fmt.Println("Error transcoding video segment:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transcode video segment"})
			return
		}

		// Serve the transcoded video segment
		c.File(segmentPath)

		// Clean up the transcoded video segment file
		err = os.Remove(segmentPath)
		if err != nil {
			// Handle error
			fmt.Println("Error removing video segment:", err)
		}
	}
}
