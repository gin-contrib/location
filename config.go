package location

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"net/url"
	"strings"
)

type (
	location struct {
		scheme  string
		host    string
		base    string
		headers headers
	}

	headers struct {
		scheme string
		host   string
	}
)

func newLocation(config Config) *location {
	return &location{
		scheme: config.Scheme,
		host:   config.Host,
		base:   config.Base,
		headers: headers{
			scheme: "X-Forwarded-Proto",
			host:   "X-Forwarded-For",
		},
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
	switch {
	case r.Header.Get(l.headers.scheme) == "https":
		return "https"
	case r.URL.Scheme == "https":
		return "https"
	case r.TLS != nil:
		return "https"
	case strings.HasPrefix(r.Proto, "HTTPS"):
		return "https"
	default:
		return l.scheme
	}
}

func (l *location) resolveHost(r *http.Request) (host string) {
	switch {
	case r.Header.Get(l.headers.host) != "":
		return r.Header.Get(l.headers.host)
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
