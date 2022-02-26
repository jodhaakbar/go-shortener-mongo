package handler

import (
	"bytes"
	"encoding/json"
	"log"
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
	data := storemongo.RetrieveInitialUrl(shortUrl)
	//fmt.Printf("Found : %s \n", initialUrl)

	if data[0] == "error" {
		c.Redirect(302, storemongo.GoDotEnvVariable("DEFAULT_URL"))
	} else {
		c.Redirect(302, data[0])
		values := map[string]string{"shortUrl": shortUrl}
		jsonData, _ := json.Marshal(values)

		go doPost(jsonData, data[1])
	}

}

func doPost(jsonData []byte, webhook string) {
	_, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}
}
