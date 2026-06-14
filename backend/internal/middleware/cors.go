package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	// Support standard Vite local development URLs
	origins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://[::1]:5173",
	}
	if raw := os.Getenv("CORS_ORIGINS"); raw != "" {
		origins = append(origins, strings.Split(raw, ",")...)
	}
	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		for _, allowedOrigin := range origins {
			if origin == allowedOrigin {
				return true
			}
		}
		// Browsers may omit port for HTTP 80 / HTTPS 443
		if origin == "http://localhost" || origin == "http://127.0.0.1" || origin == "http://[::1]" {
			return true
		}
		// Support dynamic ports for local development
		return strings.HasPrefix(origin, "http://localhost:") ||
			strings.HasPrefix(origin, "http://127.0.0.1:") ||
			strings.HasPrefix(origin, "http://[::1]:")
	}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.AllowWebSockets = true
	return cors.New(config)
}
