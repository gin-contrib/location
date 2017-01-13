package main

import (
	"github.com/marinsalinas/location"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	router := gin.Default()

	// configure to automatically detect scheme and host with
	// fallback to https://foo.com/base
	// - use https when default scheme cannot be determined
	// - use foo.com when default host cannot be determined
	// - include /base as the path
	router.Use(location.New(location.Config{
		Scheme: "https",
		Host:   "foo.com",
		Base:   "/base",
	}))

	router.GET("/", func(c *gin.Context) {
		url := location.Get(c)
		c.String(200, url.String())
	})

	router.Run()
}
