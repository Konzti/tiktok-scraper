package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"tiktok-scraper/controllers/helpers"
	"tiktok-scraper/structs"
	"tiktok-scraper/utils"
)

func ScrapeHandler(c *gin.Context) {
	var requestBody structs.URL
	// Checking for valid JSON
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Checking for valid URL ...
	parsedUrl, err := url.ParseRequestURI(requestBody.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if parsedUrl.Host != "www.tiktok.com" {
		if parsedUrl.Host != "vm.tiktok.com" && parsedUrl.Host != "vt.tiktok.com" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		} else {
			secondPart := strings.Split(parsedUrl.Path, "/")[1]
			if len(secondPart) < 4 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
				return
			}

			foundUrl, err := utils.SearchForCanonical(requestBody.URL)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
				return
			}
			newParsedUrl, err := url.ParseRequestURI(foundUrl)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			helpers.ScrapeForVideo(newParsedUrl, c)
			return
		}
	}
	helpers.ScrapeForVideo(parsedUrl, c)

}
