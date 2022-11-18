package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"os"
	"tiktok-scraper/constants"
)

func VideoHandler(c *gin.Context) {
	id := c.Param("id")
	fileName := fmt.Sprintf("tt_%s.mp4", id)
	path := constants.VideoFolder + fileName
	_, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File does not exist anymore"})
	}

	cd := mime.FormatMediaType("attachment", map[string]string{"filename": fileName})
	c.Header("Content-Disposition", cd)
	c.Header("Content-Type", "application/octet-stream")
	fmt.Println("Sending file.....", fileName)
	c.FileAttachment(path, fileName)
}
