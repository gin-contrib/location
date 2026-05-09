package location

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testBarPath = "/bar"
	testBarHost = "bar.com"
	testBarURL  = "http://bar.com/bar"
)

var defaultHeaders = Headers{
	Scheme: HeaderXForwardedProto,
	Host:   HeaderXForwardedHost,
}

var tests = []struct {
	want string
	conf Config
	req  *http.Request
}{
	// defaults
	{
		want: "http://localhost:8080",
		conf: DefaultConfig(),
		req: &http.Request{
			Header: http.Header{},
			URL:    &url.URL{},
		},
	},

	// url scheme
	{
		want: "https://localhost:8080",
		conf: DefaultConfig(),
		req: &http.Request{
			Header: http.Header{},
			URL: &url.URL{
				Scheme: "https",
			},
		},
	},

	// x-forwarded headers
	{
		want: "https://bar.com/bar",
		conf: Config{HTTP, "foo.com", testBarPath, defaultHeaders},
		req: &http.Request{
			Header: http.Header{
				HeaderXForwardedProto: {HTTPS},
				HeaderXForwardedHost:  {testBarHost},
			},
			URL: &url.URL{},
		},
	},

	// X-Host headers
	{
		want: testBarURL,
		conf: Config{HTTP, "foo.com", testBarPath, defaultHeaders},
		req: &http.Request{
			Header: http.Header{
				"X-Host": {testBarHost},
			},
			URL: &url.URL{},
		},
	},

	// URL Host
	{
		want: testBarURL,
		conf: Config{HTTP, "foo.com", testBarPath, defaultHeaders},
		req: &http.Request{
			Header: http.Header{},
			URL: &url.URL{
				Host: testBarHost,
			},
		},
	},

	// requests
	{
		want: "https://baz.com/bar",
		conf: Config{HTTP, "foo.com", testBarPath, defaultHeaders},
		req: &http.Request{
			Proto:  "HTTPS://",
			Host:   "baz.com",
			Header: http.Header{},
			URL:    &url.URL{},
		},
	},

	// tls
	{
		want: "https://foo.com/bar",
		conf: Config{HTTP, "foo.com", testBarPath, defaultHeaders},
		req: &http.Request{
			TLS:    &tls.ConnectionState{},
			Header: http.Header{},
			URL:    &url.URL{},
		},
	},

	// X-Forwarded-Host host header
	{
		want: testBarURL,
		conf: Config{HTTP, "foo.com", testBarPath, Headers{
			Scheme: HeaderXForwardedProto,
			Host:   HeaderXForwardedHost,
		}},
		req: &http.Request{
			Header: http.Header{
				HeaderXForwardedHost: {testBarHost},
			},
			URL: &url.URL{},
		},
	},
}

func TestLocation(t *testing.T) {
	for _, test := range tests {
		c := new(gin.Context)
		c.Request = test.req
		loc := newLocation(test.conf)
		loc.applyToContext(c)

		got := Get(c)

		if got.String() != test.want {
			t.Errorf("wanted location %s, got %s", got.String(), test.want)
		}
	}
}

func defaultRouter() *gin.Engine {
	router := gin.New()
	router.Use(Default())

	router.GET("/", func(c *gin.Context) {
		url := Get(c)
		c.String(200, url.String())
	})

	return router
}

func performRequest(r http.Handler, method string) *httptest.ResponseRecorder {
	req, _ := http.NewRequestWithContext(context.Background(), method, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestDefault(t *testing.T) {
	router := defaultRouter()
	w := performRequest(router, "GET")

	assert.Equal(t, "http://localhost:8080", w.Body.String())
}

func customRouter(config Config) *gin.Engine {
	router := gin.New()
	router.Use(New(config))

	router.GET("/", func(c *gin.Context) {
		url := Get(c)
		c.String(200, url.String())
	})

	return router
}

func TestCustom(t *testing.T) {
	router := customRouter(Config{
		Scheme: "https",
		Host:   "foo.com",
		Base:   "/base",
	})
	w := performRequest(router, "GET")

	assert.Equal(t, "https://foo.com/base", w.Body.String())
}
