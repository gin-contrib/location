# location

This Gin middleware can be used to automatically find and expose the server's
hostname and scheme by inspecting information in the incoming `http.Request`.
The alternative to this plugin would be explicitly providing such information to
the server as a command line argument or environment variable.

## Usage

### Default

```go
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

		// url.Scheme
		// url.Host
		// url.Path
	})

	router.Run()
}
```

### Custom

```go
package main

import (
	"github.com/gin-contrib/location"
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
		Base: "/base",
	}))

	router.GET("/", func(c *gin.Context) {
		url := location.Get(c)

		// url.Scheme
		// url.Host
		// url.Path
	})

	router.Run()
}
```

## Contributing

Fork -> Patch -> Push -> Pull Request

## License

MIT
