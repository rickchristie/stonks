package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type CORSConfig struct {
	AllowNone      bool
	AllowAll       bool
	AllowedOrigins map[string]bool
}

func ParseCORSConfigFromString(raw string) (*CORSConfig, error) {
	if raw == "" {
		return &CORSConfig{AllowNone: true}, nil
	}

	parts := strings.Split(raw, ",")
	config := &CORSConfig{AllowedOrigins: map[string]bool{}}

	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}

		if origin == "*" {
			config.AllowAll = true
			return config, nil
		}

		if err := validateOrigin(origin); err != nil {
			return nil, fmt.Errorf("invalid origin %q: %w", origin, err)
		}
		config.AllowedOrigins[origin] = true
	}

	if len(config.AllowedOrigins) == 0 {
		return &CORSConfig{AllowNone: true}, nil
	}

	return config, nil
}

func validateOrigin(origin string) error {
	u, err := url.Parse(origin)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme must be http or https")
	}
	if u.Host == "" {
		return fmt.Errorf("missing host")
	}
	if u.Path != "" && u.Path != "/" {
		return fmt.Errorf("origin should not contain path")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return fmt.Errorf("origin should not contain query or fragment")
	}
	return nil
}

func (c *CORSConfig) IsOriginAllowed(origin string) bool {
	if c.AllowNone {
		return false
	}
	if c.AllowAll {
		return true
	}
	return c.AllowedOrigins[origin]
}

func (c *CORSConfig) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin == "" {
			ctx.Next()
			return
		}

		if c.AllowNone {
			if ctx.Request.Method == "OPTIONS" {
				ctx.AbortWithStatus(403)
				return
			}
			ctx.Next()
			return
		}

		if c.AllowAll {
			ctx.Header("Access-Control-Allow-Origin", "*")
		} else if c.IsOriginAllowed(origin) {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Vary", "Origin")
		} else {
			if ctx.Request.Method == "OPTIONS" {
				ctx.AbortWithStatus(403)
				return
			}
			ctx.Next()
			return
		}

		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func MustParseCORSConfig(raw string) *CORSConfig {
	config, err := ParseCORSConfigFromString(raw)
	if err != nil {
		panic(fmt.Sprintf("CORS configuration error: %v", err))
	}
	return config
}
