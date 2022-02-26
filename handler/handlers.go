package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jodhaakbar/go-shortener-mongo/shortener"
	"github.com/jodhaakbar/go-shortener-mongo/storemongo"
)

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required,url"`
	UserId  string `json:"user_id" binding:"required"`
	Webhook string `json:"webhook" binding:"required,url"`
}

func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	storemongo.SaveUrlMapping(shortUrl, creationRequest.LongUrl, creationRequest.UserId, creationRequest.Webhook)

	host := storemongo.GoDotEnvVariable("HOST_URL")
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})

}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initialUrl := storemongo.RetrieveInitialUrl(shortUrl)
	fmt.Printf("Found : %s \n", initialUrl)
	c.Redirect(302, initialUrl)
}
