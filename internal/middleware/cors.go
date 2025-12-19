package middleware

import (
	"log"
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

	// Log for debugging - remove in production if needed
	log.Printf("CORS: Allowed origins configured: %v", origins)

	config := cors.Config{
		// IMPORTANT: Only use AllowOriginFunc, not AllowOrigins, to avoid conflicts
		AllowOriginFunc: func(origin string) bool {
			cleanedOrigin := strings.TrimRight(origin, "/")
			log.Printf("CORS: Checking origin: %s", cleanedOrigin)
			for _, o := range origins {
				if o == "*" {
					return true
				}
				if o == cleanedOrigin {
					log.Printf("CORS: Origin %s ALLOWED", cleanedOrigin)
					return true
				}
			}
			log.Printf("CORS: Origin %s DENIED", cleanedOrigin)
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}
