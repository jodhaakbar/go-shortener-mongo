package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jodhaakbar/go-shortener-mongo/handler"
	"github.com/jodhaakbar/go-shortener-mongo/storemongo"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the URL Shortener API",
		})
	})

	r.POST("/create-short-url", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	port := storemongo.GoDotEnvVariable("PORT")

	err := r.Run(":" + port)
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}

}
