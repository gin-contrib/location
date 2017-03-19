package main

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// configure to automatically detect scheme and host
	// - use http when default scheme cannot be determined
	// - use localhost:8080 when default host cannot be determined
	router.Use(location.Default())

	router.GET("/", func(c *gin.Context) {
		url := location.Get(c)
		c.String(200, url.String())
	})

	router.Run()
}
