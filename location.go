package location

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

const key = "location"

type Config struct {
	// Scheme is the default scheme that should be used when it cannot otherwise
	// be ascertained from the incoming http.Request.
	Scheme string

	// Host is the default host that should be used when it cannot otherwise
	// be ascertained from the incoming http.Request.
	Host string

	// Base is the base path that should be used in conjunction with proxy
	// servers that do path re-writing.
	Base string
}

func DefaultConfig() Config {
	return Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
}

// Default returns the location middle with defualt configuration.
func Default() gin.HandlerFunc {
	config := DefaultConfig()
	return New(config)
}

// New returns the location middleware with user-defined custom configuration.
func New(config Config) gin.HandlerFunc {
	location := newLocation(config)
	return func(c *gin.Context) {
		location.applyToContext(c)
	}
}

// Get returns the Location information for the incoming http.Request from the
// context. If the location is not set a nil value is returned.
func Get(c *gin.Context) *url.URL {
	v, ok := c.Get(key)
	if !ok {
		return nil
	}
	vv, ok := v.(*url.URL)
	if !ok {
		return nil
	}
	return vv
}
