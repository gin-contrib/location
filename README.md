# location

[![Run Tests](https://github.com/gin-contrib/location/v2/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/gin-contrib/location/v2/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/gin-contrib/location/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-contrib/location)
[![Go Report Card](https://goreportcard.com/badge/github.com/gin-contrib/location/v2)](https://goreportcard.com/report/github.com/gin-contrib/location/v2)
[![GoDoc](https://godoc.org/github.com/gin-contrib/location/v2?status.svg)](https://godoc.org/github.com/gin-contrib/location/v2)

A Gin middleware that automatically detects and exposes the server's hostname and scheme by inspecting the incoming `http.Request`. This is particularly useful when your application runs behind proxies or load balancers, as it intelligently determines the correct public-facing URL.

## Features

- **Automatic Detection**: Intelligently detects scheme (HTTP/HTTPS) and host from various sources
- **Proxy-Aware**: Respects standard proxy headers (`X-Forwarded-Proto`, `X-Forwarded-Host`)
- **Flexible Configuration**: Support for custom headers and fallback values
- **Base Path Support**: Configure base paths for applications running under sub-paths
- **Zero Dependencies**: Only depends on Gin framework

## How It Works

The middleware examines the request in the following order to determine the scheme and host:

**Scheme Detection:**

1. `X-Forwarded-Proto` header (or custom configured header)
2. Request URL scheme
3. TLS connection presence
4. Protocol string
5. Configured default scheme (fallback)

**Host Detection:**

1. Configured custom host header (e.g., `X-Forwarded-Host`)
2. `X-Host` header
3. Request `Host` header
4. Request URL host
5. Configured default host (fallback)

## Usage

### Start using it

Download and install it:

```bash
go get github.com/gin-contrib/location/v2
```

Import it in your code:

```go
import "github.com/gin-contrib/location/v2"
```

### Default Configuration

Use `location.Default()` for quick setup with sensible defaults:

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/gin-contrib/location/v2"
  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()

  // Configure with default settings:
  // - Scheme: defaults to "http"
  // - Host: defaults to "localhost:8080"
  // - Headers: X-Forwarded-Proto (scheme), X-Forwarded-Host (host)
  router.Use(location.Default())

  router.GET("/", func(c *gin.Context) {
    url := location.Get(c)

    // Access the detected URL components
    fmt.Printf("Scheme: %s\n", url.Scheme) // e.g., "https"
    fmt.Printf("Host: %s\n", url.Host)     // e.g., "example.com"
    fmt.Printf("Path: %s\n", url.Path)     // e.g., ""

    c.String(http.StatusOK, "Full URL: %s", url.String())
  })

  router.Run(":8080")
}
```

### Custom Configuration

Customize the middleware for your specific environment:

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/gin-contrib/location/v2"
  "github.com/gin-gonic/gin"
)

func main() {
  router := gin.Default()

  // Custom configuration for production environment behind a proxy
  router.Use(location.New(location.Config{
    Scheme: "https",                    // Default scheme when not detected
    Host:   "api.example.com",          // Default host when not detected
    Base:   "/v1",                      // Base path for the application
    Headers: location.Headers{
      Scheme: "X-Forwarded-Proto",      // Header for scheme detection
      Host:   "X-Forwarded-Host",       // Header for host detection
    },
  }))

  router.GET("/users", func(c *gin.Context) {
    url := location.Get(c)

    // With Base="/v1", url.Path will be "/v1"
    // Full URL will be: https://api.example.com/v1
    c.JSON(http.StatusOK, gin.H{
      "message": "API endpoint",
      "url":     url.String(),
    })
  })

  router.Run(":8080")
}
```

## Configuration Options

| Field            | Type     | Description                      | Default               |
| ---------------- | -------- | -------------------------------- | --------------------- |
| `Scheme`         | `string` | Default scheme when not detected | `"http"`              |
| `Host`           | `string` | Default host when not detected   | `"localhost:8080"`    |
| `Base`           | `string` | Base path for the application    | `""`                  |
| `Headers.Scheme` | `string` | Header name for scheme detection | `"X-Forwarded-Proto"` |
| `Headers.Host`   | `string` | Header name for host detection   | `"X-Forwarded-Host"`  |

## Use Cases

### Behind a Reverse Proxy

When running behind nginx, Apache, or other reverse proxies:

```go
router.Use(location.New(location.Config{
  Scheme: "https",
  Host:   "your-domain.com",
  Headers: location.Headers{
    Scheme: "X-Forwarded-Proto",
    Host:   "X-Forwarded-Host",
  },
}))
```

### Kubernetes Ingress

For applications deployed in Kubernetes with an Ingress controller:

```go
router.Use(location.New(location.Config{
  Scheme: "https",
  Host:   "api.myapp.com",
  Base:   "/api",  // If your ingress routes /api to this service
  Headers: location.Headers{
    Scheme: "X-Forwarded-Proto",
    Host:   "X-Forwarded-Host",
  },
}))
```

### Building Absolute URLs

Generate absolute URLs in your responses:

```go
router.GET("/users/:id", func(c *gin.Context) {
  userID := c.Param("id")
  baseURL := location.Get(c)

  // Build absolute URL for the resource
  userURL := fmt.Sprintf("%s/users/%s", baseURL.String(), userID)

  c.JSON(http.StatusOK, gin.H{
    "id":   userID,
    "self": userURL,
  })
})
```

## API Reference

### Functions

#### `Default() gin.HandlerFunc`

Returns the location middleware with default configuration (HTTP scheme, localhost:8080 host).

#### `New(config Config) gin.HandlerFunc`

Returns the location middleware with custom configuration.

#### `Get(c *gin.Context) *url.URL`

Retrieves the location information from the Gin context. Returns `nil` if the middleware is not configured.

**Returns:** A `*url.URL` with the following fields:

- `Scheme`: The detected or configured scheme (http/https)
- `Host`: The detected or configured host (including port if present)
- `Path`: The configured base path

### Types

#### `Config`

Configuration structure for the middleware:

```go
type Config struct {
    Scheme  string   // Default scheme (e.g., "https")
    Host    string   // Default host (e.g., "example.com:443")
    Base    string   // Base path (e.g., "/api/v1")
    Headers Headers  // Custom headers for detection
}
```

#### `Headers`

Header configuration for scheme and host detection:

```go
type Headers struct {
    Scheme string  // Header name for scheme (default: "X-Forwarded-Proto")
    Host   string  // Header name for host (default: "X-Forwarded-Host")
}
```

## Migration from v1

If you're upgrading from v1, simply update your import path:

```diff
- import "github.com/gin-contrib/location"
+ import "github.com/gin-contrib/location/v2"
```

The API remains fully backward compatible. No code changes are required.

## Contributing

We welcome contributions! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add some amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

Please ensure:

- All tests pass (`go test -v ./...`)
- Code is formatted (`go fmt ./...`)
- Linting passes (`golangci-lint run`)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
