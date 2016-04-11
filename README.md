# gin-location

gin middleware to provide the incoming http.Request hostname and scheme.

## Default

```go
package main

import (
	"github.com/drone/gin-location"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// configure to automatically detect scheme and host
	// - use http when default scheme cannot be determined
	// - use localhost:8080 when default host cannot be determined
	router.Use(location.Default())

	router.Get("/", func(c *gin.Context) {
		url := location.Get(c)
		// url.Scheme
		// url.Host
		// url.Path
	})

	router.Run()
}
```

## Custom

```go
package main

import (
	"github.com/drone/gin-location"
	"github.com/gin-gonic/gin"
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
		Host: "foo.com",
		Path: "/base",
	}))

	router.Get("/", func(c *gin.Context) {
		url := location.Get(c)
		// url.Scheme
		// url.Host
		// url.Path
	})

	router.Run()
}
```