package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	rawOrigins := strings.Split(allowedOrigins, ",")
	var origins []string
	for _, o := range rawOrigins {
		// Trim whitespace and trailing slash
		cleaned := strings.TrimRight(strings.TrimSpace(o), "/")
		if cleaned != "" {
			origins = append(origins, cleaned)
		}
	}

    // Add logging to verify loaded origins
    // log.Printf("CORS: Loaded allowed origins: %v", origins)

	config := cors.Config{
		AllowOrigins:     origins,
        // Fallback to allow specific common headers if needed
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
			for _, o := range origins {
				if o == "*" {
					return true
				}
				// Exact match
				if o == origin {
					return true
				}
				// Check without trailing slash (just in case config has it)
				if strings.TrimRight(o, "/") == origin {
					return true
				}
			}
			return false
		},
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}
