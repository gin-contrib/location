package location

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"

	"gopkg.in/gin-gonic/gin.v1"
)

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

	// x-formward headers
	{
		want: "https://bar.com/bar",
		conf: Config{"http", "foo.com", "/bar"},
		req: &http.Request{
			Header: http.Header{
				"X-Forwarded-Proto": {"https"},
				"X-Forwarded-For":   {"bar.com"},
			},
			URL: &url.URL{},
		},
	},

	// X-Host headers
	{
		want: "http://bar.com/bar",
		conf: Config{"http", "foo.com", "/bar"},
		req: &http.Request{
			Header: http.Header{
				"X-Host": {"bar.com"},
			},
			URL: &url.URL{},
		},
	},

	// URL Host
	{
		want: "http://bar.com/bar",
		conf: Config{"http", "foo.com", "/bar"},
		req: &http.Request{
			Header: http.Header{},
			URL: &url.URL{
				Host: "bar.com",
			},
		},
	},

	// requests
	{
		want: "https://baz.com/bar",
		conf: Config{"http", "foo.com", "/bar"},
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
		conf: Config{"http", "foo.com", "/bar"},
		req: &http.Request{
			TLS:    &tls.ConnectionState{},
			Header: http.Header{},
			URL:    &url.URL{},
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
