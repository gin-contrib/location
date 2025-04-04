package location

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	location struct {
		scheme  string
		host    string
		base    string
		headers Headers
	}
)

func newLocation(config Config) *location {
	return &location{
		scheme:  config.Scheme,
		host:    config.Host,
		base:    config.Base,
		headers: config.Headers,
	}
}

func (l *location) applyToContext(c *gin.Context) {
	value := new(url.URL)
	value.Scheme = l.resolveScheme(c.Request)
	value.Host = l.resolveHost(c.Request)
	value.Path = l.base
	c.Set(key, value)
}

func (l *location) resolveScheme(r *http.Request) string {
	const HTTPS = "https"
	switch {
	case r.Header.Get(l.headers.Scheme) == HTTPS:
		return HTTPS
	case r.URL.Scheme == HTTPS:
		return HTTPS
	case r.TLS != nil:
		return HTTPS
	case strings.HasPrefix(r.Proto, "HTTPS"):
		return HTTPS
	default:
		return l.scheme
	}
}

func (l *location) resolveHost(r *http.Request) (host string) {
	switch {
	case r.Header.Get(l.headers.Host) != "":
		return r.Header.Get(l.headers.Host)
	case r.Header.Get("X-Host") != "":
		return r.Header.Get("X-Host")
	case r.Host != "":
		return r.Host
	case r.URL.Host != "":
		return r.URL.Host
	default:
		return l.host
	}
}
