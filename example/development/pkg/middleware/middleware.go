package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs every request with method, path, status, latency, and key headers.
// Useful for verifying that global headers (CORS, Auth, custom) are actually being sent.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		log.Printf("[%s] %s %d %v | Auth: %s | Content-Type: %s | X-API-Key: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
			truncate(c.GetHeader("Authorization"), 30),
			c.GetHeader("Content-Type"),
			truncate(c.GetHeader("X-API-Key"), 20),
		)
	}
}

// BearerAuth validates that a Bearer token is present and non-empty.
// Apply to routes that require Bearer authentication.
func BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"error":   "Unauthorized",
				"message": "Missing or invalid Bearer token",
			})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"error":   "Unauthorized",
				"message": "Bearer token is empty",
			})
			return
		}
		c.Set("bearer_token", token)
		c.Next()
	}
}

// APIKeyAuth validates that the X-API-Key header is present and non-empty.
// Apply to routes that require API key authentication.
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"error":   "Unauthorized",
				"message": "Missing X-API-Key header",
			})
			return
		}
		c.Set("api_key", key)
		c.Next()
	}
}

// BasicAuth validates that a Basic auth header is present.
// Apply to routes that require Basic authentication.
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Basic ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"error":   "Unauthorized",
				"message": "Missing or invalid Basic auth credentials",
			})
			return
		}
		c.Next()
	}
}

// CORS adds standard CORS headers to every response.
// This lets you verify in Coco UI that CORS headers are correctly applied globally.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-API-Key, X-Custom-Header")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
